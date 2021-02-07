package repositories

import (
	"errors"
	"flightAPI/server/models"
	"sort"
	"sync"
)

type MapItemRepository struct {
	lock        sync.Mutex
	items       map[int]models.Item
	primaryKeys []int
}

func NewMapItemRepository() *MapItemRepository {
	repository := new(MapItemRepository)
	repository.items = make(map[int]models.Item)
	repository.primaryKeys = []int{}
	return repository
}

func (itemRepository *MapItemRepository) AddItem(item models.Item) (*models.Item, error) {
	itemRepository.lock.Lock()
	defer itemRepository.lock.Unlock()
	_, ok := itemRepository.items[int(item.ID)]
	if !ok {
		itemRepository.items[int(item.ID)] = item
		i := sort.SearchInts(itemRepository.primaryKeys, int(item.ID))
		itemRepository.primaryKeys = append(itemRepository.primaryKeys, 0)
		copy(itemRepository.primaryKeys[i+1:], itemRepository.primaryKeys[i:])
		itemRepository.primaryKeys[i] = int(item.ID)
		return &item, nil
	} else {
		return nil, errors.New("already_exist")
	}
}

func (itemRepository *MapItemRepository) FindItem(ID int32) (*models.Item, error) {
	itemRepository.lock.Lock()
	defer itemRepository.lock.Unlock()
	item, ok := itemRepository.items[int(ID)]
	if ok {
		return &item, nil
	} else {
		return nil, errors.New("not_found")
	}
}

func (itemRepository *MapItemRepository) DeleteItem(ID int32) error {
	itemRepository.lock.Lock()
	defer itemRepository.lock.Unlock()
	index := sort.SearchInts(itemRepository.primaryKeys, int(ID))
	if index < len(itemRepository.primaryKeys) && itemRepository.primaryKeys[index] == int(ID) {
		delete(itemRepository.items, int(ID))
		copy(itemRepository.primaryKeys[index:], itemRepository.primaryKeys[index+1:])
		itemRepository.primaryKeys = itemRepository.primaryKeys[:len(itemRepository.primaryKeys)-1]
		return nil
	}
	return errors.New("not_found")
}

func (itemRepository *MapItemRepository) FindItems(from int32, limit int32) ([]*models.Item, error) {
	itemRepository.lock.Lock()
	defer itemRepository.lock.Unlock()
	itemsQty := len(itemRepository.primaryKeys)
	if int(from) >= itemsQty {
		return []*models.Item{}, nil
	} else {
		to := itemsQty
		if itemsQty > int(from)+int(limit) {
			to = int(from) + int(limit)
		}

		var items []*models.Item
		for _, id := range itemRepository.primaryKeys[from:to] {
			item := itemRepository.items[id]
			items = append(items, &item)
		}

		return items, nil
	}
}

func (itemRepository *MapItemRepository) UpdateItem(item models.Item) (*models.Item, error) {
	itemRepository.lock.Lock()
	defer itemRepository.lock.Unlock()
	_, ok := itemRepository.items[int(item.ID)]
	if ok {
		itemRepository.items[int(item.ID)] = item
		return &item, nil
	} else {
		return nil, errors.New("not_found")
	}
}
