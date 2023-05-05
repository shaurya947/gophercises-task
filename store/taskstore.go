package store

import (
	"encoding/binary"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

var incompleteTasksBucket = []byte("incompleteTasks")
var completedTasksBucket = []byte("completedTasks")

type TaskStore struct {
	db *bolt.DB

	// Used for getting current time when marking as task as completed.
	// Making this a non-exported struct variable allows us to use different
	// values in tests as desired.
	nowFunc func() time.Time
}

func NewTaskStore(filepath string) (*TaskStore, error) {
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &TaskStore{db: db, nowFunc: time.Now}, nil
}

func (ts *TaskStore) Close() {
	ts.db.Close()
}

func (ts *TaskStore) AddTask(t *Task) error {
	return ts.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(incompleteTasksBucket)
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()
		buf, err := t.Marshal()
		if err != nil {
			return err
		}

		return b.Put(itob(id), buf)
	})
}

func (ts *TaskStore) GetIncompleteTasks() ([]*Task, error) {
	var incompleteTasks []*Task

	err := ts.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(incompleteTasksBucket)
		if b == nil {
			return fmt.Errorf("Bucket %s does not exist",
				incompleteTasksBucket)
		}

		err := b.ForEach(func(k, v []byte) error {
			t := &Task{}
			err := t.Unmarshal(v)
			if err != nil {
				return err
			}
			t.id = btoi(k)
			incompleteTasks = append(incompleteTasks, t)
			return nil
		})
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return incompleteTasks, nil
}

// Returns tasks that were completed since the provided time
func (ts *TaskStore) GetCompletedTasks(since time.Time) ([]*Task, error) {
	var completedTasks []*Task

	err := ts.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(completedTasksBucket)
		if b == nil {
			return nil
		}

		err := b.ForEach(func(k, v []byte) error {
			t := &Task{}
			err := t.Unmarshal(v)
			if err != nil {
				return err
			}
			if t.CompletionTime >= since.Unix() {
				t.id = btoi(k)
				completedTasks = append(completedTasks, t)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return completedTasks, nil
}

// taskNum represents 1-indexed position of task in incomplete bucket
// Returns list of tasks that were marked as completed if successful.
func (ts *TaskStore) CompleteTasks(taskNums []int) ([]*Task, error) {
	var completedTasks []*Task

	incompleteTasks, err := ts.GetIncompleteTasks()
	if err != nil {
		return nil, err
	}

	err = ts.db.Update(func(tx *bolt.Tx) error {
		bIncomplete := tx.Bucket(incompleteTasksBucket)
		if bIncomplete == nil {
			return fmt.Errorf("Bucket %s does not exist",
				incompleteTasksBucket)
		}

		bCompleted, err := tx.CreateBucketIfNotExists(
			completedTasksBucket)
		if err != nil {
			return err
		}

		now := ts.nowFunc()
		for i := 0; i < len(taskNums); i++ {
			if taskNums[i]-1 >= len(incompleteTasks) {
				return fmt.Errorf("Invalid task number %d",
					taskNums[i])
			}
			t := incompleteTasks[taskNums[i]-1]
			bIncomplete.Delete(itob(t.id))

			t.CompletionTime = now.Unix()
			id, _ := bCompleted.NextSequence()
			buf, err := t.Marshal()
			if err != nil {
				return err
			}

			err = bCompleted.Put(itob(id), buf)
			if err != nil {
				return err
			}
			completedTasks = append(completedTasks, t)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return completedTasks, nil
}

// taskNum represents 1-indexed position of task in incomplete bucket
// Returns list of tasks that were removed if successful.
func (ts *TaskStore) RemoveTasks(taskNums []int) ([]*Task, error) {
	var removedTasks []*Task

	incompleteTasks, err := ts.GetIncompleteTasks()
	if err != nil {
		return nil, err
	}

	err = ts.db.Update(func(tx *bolt.Tx) error {
		bIncomplete := tx.Bucket(incompleteTasksBucket)
		if bIncomplete == nil {
			return fmt.Errorf("Bucket %s does not exist",
				incompleteTasksBucket)
		}

		for i := 0; i < len(taskNums); i++ {
			if taskNums[i]-1 >= len(incompleteTasks) {
				return fmt.Errorf("Invalid task number %d",
					taskNums[i])
			}
			t := incompleteTasks[taskNums[i]-1]
			bIncomplete.Delete(itob(t.id))
			removedTasks = append(removedTasks, t)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return removedTasks, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
