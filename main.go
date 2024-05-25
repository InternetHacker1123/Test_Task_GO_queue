package main

import (
	"fmt"
	"io"
	"net/http"
)

func getReq(URL string) ([]byte, error){
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


func queue(f func(string) ([]byte, error), numWorkers int, numTasks int, tasks chan string, results chan string, data string) {
    for c := 1; c <= numWorkers; c++ {
        go func() {
            for task := range tasks {
                result, err := f(task)
                if err != nil {
                    fmt.Println("Ошибка при выполнении задачи:", err)
                } else {
                    results <- string(result)
                }
            }
        }()
    }

    for c := 1; c <= numTasks; c++ {
        tasks <- data
    }
    close(tasks)

    for c := 1; c <= numTasks; c++ {
        result := <-results
        fmt.Println("Результат:", result)
    }
}



func main() {
	tasks := make(chan string)
	results := make(chan string)
	num_workers := 10
	num_tasks := 1000
	data := "https://postman-rest-api-learner.glitch.me/info"


	go queue(getReq, num_workers, num_tasks, tasks, results, data)



}


