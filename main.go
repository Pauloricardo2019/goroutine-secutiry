package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type queueGeneric struct {
	msgChannel chan string
}

func NewQueue(msg chan string) *queueGeneric {
	return &queueGeneric{
		msgChannel: msg,
	}
}

func (q *queueGeneric) startQueue(ctx context.Context) {
	defer close(q.msgChannel)
	timeTicker := time.NewTicker(time.Second * 3)

	for {
		select {
		case <-timeTicker.C:
			q.msgChannel <- "enviando mensagem..."

		case <-ctx.Done():
			fmt.Println("finalizando queue...")
			return
		}
	}

}

func main() {

	queue := NewQueue(make(chan string, 10))

	ctx := context.TODO()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go queue.startQueue(ctx)

	for workerCount := 0; workerCount < 10; workerCount++ {
		wg.Add(1)
		go worker(ctx, workerCount, &wg, queue.msgChannel)
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
