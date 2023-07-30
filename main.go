package main

import (
	"github.com/BobyMCbobs/todo-list-etcd/pkg/etcd"
	"github.com/BobyMCbobs/todo-list-etcd/pkg/httpserver"
	"github.com/BobyMCbobs/todo-list-etcd/pkg/todolist"
)

func main() {
	clientset, err := etcd.NewClient()
	if err != nil {
		panic(err)
	}
	mgr := todolist.NewManager(clientset)
	httpserver.NewHTTPServer(mgr).Run()
}
