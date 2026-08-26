package store

import "errors"

type PartitionResultStreamStore struct{}

func NewPartitionResultStreamStore() *PartitionResultStreamStore {
	return &PartitionResultStreamStore{}
}
func (s *PartitionResultStreamStore) Stream(fail bool) (<-chan string, <-chan error) {
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
