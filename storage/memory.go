package storage

import (
	"time"

	"github.com/m1kola/shipsterbot/models"
)

// shoppingItemsMap stores shopping items by their id
type shoppingItemsMap map[int64]*models.ShoppingItem

// chatShoppingListsMap stores shopping lisets by chat id
type chatShoppingListsMap map[int64]*shoppingItemsMap

// MemoryStorage implements the DataStorageInterface
// to store data in memory. Useful for prototyping and, potentially,
// for tests.
type MemoryStorage struct {
	unfinishedCommands map[int]*models.UnfinishedCommand
	latestItemID       int64
	items              chatShoppingListsMap
}

// NewMemoryStorage initialises a new MemoryStorage instance
func NewMemoryStorage() *MemoryStorage {
	storage := MemoryStorage{}

	storage.unfinishedCommands = make(map[int]*models.UnfinishedCommand)
	storage.items = make(chatShoppingListsMap)

	return &storage
}

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s *MemoryStorage) AddUnfinishedCommand(chatID int64, userID int, command models.Command) {
	now := time.Now()

	s.unfinishedCommands[userID] = &models.UnfinishedCommand{
		Command:   command,
		ChatID:    chatID,
		CreatedBy: userID,
		CreatedAt: &now}
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
func (s *MemoryStorage) AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem) {
	// Set an ID
	s.latestItemID++
	item.ID = s.latestItemID

	// Set CreatedAt if not present
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	// Create a shopping list, if not present
	_, ok := s.items[chatID]
	if !ok {
		newshoppingList := make(shoppingItemsMap)
		s.items[chatID] = &newshoppingList
	}

	// Insert and item
	(*s.items[chatID])[item.ID] = item
}

// GetShoppingItems returns a shopping list for a specific chat
func (s *MemoryStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool) {
	shoppingList, ok := s.items[chatID]
	if !ok {
		return nil, false
	}

	itemsList := make([]*models.ShoppingItem, 0, len(*shoppingList))
	for _, item := range *shoppingList {
		itemsList = append(itemsList, item)
	}

	return itemsList, true
}

// GetShoppingItem returns a shopping item by id from a specific chat
func (s *MemoryStorage) GetShoppingItem(chatID, itemID int64) (*models.ShoppingItem, bool) {
	shoppingList, ok := s.items[chatID]

	if ok && itemID > 0 && itemID <= s.latestItemID {
		item, ok := (*shoppingList)[itemID]
		return item, ok
	}
	return nil, false
}

// DeleteShoppingItem deletes a shipping item from a shipping lits
// for a specific chat
func (s *MemoryStorage) DeleteShoppingItem(chatID, itemID int64) {
	shoppingList, ok := s.items[chatID]
	if ok {
		delete(*shoppingList, itemID)
	}
}

// DeleteAllShoppingItems deletes all shopping items for a specific chat
func (s *MemoryStorage) DeleteAllShoppingItems(chatID int64) {
	s.items[chatID] = &shoppingItemsMap{}
}
