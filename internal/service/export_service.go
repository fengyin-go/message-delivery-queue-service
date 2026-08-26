package service

import (
	"message-queue/internal/model"
)

type TopicExport struct {
	TopicID     string                `json:"topic_id"`
	TopicName   string                `json:"topic_name"`
	Messages    []*model.Message      `json:"messages"`
	DeadLetters []*model.DeadLetter   `json:"dead_letters"`
	Producers   []*model.Producer     `json:"producers"`
	Groups      []*model.ConsumerGroup `json:"groups"`
}

func (s *Service) ExportTopic(topicID string) (*TopicExport, error) {
	topic, err := s.store.GetTopic(topicID)
	if err != nil {
		return nil, err
	}
	export := &TopicExport{
		TopicID:   topicID,
		TopicName: topic.Name,
	}
	messages := s.store.GetMessagesByTopic(topicID)
	export.Messages = make([]*model.Message, len(messages))
	copy(export.Messages, messages)

	deadLetters := s.store.GetDeadLettersByTopic(topicID)
	export.DeadLetters = make([]*model.DeadLetter, len(deadLetters))
	copy(export.DeadLetters, deadLetters)

	allProducers := s.store.ListProducers()
	producers := make([]*model.Producer, 0)
	for _, p := range allProducers {
		if p.TopicID == topicID {
			producers = append(producers, p)
		}
	}
	export.Producers = producers

	allGroups := s.store.ListConsumerGroups()
	groups := make([]*model.ConsumerGroup, 0)
	for _, g := range allGroups {
		if g.TopicID == topicID {
			groups = append(groups, g)
		}
	}
	export.Groups = groups

	return export, nil
}
