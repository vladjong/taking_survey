package workerpool

import (
	"context"
	"fmt"

	"github.com/vladjong/taking_survey/client"
)

type worker struct {
	ID int
}

func NewWorker(ID int) *worker {
	return &worker{
		ID: ID,
	}
}

func (w *worker) work(ctx context.Context) error {
	client := client.NewClinet(ctx)
	if err := client.Run(); err != nil {
		return fmt.Errorf("error: %s worker id: %d", err.Error(), w.ID)
	}
	return nil
}
