package main

import (
	APIs "YTSearchAPI/APIs"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("starting server...")

	fmt.Println("starting worker...")
	var wg sync.WaitGroup
	wg.Add(1)
	go APIs.FetchAndStoreVideos(&wg)
	fmt.Println("registering APIs...")

	fmt.Println("good to go")

	wg.Wait()
}
