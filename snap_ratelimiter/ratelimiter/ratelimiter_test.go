package ratelimiter

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestSlidingWindowBruteForce(t *testing.T) {
	userid := "nur"
	ctx := context.Background()
	MAX := int64(3)
	rc := NewRateLimitter()
	for i := 0; i < 10; i++ {
		isGreater, err := rc.SlidingWindow(ctx, userid, time.Now(), MAX)
		if err != nil {
			log.Println(err)
		}

		if isGreater {
			log.Println("isGreater")
			break
		}
	}

}

func TestSlidingWindowDelay(t *testing.T) {
	userid := "nur"
	ctx := context.Background()
	MAX := int64(30)
	rc := NewRateLimitter()
	for i := 0; i < 10; i++ {
		isGreater, err := rc.SlidingWindow(ctx, userid, time.Now(), MAX)
		if err != nil {
			log.Println(err)
		}

		if isGreater {
			log.Println("isGreater")
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func TestFixWindowDelay(t *testing.T) {
	userid := "nur"
	ctx := context.Background()
	MAX := int64(30)
	rc := NewRateLimitter()
	for i := 0; i < 100; i++ {
		isGreater, err := rc.FixedWindow(ctx, userid, time.Now(), MAX)
		if err != nil {
			log.Println(err)
		}

		if isGreater {
			log.Println("isGreater")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
