package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/BobyMCbobs/todo-list-etcd/pkg/common"
	"github.com/BobyMCbobs/todo-list-etcd/pkg/todolist"
)

type HTTPServer struct {
	server          *http.Server
	todolistManager *todolist.Manager
	port            string
}

func NewHTTPServer(todolistManager *todolist.Manager) *HTTPServer {
	h := &HTTPServer{
		port:            common.GetAppPort(),
		todolistManager: todolistManager,
	}
	c := cors.New(cors.Options{})
	router := mux.NewRouter().StrictSlash(true)
	router.Use(Logging)
	router.Use(gziphandler.GzipHandler)
	router.Use(c.Handler)
	apiRouters := router.PathPrefix("/api").Subrouter()
	// apiRouters.HandleFunc("", h.apiGet())
	apiListsRouters := apiRouters.PathPrefix("/list").Subrouter()
	// apiListsRouters.Use(h.MiddlewareValidateJWT)
	apiListsRouters.HandleFunc("", h.apiListLists()).Methods(http.MethodGet)
	apiListsRouters.HandleFunc("/{id}", h.apiGetList()).Methods(http.MethodGet)
	apiListsRouters.HandleFunc("/{id}", h.apiPutList()).Methods(http.MethodPut)
	apiListsRouters.HandleFunc("/{id}", h.apiDeleteList()).Methods(http.MethodDelete)
	apiListsRouters.HandleFunc("", h.apiPostList()).Methods(http.MethodPost)
	apiListsRouters.HandleFunc("/{listid}/item", h.apiListItems()).Methods(http.MethodGet)
	apiListsRouters.HandleFunc("/{listid}/item/{id}", h.apiGetItem()).Methods(http.MethodGet)
	apiListsRouters.HandleFunc("/{listid}/item/{id}", h.apiPutItem()).Methods(http.MethodPut)
	apiListsRouters.HandleFunc("/{listid}/item/{id}", h.apiDeleteItem()).Methods(http.MethodDelete)
	apiListsRouters.HandleFunc("/{listid}/item", h.apiDeleteItemAll()).Methods(http.MethodDelete)
	apiListsRouters.HandleFunc("/{listid}/item", h.apiPostItem()).Methods(http.MethodPost)

	// webFolderPath := common.GetWebFolder()
	// webHandler := http.FileServer(http.Dir(webFolderPath))
	// router.PathPrefix("/").Handler(webHandler).Methods(http.MethodGet)

	s := &http.Server{
		Handler:           router,
		Addr:              h.port,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	h.server = s

	return h
}

func (h *HTTPServer) Run() {
	log.Println("HTTP listening on", h.server.Addr)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	log.Println("Shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server didn't exit gracefully %v", err)
	}
}
