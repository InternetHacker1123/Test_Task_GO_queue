package main

import (
	"fmt"
	"io"
	"net/http"
)

func getReq(URL string) ([]byte, error) {
    resp, err := http.Get(URL)
    if err != nil {
        fmt.Println("Ошибка при выполнении GET запроса:", err)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Ошибка при чтении ответа:", err)
        return nil, err
    }
    return body, nil
}

func worker(tasks chan string, results chan []byte, f func(string) ([]byte, error)) {
    for task := range tasks {
        result, error := f(task)
		if error != nil {
			fmt.Println(error)
		}
        results <- result
    }
}


func queue(numWorkers int, numTasks int, tasks chan string, results chan []byte, data string) {

    for c := 1; c <= numWorkers; c++ {
		go worker(tasks, results, getReq)
		fmt.Println("Горутина запущена")
	}
        

	for c := 1; c <= numTasks; c++ {
        tasks <- data
    }
	close(tasks)

    for c := 1; c <= numTasks; c++ {
        result := <-results
        fmt.Println("Результат:", string(result))
    }
}

func main() {
	numWorkers := 3
    numTasks := 20
	
    tasks := make(chan string, numTasks)
    results := make(chan []byte, numTasks)

    data := "https://postman-rest-api-learner.glitch.me/info"

    queue(numWorkers, numTasks, tasks, results, data)
}