package store

import "errors"

type RetryBatchStreamStore struct{}

func NewRetryBatchStreamStore() *RetryBatchStreamStore { return &RetryBatchStreamStore{} }
func (s *RetryBatchStreamStore) Stream(fail bool) (<-chan string, <-chan error) {
	results := make(chan string)
	errs := make(chan error, 1)
	go func() {
		// 退出前关闭两个通道：results 保证消费端的 for-range 收尾，
		// errs 保证消费端最终读取错误时不会因无数据而永久阻塞。
		defer close(results)
		defer close(errs)
		if fail {
			errs <- errors.New("partition unavailable")
			return
		}
		results <- "ready"
	}()
	return results, errs
}
