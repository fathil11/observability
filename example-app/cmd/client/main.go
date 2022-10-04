package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func() {
			_, err := http.Get("http://localhost:5000/product")
			if err != nil {
				log.Println(err)
				return
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
