package storage

import (
	"time"

	"github.com/m1kola/telegram_shipsterbot/models"
)

type DatabaseStorage struct{}

// TODO: Use RDBMS for storing data

// Note that it will be fun, if an user starts opperation in a private chat
// and finishes in a group chat
var unfinishedCommands = make(map[int]*models.UnfinishedCommand)
var items = make(map[int64][]*models.ShoppingItem)

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s DatabaseStorage) AddUnfinishedCommand(UserID int, command models.Command) {
	now := time.Now()

	unfinishedCommands[UserID] = &models.UnfinishedCommand{
		Command:   command,
		CreatedBy: UserID,
		CreatedAt: &now}
}

func (s DatabaseStorage) GetUnfinishedCommand(UserID int) (*models.UnfinishedCommand, bool) {
	item, ok := unfinishedCommands[UserID]
	return item, ok
}

func (s DatabaseStorage) DeleteUnfinishedCommand(UserID int) {
	delete(unfinishedCommands, UserID)
}

func (s DatabaseStorage) AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem) {
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	items[chatID] = append(items[chatID], item)
}

func (s DatabaseStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool) {
	shoppingList, ok := items[chatID]
	return shoppingList, ok
}

func (s DatabaseStorage) GetShoppingItem(chatID int64, itemID int64) (*models.ShoppingItem, bool) {
	itemID--

	shoppingList, ok := items[chatID]
	if ok && itemID >= 0 && itemID < int64(len(shoppingList)) {
		return shoppingList[itemID], true
	}
	return nil, false
}

func (s DatabaseStorage) DeleteShoppingItem(chatID int64, itemID int64) {
	itemID--

	shoppingList, ok := items[chatID]
	if ok {
		items[chatID] = append(shoppingList[:itemID], shoppingList[itemID+1:]...)
	}
}

func (s DatabaseStorage) DeleteAllShoppingItems(chatID int64) {
	items[chatID] = []*models.ShoppingItem{}
}
