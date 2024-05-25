package main

import (
	"fmt"
	"io"
	"net/http"
)

// func getReq() {
// 	URL := "https://postman-rest-api-learner.glitch.me/info"

// }


func getReq(URLs <-chan string, results chan<- string) {
	for URL := range URLs {
		resp, err := http.Get(URL)
    if err != nil {
        fmt.Println("Ошибка при выполнении GET запроса:", err)
		defer resp.Body.Close()
        return
    }

	body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Ошибка при чтении ответа:", err)
		defer resp.Body.Close()
        return
    }
	results <- string(body)
    defer resp.Body.Close()
	}
}

func main() {
	URLs := make(chan string)
	results := make(chan string)
	URL := "https://postman-rest-api-learner.glitch.me/info"

	for c := 1; c <= 1000; c ++ {
		go getReq(URLs, results)
	}

	for c := 1; c <= 1000; c ++ {
		URLs <- URL 
	}
	close(URLs)

	for c := 1; c <= 1000; c++ {
        result := <-results
        fmt.Println("Result:", result)
    }
}

