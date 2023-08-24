package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	test2()
}

func test1() {
	forever := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				forever <- struct{}{}
				return
			default:
				fmt.Println("for loop")
			}

			time.Sleep(500 * time.Millisecond)
		}
	}(ctx)

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	<-forever
	fmt.Println("finish")
}

func test2() {

	// var g errgroup.Group

	eg, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < 5; i++ {
		i := i
		eg.Go(func() error {
			return insertJob(ctx, i)
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Println("error group")
	} else {
		fmt.Println("success group")
	}
}

func insertJob(ctx context.Context, i int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := insert(i); err != nil {
			return err
		}
		return nil
	}
}
func insert(i int) error {
	if i == 2 {
		return errors.New("error insert 3")
	}
	fmt.Printf("do job %d.\n", i)
	return nil
}
