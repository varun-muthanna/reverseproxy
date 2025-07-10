package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/varun-muthanna/loadbalancer/test/handler"
	"github.com/gorilla/mux"
)

var connectionCount int32

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run backend.go <port>")
		return
	}
	port := os.Args[1]
	address := ":" + port

	r := mux.NewRouter()

	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/",handler.GetRouter)

	putRouter := r.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/",handler.PutRouter)

	postRouter := r.Methods("POST").Subrouter()
	postRouter.HandleFunc("/",handler.PostRouter)

	deleteRouter := r.Methods("POST").Subrouter()
	deleteRouter.HandleFunc("/",handler.DeleteRouter)

	s := &http.Server{
		Addr: address,
		Handler: r,
	}
	
	ch := make(chan os.Signal,1)

	go func(){
		fmt.Printf("Server active and listening on %s \n",address)
		
		err := s.ListenAndServe()

		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", address, err)
		}

	}()
	
	signal.Notify(ch,os.Interrupt)
	signal.Notify(ch,syscall.SIGTERM)

	sig := <-ch

	fmt.Println("Gracefull initiated afer recieving ",sig)

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	defer s.Shutdown(ctx)
}
