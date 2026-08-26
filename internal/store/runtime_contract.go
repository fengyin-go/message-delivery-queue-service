package store

import "sync"

type VersionedState struct {
	Version int
	State   string
}
type MessageDeliveryVersionStore struct {
	mu     sync.Mutex
	states map[string]VersionedState
}

func NewMessageDeliveryVersionStore() *MessageDeliveryVersionStore {
	return &MessageDeliveryVersionStore{states: make(map[string]VersionedState)}
}
// Update 以单调递增的版本号提交投递终态。
// 仅当 version 严格大于当前已提交版本时才接受写入并返回 true；
// 任何来自旧回调的 version<=current 写入（版本倒退或同版本重放）
// 一律拒绝并返回 false，保证已提交的终态不被迟到回调覆盖。
func (s *MessageDeliveryVersionStore) Update(key string, version int, state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.states[key]
	if ok && version <= cur.Version {
		// 旧回调晚到：版本倒退或同版本重放，拒绝覆盖已提交终态。
		return false
	}
	s.states[key] = VersionedState{Version: version, State: state}
	return true
}
func (s *MessageDeliveryVersionStore) Get(key string) VersionedState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[key]
}
