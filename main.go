package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetch(ctx context.Context, n int, fail bool) (string, error) {

	// time.Sleep(time.Duration(n) * time.Second)

	select {
	case <-time.After(time.Duration(n) * time.Second):
	// sleep finished uninterrupted
	case <-ctx.Done():
		return "", ctx.Err()
	}

	if fail {
		return "", errors.New("an error") //fmt.Errorf("an error")
	} else {
		return "Hello", nil
	}
}

func getData(ctx context.Context, ch chan string, n int, fail bool) error {
	fmt.Println("fetch data" + strconv.Itoa(n))

	if response, err := fetch(ctx, n, fail); err != nil {
		fmt.Println("error encountered at ", strconv.Itoa(n))
		return err
	} else {
		ch <- response
		fmt.Println("fetched data" + strconv.Itoa(n))
		return nil
	}

}

func main() {
	res, err := func() (string, error) {
		ctx := context.Background()
		g, ctx := errgroup.WithContext(ctx)

		ch1 := make(chan string, 1)
		ch2 := make(chan string, 1)
		ch3 := make(chan string, 1)

		g.Go(func() error {
			defer close(ch1)
			return getData(ctx, ch1, 1, true)
		})

		g.Go(func() error {
			defer close(ch2)
			return getData(ctx, ch2, 2, false)
		})

		g.Go(func() error {
			defer close(ch3)
			return getData(ctx, ch3, 3, false)
		})

		result := <-ch1 + <-ch2 + <-ch3
		return result, g.Wait()
	}()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
