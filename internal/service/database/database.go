package database

import (
	"github.com/mnogokotin/chat-gpt-test/internal/model"
	"github.com/mnogokotin/golang-packages/database/postgres"
)

func GetEventModelsWithGreaterId(pg *postgres.Postgres, id string) []model.Event {
	var eventModels []model.Event
	pg.Db.Where("id > ?", id).Find(&eventModels)
	return eventModels
}
