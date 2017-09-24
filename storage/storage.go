package storage

import (
	"time"

	"github.com/m1kola/telegram_shipsterbot/types"
)

// TODO: Use RDBMS for storing data

// Note that it will be fun, if an user starts opperation in a private chat
// and finishes in a group chat
var unfinishedCommands = make(map[int]*types.UnfinishedCommand)
var items = make(map[int64][]*types.ShoppingItem)

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func AddUnfinishedCommand(UserID int, command types.Command) {
	now := time.Now()

	unfinishedCommands[UserID] = &types.UnfinishedCommand{
		Command:   command,
		CreatedBy: UserID,
		CreatedAt: &now}
}

func GetUnfinishedCommand(UserID int) (*types.UnfinishedCommand, bool) {
	item, ok := unfinishedCommands[UserID]
	return item, ok
}

func DeleteUnfinishedCommand(UserID int) {
	delete(unfinishedCommands, UserID)
}

func AddShoppingItemIntoShoppingList(chatID int64, item *types.ShoppingItem) {
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	items[chatID] = append(items[chatID], item)
}

func GetShoppingItems(chatID int64) ([]*types.ShoppingItem, bool) {
	shoppingList, ok := items[chatID]
	return shoppingList, ok
}

func GetShoppingItem(chatID int64, itemID int64) (*types.ShoppingItem, bool) {
	itemID--

	shoppingList, ok := items[chatID]
	if ok && itemID >= 0 && itemID < int64(len(shoppingList)) {
		return shoppingList[itemID], true
	}
	return nil, false
}

func DeleteShoppingItem(chatID int64, itemID int64) {
	itemID--

	shoppingList, ok := items[chatID]
	if ok {
		items[chatID] = append(shoppingList[:itemID], shoppingList[itemID+1:]...)
	}
}
