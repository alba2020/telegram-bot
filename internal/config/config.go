package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Token                string   `yaml:"token"`
	PGDSN                string   `yaml:"pg_dsn"`
	ReportsServerAddress string   `yaml:"reports_server_address"`
	Topic                string   `yaml:"topic"`
	BrokersList          []string `yaml:"brokers_list"`
	ConsumerGroup        string   `yaml:"consumer_group"`
}

type Service struct {
	config Config
}

func New(configFile string) (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) PGDSN() string {
	return s.config.PGDSN
}

func (s *Service) ReportsServerAddress() string {
	return s.config.ReportsServerAddress
}

func (s *Service) Topic() string {
	return s.config.Topic
}

func (s *Service) BrokersList() []string {
	return s.config.BrokersList
}

func (s *Service) ConsumerGroup() string {
	return s.config.ConsumerGroup
}
