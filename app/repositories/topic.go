package repositories

import (
	"go-shop-v2/pkg/repository"
)

type TopicRep struct {
	repository.IRepository
}

func NewTopicRep(rep repository.IRepository) *TopicRep {
	return &TopicRep{rep}
}
