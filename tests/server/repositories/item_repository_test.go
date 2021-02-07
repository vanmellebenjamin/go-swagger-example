package repositories

import (
	"flightAPI/server/models"
	"flightAPI/server/repositories"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestInsertGet(t *testing.T) {
	// given
	var itemRepository repositories.ItemRepository = repositories.NewMapItemRepository()
	item := new(models.Item)
	item.ID = 5
	completed := false
	item.Completed = &completed
	desc := "Task 1"
	item.Description = &desc

	// when
	_, err := itemRepository.AddItem(*item)

	// then
	if err != nil {
		t.Errorf("error append during insert %s", err)
	}
	itemResult, errFind := itemRepository.FindItem(5)
	if errFind != nil {
		t.Errorf("error append during get %s", errFind)
	}
	if itemResult.ID != item.ID {
		t.Errorf("expected %v, actual %v", item.ID, itemResult.ID)
	}
}

func TestInsertGetNotFound(t *testing.T) {
	// given
	var itemRepository repositories.ItemRepository = repositories.NewMapItemRepository()
	item := new(models.Item)
	item.ID = 5
	completed := false
	item.Completed = &completed
	desc := "Task 1"
	item.Description = &desc

	// when
	_, err := itemRepository.AddItem(*item)

	// then
	if err != nil {
		t.Errorf("error append during insert %s", err)
	}
	_, errFind := itemRepository.FindItem(6)
	if errFind == nil {
		t.Errorf("error append during get %s", errFind)
	}
	if errFind.Error() != "not_found" {
		assert.Equal(t, "not_found", errFind)
	}
}

func TestInsertAndGetList(t *testing.T) {
	// given
	var itemRepository repositories.ItemRepository = repositories.NewMapItemRepository()
	completed := false
	itemRepository.AddItem(models.Item{ID: 1, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 6, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 3, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 2, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 5, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 4, Completed: &completed, Description: nil})

	// when
	it, err := itemRepository.FindItems(0, 10)
	if err != nil {
		t.Errorf("error append during get %s", err)
	}

	items := it
	assert.Equal(t, int32(1), items[0].ID)
	assert.Equal(t, int32(2), items[1].ID)
	assert.Equal(t, int32(3), items[2].ID)
	assert.Equal(t, int32(4), items[3].ID)
	assert.Equal(t, int32(5), items[4].ID)
	assert.Equal(t, int32(6), items[5].ID)
}

func TestInsertDeleteAndGetList(t *testing.T) {
	// given
	var itemRepository repositories.ItemRepository = repositories.NewMapItemRepository()
	completed := false
	itemRepository.AddItem(models.Item{ID: 1, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 6, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 3, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 2, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 5, Completed: &completed, Description: nil})
	itemRepository.AddItem(models.Item{ID: 4, Completed: &completed, Description: nil})
	itemRepository.DeleteItem(4)
	itemRepository.DeleteItem(2)

	// when
	it, err := itemRepository.FindItems(0, 10)
	if err != nil {
		t.Errorf("error append during get %s", err)
	}

	items := it
	assert.Equal(t, int32(1), items[0].ID)
	assert.Equal(t, int32(3), items[1].ID)
	assert.Equal(t, int32(5), items[2].ID)
	assert.Equal(t, int32(6), items[3].ID)
}

func TestInsertUpdate(t *testing.T) {
	// given
	var itemRepository repositories.ItemRepository = repositories.NewMapItemRepository()
	completed := false
	itemRepository.AddItem(models.Item{ID: 1, Completed: &completed, Description: nil})

	// when
	completed = true
	itemRepository.UpdateItem(models.Item{ID: 1, Completed: &completed, Description: nil})

	// then
	item, err := itemRepository.FindItem(1)
	if err != nil {
		t.Errorf("error append during get %s", err)
	}

	assert.Equal(t, true, *item.Completed)
}

func TestConcurrentInsert(t *testing.T) {
	// given
	workers := 10
	qtyPerWorker := 1_000
	var itemRepository repositories.ItemRepository = repositories.NewMapItemRepository()
	var wg sync.WaitGroup

	// when
	for workerNumber := 0; workerNumber < workers; workerNumber++ {
		wg.Add(1)
		go insertItems(workerNumber, qtyPerWorker, itemRepository, &wg)
	}
	wg.Wait()

	// then
	for i := 0; i < workers*qtyPerWorker; i++ {
		item, err := itemRepository.FindItem(int32(i))
		if err != nil {
			t.Errorf("error append during get %v %s", i, err)
		}
		assert.Equal(t, int32(i), item.ID)
	}
}

func insertItems(workerNumber int, qty int, itemRepository repositories.ItemRepository, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(200) * time.Millisecond)
	completed := false
	for i := 0; i < qty; i++ {
		itemRepository.AddItem(models.Item{ID: int32(workerNumber*qty + i), Completed: &completed, Description: nil})
	}
}
