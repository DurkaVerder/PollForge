package service

import (
	"auth/internal/models"
	"errors"
	"log"

)

// Типы заданий
type JobType string

const (
	JobRegister JobType = "register"
	JobLogin    JobType = "login"
)

type Job struct {
	Type    JobType
	Request models.UserRequest
	Result  chan JobResult
}

type JobResult struct {
	Token string
	Error error
}

var jobQueue chan Job

// Запуск пула воркеров
func StartWorkerPool(numWorkers int) {
	jobQueue = make(chan Job, 100)

	for i := 0; i < numWorkers; i++ {
		go worker(i, jobQueue)
	}
	log.Printf("Запущен пул из %d воркеров", numWorkers)
}
// Обработка заданий
func worker(id int, jobs <-chan Job) {
	for job := range jobs {

		log.Printf("Worker %d: обрабатывает задачу %s", id, job.Type)
		
		var token string
		var err error

		switch job.Type {
		case JobRegister:
			token, err = handleRegistration(job.Request)
		case JobLogin:
			token, err = handleLogin(job.Request)
		default:
			err = errors.New("неизвестный тип задачи")
		}

		job.Result <- JobResult{Token: token, Error: err}
	}
}


func AsyncRegisterUser(request models.UserRequest) (string, error) {
	resultChan := make(chan JobResult)
	jobQueue <- Job{Type: JobRegister, Request: request, Result: resultChan}
	res := <-resultChan
	return res.Token, res.Error
}

func AsyncLoginUser(request models.UserRequest) (string, error) {
	resultChan := make(chan JobResult)
	jobQueue <- Job{Type: JobLogin, Request: request, Result: resultChan}
	res := <-resultChan
	return res.Token, res.Error
}
