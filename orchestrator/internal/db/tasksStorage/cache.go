package tasksStorage

import (
	"context"
	"github.com/valkey-io/valkey-go"
	"orchestrator/internal/controllers/utils"
)

const dbConn = "127.0.0.1:6379"

type Cache struct {
	Storage valkey.Client
}

func (c *Cache) GetTasks(ctx context.Context) (*[]utils.InternalTask, error) {
	array, err := c.Storage.Do(ctx, c.Storage.B().Keys().Pattern("*").Build()).ToArray()
	if err != nil {
		return nil, err
	}

	tasks := []utils.InternalTask{}
	for _, item := range array {
		id, err := item.ToString()
		if err != nil {
			return nil, err
		}

		task, err := c.GetTask(ctx, id)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	return &tasks, nil
}

func (c *Cache) GetTask(ctx context.Context, id string) (*utils.InternalTask, error) {
	item, err := c.Storage.Do(ctx, c.Storage.B().Get().Key(id).Build()).ToString()
	if err != nil {
		return nil, err
	}

	task, err := utils.Decode(item)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (c *Cache) SetTask(ctx context.Context, task *utils.InternalTask) error {
	if err := c.Storage.Do(ctx, c.Storage.B().Set().Key(task.ID).
		Value(utils.Encode(*task)).Build()).Error(); err != nil {
		return err
	}

	if err := c.Storage.Do(ctx, c.Storage.B().Expire().Key(task.ID).
		Seconds(120).Build()).Error(); err != nil {
		return err
	}

	return nil
}

func NewCache() (*Cache, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{dbConn}})
	if err != nil {
		return nil, err
	}
	return &Cache{Storage: client}, nil
}
