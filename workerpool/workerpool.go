package workerpool

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

func StartWorkerpool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	fmt.Println("Start")
	for i := 0; i <= viper.GetInt("cnt_workers"); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker := NewWorker(i)
			if err := worker.work(ctx); err != nil {
				log.Fatal(err)
			}
			log.Printf("work done: %d", worker.ID)
		}()
	}
	fmt.Println("END")
}
