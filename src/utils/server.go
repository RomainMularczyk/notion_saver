package utils

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"notion_saver/src/database"
	"notion_saver/src/queue"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Server struct {
	Database   *gorm.DB
	Queue      *amqp.Connection
	HTTPClient *http.Client
	Logger     *slog.Logger
	Env        map[string]string
	Context    context.Context
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func (s *Server) ConnectDatabase() {
	db := database.Database{}
	db.Connect(
		s.Env["POSTGRES_HOST"],
		s.Env["POSTGRES_USER"],
		s.Env["POSTGRES_PASSWORD"],
		s.Env["POSTGRES_DB"],
		s.Env["POSTGRES_PORT"],
		s.Env["POSTGRES_SSLMODE"],
		s.Env["POSTGRES_TIMEZONE"],
	)
	s.Database = db.Client
}

func (s *Server) ConnectQueue() {
	q := queue.Queue{}
	q.Connect(
		s.Env["RABBITMQ_HOST"],
		s.Env["RABBITMQ_USER"],
		s.Env["RABBITMQ_PASSWORD"],
		s.Env["RABBITMQ_PORT"],
	)
	q.SetupTopology()
	s.Queue = q.Client
}
