package store

import "errors"

type RetryBatchStreamStore struct{}

func NewRetryBatchStreamStore() *RetryBatchStreamStore { return &RetryBatchStreamStore{} }
func (s *RetryBatchStreamStore) Stream(fail bool) (<-chan string, <-chan error) {
	results := make(chan string)
	errs := make(chan error, 1)
	go func() {
		if fail {
			errs <- errors.New("partition unavailable")
			return
		}
		results <- "ready"
	}()
	return results, errs
}
