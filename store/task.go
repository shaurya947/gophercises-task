package store

import (
	"github.com/shaurya947/gophercises-task/store/internal"
	"google.golang.org/protobuf/proto"
)

type Task struct {
	id             uint64
	Description    string
	CompletionTime int64
}

func (t *Task) Marshal() ([]byte, error) {
	return proto.Marshal(&internal.Task{
		Description:    t.Description,
		CompletionTime: t.CompletionTime,
	})
}

func (t *Task) Unmarshal(data []byte) error {
	var pb internal.Task
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	t.Description = pb.Description
	t.CompletionTime = pb.CompletionTime
	return nil
}
