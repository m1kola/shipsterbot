package storage

import (
	"time"

	"github.com/m1kola/shipsterbot/models"
)

// shoppingItemsMap stores shopping items by their id
type shoppingItemsMap map[int64]*models.ShoppingItem

// MemoryStorage implements the DataStorageInterface
// to store data in memory. Useful for prototyping and, potentially,
// for tests.
type MemoryStorage struct {
	unfinishedCommands map[int]*models.UnfinishedCommand
	latestItemID       int64
	items              shoppingItemsMap
}

// NewMemoryStorage initialises a new MemoryStorage instance
func NewMemoryStorage() *MemoryStorage {
	storage := MemoryStorage{}

	storage.unfinishedCommands = make(map[int]*models.UnfinishedCommand)
	storage.items = make(shoppingItemsMap)

	return &storage
}

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s *MemoryStorage) AddUnfinishedCommand(command models.UnfinishedCommand) {
	// Set CreatedAt if not present
	if command.CreatedAt == nil {
		now := time.Now()
		command.CreatedAt = &now
	}

	// Delete a previous unfinshed command (if any)
	s.DeleteUnfinishedCommand(command.ChatID, command.CreatedBy)

	// Add a new unfinshed command
	s.unfinishedCommands[command.CreatedBy] = &command
}

// GetUnfinishedCommand returns an unfinished operaiont from the storage
func (s *MemoryStorage) GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, bool) {
	item, ok := s.unfinishedCommands[userID]
	return item, ok
}

// DeleteUnfinishedCommand deletes an unfinished operaiont from the storage
func (s *MemoryStorage) DeleteUnfinishedCommand(chatID int64, userID int) {
	delete(s.unfinishedCommands, userID)
}

// AddShoppingItemIntoShoppingList adds a shoping item into a shipping list
// of a specific chat
func (s *MemoryStorage) AddShoppingItemIntoShoppingList(item models.ShoppingItem) {
	// Set an ID
	s.latestItemID++
	item.ID = s.latestItemID

	// Set CreatedAt if not present
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	// Insert an item
	s.items[item.ID] = &item
}

// GetShoppingItems returns a shopping list for a specific chat
func (s *MemoryStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool) {
	var itemsList []*models.ShoppingItem
	for _, item := range s.items {
		if item.ChatID == chatID {
			itemsList = append(itemsList, item)
		}
	}

	if len(itemsList) == 0 {
		return nil, false
	}

	return itemsList, true
}

// GetShoppingItem returns a shopping item by id from a specific chat
func (s *MemoryStorage) GetShoppingItem(itemID int64) (*models.ShoppingItem, bool) {
	item, ok := s.items[itemID]
	return item, ok
}

// DeleteShoppingItem deletes a shipping item from a shipping lits
// for a specific chat
func (s *MemoryStorage) DeleteShoppingItem(itemID int64) {
	delete(s.items, itemID)
}

// DeleteAllShoppingItems deletes all shopping items for a specific chat
func (s *MemoryStorage) DeleteAllShoppingItems(chatID int64) {
	for _, item := range s.items {
		if item.ChatID == chatID {
			s.DeleteShoppingItem(item.ID)
		}
	}
}
