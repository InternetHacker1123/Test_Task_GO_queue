package main

import (
	"fmt"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d started job %d\n", id, job)
        // Выполнение работы здесь
        results <- job * 2 // Результат работы
    }
}

func main() {
    numJobs := 20
    numWorkers := 3

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // Запуск нескольких горутин-рабочих
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // Подготовка и отправка задач в очередь
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // Получение результатов
    for a := 1; a <= numJobs; a++ {
        result := <-results
        fmt.Println("Result:", result)
    }
}
