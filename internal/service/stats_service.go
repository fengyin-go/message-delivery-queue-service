package service

import (
	"sort"

	"message-queue/internal/model"
)

type GlobalStats struct {
	TopicCount      int `json:"topic_count"`
	MessageCount    int `json:"message_count"`
	DeadLetterCount int `json:"dead_letter_count"`
	ProducerCount   int `json:"producer_count"`
	ConsumerCount   int `json:"consumer_count"`
	GroupCount      int `json:"group_count"`
	PendingCount    int `json:"pending_count"`
	DeliveredCount  int `json:"delivered_count"`
	AckedCount      int `json:"acked_count"`
}

func (s *Service) GetGlobalStats() (*GlobalStats, error) {
	stats := &GlobalStats{}
	stats.TopicCount = len(s.store.ListTopics())
	stats.ProducerCount = len(s.store.ListProducers())
	stats.ConsumerCount = len(s.store.ListConsumers())
	stats.GroupCount = len(s.store.ListConsumerGroups())
	messages := s.store.ListMessages()
	stats.MessageCount = len(messages)
	for _, m := range messages {
		switch m.Status {
		case model.MessageStatusPending:
			stats.PendingCount++
		case model.MessageStatusDelivered:
			stats.DeliveredCount++
		case model.MessageStatusAcknowledged:
			stats.AckedCount++
		}
	}
	stats.DeadLetterCount = len(s.store.ListDeadLetters())
	return stats, nil
}

type TopicMessageStats struct {
	TopicID        string `json:"topic_id"`
	TopicName      string `json:"topic_name"`
	MessageCount   int    `json:"message_count"`
	PendingCount   int    `json:"pending_count"`
	DeliveredCount int    `json:"delivered_count"`
	AckedCount     int    `json:"acked_count"`
	DeadCount      int    `json:"dead_count"`
}

func (s *Service) GetTopicMessageStats() ([]*TopicMessageStats, error) {
	topics := s.store.ListTopics()
	messages := s.store.ListMessages()
	deadLetters := s.store.ListDeadLetters()

	statsMap := make(map[string]*TopicMessageStats, len(topics))
	for _, t := range topics {
		statsMap[t.ID] = &TopicMessageStats{
			TopicID:   t.ID,
			TopicName: t.Name,
		}
	}
	for _, m := range messages {
		st, ok := statsMap[m.TopicID]
		if !ok {
			continue
		}
		st.MessageCount++
		switch m.Status {
		case model.MessageStatusPending:
			st.PendingCount++
		case model.MessageStatusDelivered:
			st.DeliveredCount++
		case model.MessageStatusAcknowledged:
			st.AckedCount++
		case model.MessageStatusDead:
			st.DeadCount++
		}
	}
	for _, d := range deadLetters {
		st, ok := statsMap[d.TopicID]
		if ok {
			st.DeadCount++
		}
	}
	result := make([]*TopicMessageStats, 0, len(statsMap))
	for _, st := range statsMap {
		result = append(result, st)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].MessageCount > result[j].MessageCount
	})
	return result, nil
}

type TopicRateStats struct {
	TopicID    string  `json:"topic_id"`
	TopicName  string  `json:"topic_name"`
	Total      int     `json:"total"`
	Acked      int     `json:"acked"`
	AckRate    float64 `json:"ack_rate"`
	Delivered  int     `json:"delivered"`
	Throughput float64 `json:"throughput"`
}

func (s *Service) GetTopicRateStats() ([]*TopicRateStats, error) {
	topics := s.store.ListTopics()
	messages := s.store.ListMessages()

	statsMap := make(map[string]*TopicRateStats, len(topics))
	for _, t := range topics {
		statsMap[t.ID] = &TopicRateStats{
			TopicID:   t.ID,
			TopicName: t.Name,
		}
	}
	for _, m := range messages {
		st, ok := statsMap[m.TopicID]
		if !ok {
			continue
		}
		st.Total++
		if m.Status == model.MessageStatusAcknowledged {
			st.Acked++
		}
		if m.Status == model.MessageStatusDelivered || m.Status == model.MessageStatusAcknowledged {
			st.Delivered++
		}
	}
	result := make([]*TopicRateStats, 0, len(statsMap))
	for _, st := range statsMap {
		if st.Total > 0 {
			st.AckRate = float64(st.Acked) / float64(st.Total)
			st.Throughput = float64(st.Delivered) / float64(st.Total)
		}
		result = append(result, st)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})
	return result, nil
}

type DeadLetterTopN struct {
	TopicID      string `json:"topic_id"`
	TopicName    string `json:"topic_name"`
	DeadCount    int    `json:"dead_count"`
}

func (s *Service) GetDeadLetterTopN(n int) ([]*DeadLetterTopN, error) {
	topics := s.store.ListTopics()
	deadLetters := s.store.ListDeadLetters()

	countMap := make(map[string]int, len(topics))
	nameMap := make(map[string]string, len(topics))
	for _, t := range topics {
		countMap[t.ID] = 0
		nameMap[t.ID] = t.Name
	}
	for _, d := range deadLetters {
		countMap[d.TopicID]++
	}
	result := make([]*DeadLetterTopN, 0, len(topics))
	for tid, cnt := range countMap {
		result = append(result, &DeadLetterTopN{
			TopicID:   tid,
			TopicName: nameMap[tid],
			DeadCount: cnt,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].DeadCount > result[j].DeadCount
	})
	if n > 0 && n < len(result) {
		result = result[:n]
	}
	return result, nil
}
