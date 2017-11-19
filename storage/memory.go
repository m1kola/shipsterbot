package storage

import (
	"time"

	"github.com/m1kola/telegram_shipsterbot/models"
)

// MemoryStorage implements the DataStorageInterface DataStorageInterface
// to store data in memory. Useful for prototyping and, potentially,
// for tests.
type MemoryStorage struct{}

// Note that it will be fun, if an user starts opperation in a private chat
// and finishes in a group chat
var unfinishedCommands = make(map[int]*models.UnfinishedCommand)
var items = make(map[int64][]*models.ShoppingItem)

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s MemoryStorage) AddUnfinishedCommand(UserID int, command models.Command) {
	now := time.Now()

	unfinishedCommands[UserID] = &models.UnfinishedCommand{
		Command:   command,
		CreatedBy: UserID,
		CreatedAt: &now}
}

// GetUnfinishedCommand returns an unfinished operaiont from the storage
func (s MemoryStorage) GetUnfinishedCommand(UserID int) (*models.UnfinishedCommand, bool) {
	item, ok := unfinishedCommands[UserID]
	return item, ok
}

// DeleteUnfinishedCommand deletes an unfinished operaiont from the storage
func (s MemoryStorage) DeleteUnfinishedCommand(UserID int) {
	delete(unfinishedCommands, UserID)
}

// AddShoppingItemIntoShoppingList adds a shoping item into a shipping list
// of a specific chat
func (s MemoryStorage) AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem) {
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	items[chatID] = append(items[chatID], item)
}

// GetShoppingItems returns a shopping list for a specific chat
func (s MemoryStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool) {
	shoppingList, ok := items[chatID]
	return shoppingList, ok
}

// GetShoppingItem returns a shopping item by id from a specific chat
func (s MemoryStorage) GetShoppingItem(chatID int64, itemID int64) (*models.ShoppingItem, bool) {
	itemID--

	shoppingList, ok := items[chatID]
	if ok && itemID >= 0 && itemID < int64(len(shoppingList)) {
		return shoppingList[itemID], true
	}
	return nil, false
}

// DeleteShoppingItem deletes a shipping item from a shipping lits
// for a specific chat
func (s MemoryStorage) DeleteShoppingItem(chatID int64, itemID int64) {
	itemID--

	shoppingList, ok := items[chatID]
	if ok {
		items[chatID] = append(shoppingList[:itemID], shoppingList[itemID+1:]...)
	}
}

// DeleteAllShoppingItems deletes all shopping items for a specific chat
func (s MemoryStorage) DeleteAllShoppingItems(chatID int64) {
	items[chatID] = []*models.ShoppingItem{}
}
