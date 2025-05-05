package worker

import (
	"agent/internal/tasks/gen"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

const errorResult = "FAILED"

func Work(taskCh <-chan *gen.TaskResponse, client gen.OrchestratorClient) {
	for task := range taskCh {
		processTask(task, client)
	}
}

func processTask(task *gen.TaskResponse, client gen.OrchestratorClient) {
	resultChan := make(chan float64, 1)
	errorChan := make(chan error, 1)

	go func() {
		res, err := calculateResult(task)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- res
	}()

	select {
	case <-time.After(time.Duration(task.Time) * time.Millisecond):
		log.Printf("Task %s timed out\n", task.Id)
		err := sendResult(client, task.Id, errorResult)
		if err != nil {
			log.Printf("Error sending result for task %s: %s\n", task.Id, err)
		}
	case err := <-errorChan:
		log.Printf("Error calculating result: %v\n", err)
		err = sendResult(client, task.Id, errorResult)
		if err != nil {
			log.Printf("Error sending result for task %s: %s\n", task.Id, err)
		}
	case result := <-resultChan:
		err := sendResult(client, task.Id, result)
		if err != nil {
			log.Printf("Error sending result for task %s: %s\n", task.Id, err)
		}
	}
}

// calculateResult is a method for calculating result of task
func calculateResult(task *gen.TaskResponse) (float64, error) {
	arg1, arg2, err := parseArgs(task.Arg1, task.Arg2)
	if err != nil {
		return 0, err
	}

	switch task.Operator {
	case gen.Operator_ADDICTION:
		return arg1 + arg2, nil
	case gen.Operator_SUBTRACTION:
		return arg1 - arg2, nil
	case gen.Operator_MULTIPLICATION:
		return arg1 * arg2, nil
	case gen.Operator_DIVISION:
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return arg1 / arg2, nil
	}

	return 0, fmt.Errorf("unknown operator")
}

func parseArgs(arg1, arg2 string) (float64, float64, error) {
	if arg1 == "" || arg2 == "" {
		return 0, 0, fmt.Errorf("missing arguments")
	}

	val1, err := strconv.ParseFloat(arg1, 64)
	if err != nil {
		return 0, 0, err
	}

	val2, err := strconv.ParseFloat(arg2, 64)
	if err != nil {
		return 0, 0, err
	}

	return val1, val2, nil
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// sendResult is a method for sending calculation result to the API
func sendResult(client gen.OrchestratorClient, taskID string, result interface{}) error {
	data := gen.UpdateTaskRequest{
		Id:     taskID,
		Result: toString(result),
	}

	_, err := client.UpdateTask(context.Background(), &data)
	return err
}
