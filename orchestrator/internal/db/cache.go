package db

import (
	"context"
	"errors"
	"github.com/valkey-io/valkey-go"
	"orchestrator/internal/controllers/dto"
	"orchestrator/pkg/calc"
)

var errStillProcessing = errors.New("still processing")

type Cache struct {
	Storage valkey.Client
}

func (c *Cache) GetTasks(ctx context.Context) (*[]calc.Task, error) {
	array, err := c.Storage.Do(ctx, c.Storage.B().Keys().Pattern("*").Build()).ToArray()
	if err != nil {
		return nil, err
	}

	tasks := []calc.Task{}
	tasksDone := 0
	for _, item := range array {
		id, err := item.ToString()
		if err != nil {
			return nil, err
		}

		task, err := c.GetTask(ctx, id)
		if err != nil {
			return nil, err
		}

		if val, ok := task.Result.(string); ok {
			if val == dto.Failed {
				tasksDone++
			}
		} else {
			tasksDone++
		}

		tasks = append(tasks, *task)
	}

	if len(tasks) == tasksDone {
		return &[]calc.Task{}, c.Clear(ctx)
	}

	return &tasks, nil
}

func (c *Cache) GetTask(ctx context.Context, id string) (*calc.Task, error) {
	item, err := c.Storage.Do(ctx, c.Storage.B().Get().Key(id).Build()).ToString()
	if err != nil {
		return nil, err
	}

	task, err := calc.Decode(item)
	if err != nil {
		return nil, err
	}

	if err := c.UpdateTaskArgs(ctx, task); err != nil && !errors.Is(err, errStillProcessing) {
		return nil, err
	}

	return task, nil
}

func (c *Cache) SetTask(ctx context.Context, task *calc.Task) error {
	if err := c.Storage.Do(ctx, c.Storage.B().Set().Key(task.ID).
		Value(calc.Encode(*task)).Build()).Error(); err != nil {
		return err
	}

	/*
		if err := c.Storage.Do(ctx, c.Storage.B().Expire().Key(task.ID).
			Seconds(120).Build()).Error(); err != nil {
			return err
		}
	*/

	return nil
}

/*
func (c *Cache) DeleteTask(ctx context.Context, id string) error {
	return c.Storage.Do(ctx, c.Storage.B().Del().Key(id).Build()).Error()
}
*/

func (c *Cache) UpdateTaskArgs(ctx context.Context, task *calc.Task) error {
	hasErr := false

	if arg1Str, ok := task.Arg1.(string); ok {
		arg1task, err := c.GetTask(ctx, arg1Str)
		if err != nil {
			return err
		}
		if val, ok := arg1task.Result.(string); ok {
			if val == dto.Processing {
				return errStillProcessing
			} else {
				hasErr = true
			}
		} else {
			task.Arg1 = arg1task.Result
		}
	}

	if arg2Str, ok := task.Arg2.(string); ok {
		arg2task, err := c.GetTask(ctx, arg2Str)
		if err != nil {
			return err
		}
		if val, ok := arg2task.Result.(string); ok {
			if val == dto.Processing {
				return errStillProcessing
			} else {
				hasErr = true
			}
		} else {
			task.Arg2 = arg2task.Result
		}
	}

	if hasErr {
		task.Result = dto.Failed
	}

	if err := c.SetTask(ctx, task); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Clear(ctx context.Context) error {
	return c.Storage.Do(ctx, c.Storage.B().Flushdb().Build()).Error()
}

func NewCache(conn string) (*Cache, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{conn}})
	if err != nil {
		return nil, err
	}
	return &Cache{Storage: client}, nil
}
