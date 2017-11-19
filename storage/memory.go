package storage

import (
	"time"

	"github.com/m1kola/telegram_shipsterbot/models"
)

// MemoryStorage implements the DataStorageInterface DataStorageInterface
// to store data in memory. Useful for prototyping and, potentially,
// for tests.
type MemoryStorage struct {
	unfinishedCommands map[int]*models.UnfinishedCommand
	items              map[int64][]*models.ShoppingItem
}

// NewMemoryStorage initialises a new MemoryStorage instance
func NewMemoryStorage() *MemoryStorage {
	storage := MemoryStorage{}

	// Note that it will be fun, if an user starts opperation in a private chat
	// and finishes in a group chat. So, probably, we should
	// store unfinished commands per chatID
	storage.unfinishedCommands = make(map[int]*models.UnfinishedCommand)
	storage.items = make(map[int64][]*models.ShoppingItem)

	return &storage
}

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s MemoryStorage) AddUnfinishedCommand(UserID int, command models.Command) {
	now := time.Now()

	s.unfinishedCommands[UserID] = &models.UnfinishedCommand{
		Command:   command,
		CreatedBy: UserID,
		CreatedAt: &now}
}

// GetUnfinishedCommand returns an unfinished operaiont from the storage
func (s MemoryStorage) GetUnfinishedCommand(UserID int) (*models.UnfinishedCommand, bool) {
	item, ok := s.unfinishedCommands[UserID]
	return item, ok
}

// DeleteUnfinishedCommand deletes an unfinished operaiont from the storage
func (s MemoryStorage) DeleteUnfinishedCommand(UserID int) {
	delete(s.unfinishedCommands, UserID)
}

// AddShoppingItemIntoShoppingList adds a shoping item into a shipping list
// of a specific chat
func (s MemoryStorage) AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem) {
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	s.items[chatID] = append(s.items[chatID], item)
}

// GetShoppingItems returns a shopping list for a specific chat
func (s MemoryStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool) {
	shoppingList, ok := s.items[chatID]
	return shoppingList, ok
}

// GetShoppingItem returns a shopping item by id from a specific chat
func (s MemoryStorage) GetShoppingItem(chatID int64, itemID int64) (*models.ShoppingItem, bool) {
	itemID--

	shoppingList, ok := s.items[chatID]
	if ok && itemID >= 0 && itemID < int64(len(shoppingList)) {
		return shoppingList[itemID], true
	}
	return nil, false
}

// DeleteShoppingItem deletes a shipping item from a shipping lits
// for a specific chat
func (s MemoryStorage) DeleteShoppingItem(chatID int64, itemID int64) {
	itemID--

	shoppingList, ok := s.items[chatID]
	if ok {
		s.items[chatID] = append(shoppingList[:itemID], shoppingList[itemID+1:]...)
	}
}

// DeleteAllShoppingItems deletes all shopping items for a specific chat
func (s MemoryStorage) DeleteAllShoppingItems(chatID int64) {
	s.items[chatID] = []*models.ShoppingItem{}
}
