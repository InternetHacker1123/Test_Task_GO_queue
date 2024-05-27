package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)


type Task struct {
    id int
    job Data
}

type Result struct {
    id int
    job Data
    result Data
}


type Data interface{}

//функция обработчик тасков 2
func inc(task Task) (Result, error) {
    res := task.job.(int)
    res++
    result := Result{id: task.id, job: task.job, result: res}
    return result, nil
}



// функция обработчик тасков
func getReq(task Task) (Result, error) {
    resp, err := http.Get(task.job.(string))
    if err != nil {
        return Result{id: task.id, job: task.job, result: "Ошибка при выполнении запроса"}, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return Result{id: task.id, job: task.job, result: "Ошибка при чтении ответа"}, err
    }
    result := Result{id: task.id, job: task.job, result: string(body)}
    return result, nil
}

// функция-воркер в очереди 
func worker(tasks chan Task, results chan Result, f func(task Task) (Result, error), wg *sync.WaitGroup) {
    defer wg.Done()
    for task := range tasks {
        result, err := f(task)

		if err != nil {
            fmt.Println(err)
            continue
		}

        taskResult := result

        results <- taskResult
    }
}


// Функция запуска воркеров
func runWorkers(numWorkers int, tasks chan Task, results chan Result, f func(task Task) (Result, error)) {
    wg := new(sync.WaitGroup)
    for c := 1; c <= numWorkers; c++ {
        wg.Add(1)
		go worker(tasks, results, f, wg)
		fmt.Println("Горутина запущена")
	}
    wg.Wait()
}


// функция показа результатов выполнения тасков
func showResults(results chan Result) {
    defer close(results)
    for c := 1; c <= len(results); c++ {
        result := <-results
        fmt.Printf("Результат: id:%d\ntask: %s\nresult: %s\n\n\n", rune(result.id), result.job, result.result)
        
    }
    
    
}

// очередь
func queue(numWorkers int, numTasks int, tasks chan Task, results chan Result, data Data, f func(task Task) (Result, error)) {
    
    switch data.(type) {

    case int: formattingData := int(data.(int))
        for c := 1; c <= numTasks; c++ {
            task := Task{id: c, job: formattingData}
            tasks <- task
        }
        close(tasks)
        runWorkers(numWorkers, tasks, results, f)
        showResults(results)


    case string: formattingData := string(data.(string))
        for c := 1; c <= numTasks; c++ {
            task := Task{id: c, job: formattingData}
            tasks <- task
        }
        close(tasks)
        runWorkers(numWorkers, tasks, results, f)
        showResults(results)


    case []int: formattingData := []int(data.([]int))
        for c := 0; c <= len(formattingData) - 1; c++ {
            task := Task{id: c, job: formattingData[c]}
            tasks <- task
        }
        close(tasks)
        runWorkers(numWorkers, tasks, results, f)
        showResults(results)


    case []string: formattingData := []string(data.([]string))
        for c := 0; c <= len(formattingData) - 1; c++ {
            task := Task{id: c, job: formattingData[c]}
            tasks <- task
        }
        close(tasks)
        runWorkers(numWorkers, tasks, results, f)
        showResults(results)
    }
    
}

func main() {
	numWorkers := 3
    numTasks := 20
	
    tasks := make(chan Task, numTasks)
    results := make(chan Result, numTasks)

    data := []string{"123", "43124413", "https://postman-rest-api-learner.glitch.me/info"} 

    queue(numWorkers, numTasks, tasks, results, data, getReq)
}