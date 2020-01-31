package repositories

import (
	"go-shop-v2/app/models"
	"go-shop-v2/pkg/db/mongodb"
)

type TopicRep struct {
	*mongoRep
}

func NewTopicRep(con *mongodb.Connection) *TopicRep {
	return &TopicRep{NewBasicMongoRepositoryByDefault(&models.Topic{}, con)}
}
