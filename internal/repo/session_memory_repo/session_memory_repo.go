package sessionmemoryrepo

import (
	"context"
	"errors"
	"sync"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

type Repo struct {
	mu       *sync.RWMutex
	sessions map[int64]domain.Session
}

func New() (*Repo, error) {
	return &Repo{
		mu:       &sync.RWMutex{},
		sessions: make(map[int64]domain.Session),
	}, nil
}

func (r *Repo) Find(
	ctx context.Context,
	userId int64) (*domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sessions[userId]
	if ok {
		return &s, nil
	}
	return nil, errors.New("session not found")
}

func (r *Repo) GetOrCreate(
	ctx context.Context,
	userId int64) (*domain.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.sessions[userId]
	if ok {
		logger.Debug("Session found", s)
		return &s, nil
	} else {
		logger.Debug("REPO session not found")
		psession, _ := domain.NewSession(userId)
		r.sessions[userId] = *psession
		return psession, nil
	}
}

func (r *Repo) Save(ctx context.Context, psession *domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if psession.UserId == 0 {
		return errors.New("invalid user id for session")
	}
	r.sessions[psession.UserId] = *psession
	return nil
}
