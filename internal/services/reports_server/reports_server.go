package reportsserver

import (
	"context"
	"fmt"
	"net"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	pb "gitlab.ozon.dev/albatros2002/telegram-bot/reports/reports"
	"google.golang.org/grpc"
)

var port int = 50051

type Sender interface {
	SendMessage(domain.Message) error
}

type server struct {
	pb.UnimplementedReporterServer
	sender Sender
}

func (s *server) SendReport(ctx context.Context, in *pb.ReportRequest) (*pb.ReportResponse, error) {
	logger.Infof("gRPC received for user %d: %v",
		in.GetUserId(),
		in.GetMessage(),
	)

	msg := domain.Message{
		Text:   in.GetMessage(),
		UserID: in.GetUserId(),
	}
	err := s.sender.SendMessage(msg)
	if err != nil {
		logger.Error(err)
	}

	return &pb.ReportResponse{Status: "ok"}, nil
}

func StartReportsServer(sender Sender) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("Failed to listen", err)
	}
	s := grpc.NewServer()

	pb.RegisterReporterServer(s, &server{sender: sender})
	logger.Infof("gRPC server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", err)
	}
}
