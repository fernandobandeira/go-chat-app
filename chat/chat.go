package chat

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/fernandomalmeida/go-chat-app/chat/hub"
	"github.com/fernandomalmeida/go-chat-app/dbgen/dbchat"
	"github.com/fernandomalmeida/go-chat-app/util/fileserver"
	"github.com/fernandomalmeida/go-chat-app/util/views"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var (
	viewsFolder  = filepath.FromSlash("chat/views")
	sharedFolder = "shared"
	staticFolder = "static"
)

type chatServer struct {
	baseView *views.Base
	store    dbchat.Store
}

func New(chatStore dbchat.Store) (cs *chatServer, err error) {
	cs = &chatServer{
		baseView: views.NewBase(
			viewsFolder,
			sharedFolder,
		),
		store: chatStore,
	}

	return
}

func (cs *chatServer) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", cs.Index())
	r.Get("/room", cs.Room())

	r.Get("/ws", cs.ServeWs())

	staticPath := filepath.Join(viewsFolder, staticFolder)
	fileserver.Server(r, "/static", staticPath)

	return r
}

func (cs *chatServer) Index() http.HandlerFunc {
	tIndex := cs.baseView.Parse(filepath.Join("pages", "index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		tIndex.ExecuteTemplate(w, "master", nil)
	}
}

func (cs *chatServer) Room() http.HandlerFunc {
	tIndex := cs.baseView.Parse(filepath.Join("pages", "room.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		tIndex.ExecuteTemplate(w, "master", nil)
	}
}

func (cs *chatServer) ServeWs() http.HandlerFunc {
	h := hub.New(cs.store)
	go h.Run()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("websocket upgrade!")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error on upgrade websocket: %v", err)
		}

		hub.NewClient(h, conn)
	}
}
