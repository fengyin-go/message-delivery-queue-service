package service

import "message-queue/internal/store"

type RetryBatchStreamCoordinator struct{ backend *store.RetryBatchStreamStore }

func NewRetryBatchStreamCoordinator(b *store.RetryBatchStreamStore) *RetryBatchStreamCoordinator {
	return &RetryBatchStreamCoordinator{backend: b}
}
func (c *RetryBatchStreamCoordinator) Collect(fail bool) (values []string, err error) {
	results, errs := c.backend.Stream(fail)
	for value := range results {
		values = append(values, value)
	}
	if err := <-errs; err != nil {
		return values, err
	}
	return values, nil
}
