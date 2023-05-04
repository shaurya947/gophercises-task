package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTask(t *testing.T) {
	s := createTestStore(t)
	defer s.close()

	t1 := &Task{
		Description: "do dishes",
	}

	s.TaskStore.AddTask(t1)

	tasks, err := s.TaskStore.GetIncompleteTasks()
	require.Nil(t, err)
	require.Len(t, tasks, 1)
	require.Equal(t, t1.Description, tasks[0].Description)
}

type store struct {
	*TaskStore
}

func createTestStore(t *testing.T) *store {
	ts, err := NewTaskStore(t.TempDir() + "/test.db")
	require.Nil(t, err)
	return &store{ts}
}

func (s *store) close() {
	s.TaskStore.Close()
}
