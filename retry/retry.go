package retry

import (
	"math/rand"
	"time"

	iLog "github.com/g-airport/go-infra/log"
)

func Retry(attempts int, sleep time.Duration, f func() error,
	logFunc func(error) string) error {
	err := f()
	if err == nil {
		return nil
	}

	if attempts--; attempts <= 0 {
		return err
	}

	// Output warn logs
	output := logFunc(err)
	if output != "" {
		iLog.Warnw("retry", output)
	}

	// Add some randomness to prevent creating a Thundering Herd
	jitter := time.Duration(rand.Int63n(int64(sleep)))
	sleep = sleep + jitter/2

	time.Sleep(sleep)
	return Retry(attempts, 2*sleep, f, logFunc)
}
