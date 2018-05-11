package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

const maxInterrupts = 2

var (
	interruptCounter = 0
)

func main() {
	mainCtx, cancle := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-mainCtx.Done():
				fmt.Println("I'm done here")
				return

			case t := <-ticker.C:
				fmt.Printf("Tick at %v\n", t)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			interruptCounter++
			fmt.Printf("This is interrupt #%d\n", interruptCounter)

			if interruptCounter >= maxInterrupts {
				cancle()

				fmt.Println("I have to leave. Bye")
			}
		}
	}()

	<-mainCtx.Done()
}
