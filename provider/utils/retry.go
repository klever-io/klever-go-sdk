package utils

import (
	"fmt"
	"time"
)

func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)
		sleep *= 2
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
