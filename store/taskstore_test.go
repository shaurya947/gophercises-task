package store

import (
	"testing"
	"time"

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

func TestGetIncompleteTasks(t *testing.T) {
	s := createTestStore(t)
	defer s.close()

	tasks := []*Task{
		{Description: "do dishes"},
		{Description: "wash clothes"},
		{Description: "buy groceries"},
	}

	for _, t := range tasks {
		s.TaskStore.AddTask(t)
	}

	incompleteTasks, err := s.TaskStore.GetIncompleteTasks()
	require.Nil(t, err)
	require.Len(t, incompleteTasks, 3)
	require.Equal(t, "do dishes", incompleteTasks[0].Description)
	require.Equal(t, "wash clothes", incompleteTasks[1].Description)
	require.Equal(t, "buy groceries", incompleteTasks[2].Description)
}

func TestGetCompletedTasksAndCompleteTasks(t *testing.T) {
	s := createTestStore(t)
	defer s.close()

	tasks := []*Task{
		{Description: "do dishes"},
		{Description: "wash clothes"},
		{Description: "buy groceries"},
		{Description: "vacuum the basement"},
	}

	for _, t := range tasks {
		s.TaskStore.AddTask(t)
	}

	fullDayAgo := time.Now().Add(-time.Hour * 24)
	completedTasks, err := s.TaskStore.GetCompletedTasks(fullDayAgo)
	require.Nil(t, err)
	require.Empty(t, completedTasks)

	completedTasks, err = s.TaskStore.CompleteTasks([]int{1, 4})
	require.Nil(t, err)
	require.Len(t, completedTasks, 2)
	require.Equal(t, "do dishes", completedTasks[0].Description)
	require.Equal(t, "vacuum the basement", completedTasks[1].Description)

	completedTasks, err = s.TaskStore.GetCompletedTasks(fullDayAgo)
	require.Len(t, completedTasks, 2)
	require.Equal(t, "do dishes", completedTasks[0].Description)
	require.Equal(t, "vacuum the basement", completedTasks[1].Description)
	require.NotZero(t, completedTasks[0].CompletionTime)
	require.NotZero(t, completedTasks[1].CompletionTime)

	incompleteTasks, err := s.TaskStore.GetIncompleteTasks()
	require.Nil(t, err)
	require.Len(t, incompleteTasks, 2)
	require.Equal(t, "wash clothes", incompleteTasks[0].Description)
	require.Equal(t, "buy groceries", incompleteTasks[1].Description)
	require.Zero(t, incompleteTasks[0].CompletionTime)
	require.Zero(t, incompleteTasks[1].CompletionTime)

	// change nowFunc to be 25 hours ago
	s.TaskStore.nowFunc = func() time.Time {
		return time.Now().Add(-time.Hour * 25)
	}

	completedTasks, err = s.TaskStore.CompleteTasks([]int{2})
	require.Nil(t, err)
	require.Len(t, completedTasks, 1)
	require.Equal(t, "buy groceries", completedTasks[0].Description)

	// completed tasks since 24 hours ago shouldn't include "buy groceries"
	completedTasks, err = s.TaskStore.GetCompletedTasks(fullDayAgo)
	require.Len(t, completedTasks, 2)
	require.Equal(t, "do dishes", completedTasks[0].Description)
	require.Equal(t, "vacuum the basement", completedTasks[1].Description)
	require.NotZero(t, completedTasks[0].CompletionTime)
	require.NotZero(t, completedTasks[1].CompletionTime)

	incompleteTasks, err = s.TaskStore.GetIncompleteTasks()
	require.Nil(t, err)
	require.Len(t, incompleteTasks, 1)
	require.Equal(t, "wash clothes", incompleteTasks[0].Description)
	require.Zero(t, incompleteTasks[0].CompletionTime)

	// but completed tasks since 26 hours ago should include "buy groceries"
	completedTasks, err = s.TaskStore.GetCompletedTasks(fullDayAgo.Add(
		-time.Hour * 2))
	require.Len(t, completedTasks, 3)
	require.Equal(t, "do dishes", completedTasks[0].Description)
	require.Equal(t, "vacuum the basement", completedTasks[1].Description)
	require.Equal(t, "buy groceries", completedTasks[2].Description)
	require.NotZero(t, completedTasks[0].CompletionTime)
	require.NotZero(t, completedTasks[1].CompletionTime)
	require.NotZero(t, completedTasks[2].CompletionTime)
}

func TestRemoveTasks(t *testing.T) {
	s := createTestStore(t)
	defer s.close()

	tasks := []*Task{
		{Description: "do dishes"},
		{Description: "wash clothes"},
		{Description: "buy groceries"},
		{Description: "vacuum the basement"},
	}

	for _, t := range tasks {
		s.TaskStore.AddTask(t)
	}

	// ...
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
