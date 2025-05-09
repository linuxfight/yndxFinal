package db

import (
	"context"
	"github.com/valkey-io/valkey-go"
	"testing"
	"time"

	"orchestrator/pkg/calc"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCacheIntegration(t *testing.T) {
	ctx := context.Background()

	// Start Valkey container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "valkey/valkey:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer func(container testcontainers.Container, ctx context.Context, opts ...testcontainers.TerminateOption) {
		err := container.Terminate(ctx)
		if err != nil {
			panic(err)
		}
	}(container, ctx)

	// Get connection string
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)
	connStr := host + ":" + port.Port()

	// Create cache client
	cache, err := NewCache(connStr)
	require.NoError(t, err)
	defer cache.Storage.Close()

	t.Run("TestGetTasks", func(t *testing.T) {
		// Setup test data
		task1 := createTestTask(t, "1+2")
		task2 := createTestTask(t, "3*4")
		require.NoError(t, cache.SetTask(ctx, task1))
		require.NoError(t, cache.SetTask(ctx, task2))

		// Test GetTasks
		tasks, err := cache.GetTasks(ctx)
		require.NoError(t, err)
		assert.Len(t, *tasks, 2)

		// Cleanup
		require.NoError(t, cache.Clear(ctx))
	})

	t.Run("TestGetTask", func(t *testing.T) {
		task := createTestTask(t, "5/2")
		require.NoError(t, cache.SetTask(ctx, task))

		// Test existing task
		fetched, err := cache.GetTask(ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, fetched.ID)

		// Test non-existing task
		_, err = cache.GetTask(ctx, "non-existent")
		assert.Error(t, err)

		require.NoError(t, cache.Clear(ctx))
	})

	t.Run("TestSetTask", func(t *testing.T) {
		task := createTestTask(t, "1+16")

		err := cache.SetTask(ctx, task)
		require.NoError(t, err)

		// Verify storage
		res, err := cache.Storage.Do(ctx, cache.Storage.B().Get().Key(task.ID).Build()).ToString()
		require.NoError(t, err)
		assert.NotEmpty(t, res)

		require.NoError(t, cache.Clear(ctx))
	})

	t.Run("TestUpdateTaskArgs", func(t *testing.T) {
		// Create dependent tasks
		arg1Task := createTestTask(t, "2+3")
		arg2Task := createTestTask(t, "4*5")
		parentTask := &calc.Task{
			ID:   uuid.NewString(),
			Arg1: arg1Task.ID,
			Arg2: arg2Task.ID,
		}

		// Test with unresolved dependencies
		err := cache.UpdateTaskArgs(ctx, parentTask)
		assert.ErrorIs(t, err, valkey.Nil)

		// Resolve dependencies
		arg1Task.Result = 5.0
		arg2Task.Result = 20.0
		require.NoError(t, cache.SetTask(ctx, arg1Task))
		require.NoError(t, cache.SetTask(ctx, arg2Task))

		// Test with resolved dependencies
		require.NoError(t, cache.UpdateTaskArgs(ctx, parentTask))
		assert.IsType(t, float64(0), parentTask.Arg1)
		assert.IsType(t, float64(0), parentTask.Arg2)

		require.NoError(t, cache.Clear(ctx))
	})

	t.Run("TestClear", func(t *testing.T) {
		task := createTestTask(t, "1+1")
		require.NoError(t, cache.SetTask(ctx, task))

		require.NoError(t, cache.Clear(ctx))

		tasks, err := cache.GetTasks(ctx)
		require.NoError(t, err)
		assert.Empty(t, *tasks)
	})
}

func createTestTask(t *testing.T, expr string) *calc.Task {
	tasks, err := calc.ParseExpression(expr)
	require.NoError(t, err)
	require.NotEmpty(t, tasks)
	return &tasks[0]
}
