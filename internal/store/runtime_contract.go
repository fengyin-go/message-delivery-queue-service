package store

import "errors"

type ConsumerGroupStreamStore struct{}

func NewConsumerGroupStreamStore() *ConsumerGroupStreamStore { return &ConsumerGroupStreamStore{} }
func (s *ConsumerGroupStreamStore) Stream(fail bool) (<-chan string, <-chan error) {
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
