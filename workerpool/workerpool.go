package workerpool

import (
	"context"
	"log"
	"sync"

	"github.com/spf13/viper"
)

func StartWorkerpool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	log.Println("workerpool start")
	for i := 1; i <= viper.GetInt("cnt_workers"); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			worker := NewWorker(i)
			if err := worker.work(ctx); err != nil {
				cancel()
				wg.Wait()
				return
			}
			log.Printf("work done: %d", worker.ID)
		}(i)
	}
	wg.Wait()
	log.Println("workerpool end")
}
