package workerpool

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/viper"
)

func StartWorkerpool() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Println("workerpool start")
	for i := 1; i <= viper.GetInt("cnt_workers"); i++ {
		select {
		case signal := <-sigs:
			log.Printf("signal: %d received", signal)
			cancel()
			wg.Wait()
			return
		default:
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				worker := NewWorker(i)
				if err := worker.work(ctx); err != nil {
					cancel()
					log.Fatal(err)
				}
				log.Printf("work done: %d", worker.ID)
			}(i)
		}
	}
	wg.Wait()
	log.Println("workerpool end")
}
