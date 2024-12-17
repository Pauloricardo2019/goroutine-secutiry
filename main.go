package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type queueGeneric struct {
}

func NewQueue() *queueGeneric {
	return &queueGeneric{}
}

func (q *queueGeneric) startQueue(ctx context.Context, wg *sync.WaitGroup, msg chan string) {
	defer func() {
		close(msg)
		wg.Done()
	}()

	timeTicker := time.NewTicker(time.Second * 3)

	for {
		select {
		case <-timeTicker.C:
			msg <- "enviando mensagem..."

		case <-ctx.Done():
			fmt.Println("finalizando queue...")
			return
		}
	}

}

var msg = make(chan string, 10)

func main() {

	queue := NewQueue()

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*8)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go queue.startQueue(ctx, &wg, msg)

	for workerCount := 0; workerCount < 10; workerCount++ {
		wg.Add(1)
		go worker(ctx, workerCount, &wg, msg)
	}

	wg.Wait()

}

func worker(ctx context.Context, workerId int, group *sync.WaitGroup, msg chan string) {
	fmt.Println("Initialing workerId: ", workerId)
	defer group.Done()

	for {
		select {
		case mensagem, ok := <-msg:
			if !ok {
				fmt.Println("closing worker, workerId: ", workerId)
				return
			}

			fmt.Printf("Recebi a msg usando o worker: %d, a msg é: %s\n", workerId, mensagem)

		case <-ctx.Done():
			fmt.Println("closing worker, workerId: ", workerId)
			return
		}
	}

}
