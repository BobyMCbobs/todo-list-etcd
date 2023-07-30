package main

import (
	"fmt"

	"github.com/BobyMCbobs/todo-list-etcd/pkg/etcd"
)

func main() {
	clientset, err := etcd.NewClient()
	if err != nil {
		panic(err)
	}
	kvs, err := clientset.ListWithPrefix("")
	if err != nil {
		panic(err)
	}
	fmt.Println("kvs length:", len(kvs))
	for _, v := range kvs {
		fmt.Println("key:", string(v.Key))
		fmt.Println("value:", string(v.Value))
	}
	if _, err := clientset.Put("/thing/a", "hello"); err != nil {
		panic(err)
	}
	if _, err := clientset.Put("/thing/b", "hello"); err != nil {
		panic(err)
	}
	val, err := clientset.Get("/thing/a")
	if err != nil {
		panic(err)
	}
	fmt.Println("value:", string(val.Value))
	kvs, err = clientset.ListWithPrefix("/thing")
	if err != nil {
		panic(err)
	}
	fmt.Println("kvs length:", len(kvs))
	for _, v := range kvs {
		fmt.Println("key:", string(v.Key))
		fmt.Println("value:", string(v.Value))
	}
	if _, err := clientset.Delete("/thing/a"); err != nil {
		panic(err)
	}
	if _, err := clientset.Delete("/thing/b"); err != nil {
		panic(err)
	}
}
