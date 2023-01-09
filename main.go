package main

import (
	APIs "YTSearchAPI/APIs"
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("starting server...")
	fmt.Println("starting worker...")
	var wg sync.WaitGroup
	wg.Add(1)
	go APIs.FetchAndStoreVideos(&wg)
	fmt.Println("registering APIs...")
	router := mux.NewRouter()
	router.Use(APIs.SetContentTypeAsJson)
	router.HandleFunc("/api/GetVideos", APIs.GetVideos).Methods("GET")
	router.HandleFunc("/api/Search", APIs.Search).Methods("POST")
	http.Handle("/", router)
	http.ListenAndServe(":7777", router)
	fmt.Println("good to go....")
	wg.Wait()
}
