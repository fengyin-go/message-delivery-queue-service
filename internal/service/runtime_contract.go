package service

import "message-queue/internal/store"

type RetryBatchStreamCoordinator struct{ backend *store.RetryBatchStreamStore }

func NewRetryBatchStreamCoordinator(b *store.RetryBatchStreamStore) *RetryBatchStreamCoordinator {
	return &RetryBatchStreamCoordinator{backend: b}
}
func (c *RetryBatchStreamCoordinator) Collect(fail bool) (values []string, err error) {
	results, errs := c.backend.Stream(fail)
	// 排空结果流，待生产端关闭 results 后再读错误通道。
	// 失败分支下生产端先发错误、再关闭 results，排空后这里才能读到错误并及时返回，
	// 否则会一直挂在 for-range 上，整次 Collect 无法返回。
	for value := range results {
		values = append(values, value)
	}
	if err := <-errs; err != nil {
		return values, err
	}
	return values, nil
}
