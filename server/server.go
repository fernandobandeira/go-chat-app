package server

import (
	"database/sql"
	"net/http"

	"github.com/fernandomalmeida/go-chat-app/chat"
	"github.com/fernandomalmeida/go-chat-app/config"
	"github.com/fernandomalmeida/go-chat-app/dbgen/dbchat"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type server struct {
}

func New(conf config.Config) (srv *http.Server, err error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	sqlDB, err := sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		return
	}
	chatStore := dbchat.NewStore(sqlDB)
	chatServer, err := chat.New(chatStore)
	if err != nil {
		return
	}
	r.Mount("/chat", chatServer.Routes())

	srv = &http.Server{
		Addr:    conf.ServerAddress,
		Handler: r,
	}

	return srv, nil
}
