package repositories

import (
	"flightAPI/server/models"
)

type ItemRepository interface {
	AddItem(item models.Item) (*models.Item, error)
	FindItem(ID int32)  (*models.Item, error)
	DeleteItem(ID int32) error
	FindItems(from int32, limit int32)  ([]*models.Item, error)
	UpdateItem(item models.Item) (*models.Item, error)
}
