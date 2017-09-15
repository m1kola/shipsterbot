package storage

import (
	"fmt"
	"time"

	"github.com/m1kola/telegram_shipsterbot/types"
)

// TODO: Use RDBMS for storing data

// Note that it will be fun, if an user starts opperation in a private chat
// and finishes in a group chat
var unfinishedOperations = make(map[int]*types.UnfinishedOperation)
var items = make(map[int64]types.ShoppingList)

// AddUnfinishedOperation inserts an unfinished operaiont into the database
func AddUnfinishedOperation(UserID int, operation types.Operation) {
	now := time.Now()

	unfinishedOperations[UserID] = &types.UnfinishedOperation{
		Operation: operation,
		CreatedAt: &now}
}

func GetUnfinishedOperation(UserID int) (*types.UnfinishedOperation, bool) {
	item, ok := unfinishedOperations[UserID]
	return item, ok
}

func DeleteUnfinishedOperation(UserID int) {
	delete(unfinishedOperations, UserID)
}

func AddShoppingItemIntoShoppingList(chatID int64, item *types.ShoppingItem) {
	if item.CreatedAt == nil {
		now := time.Now()
		item.CreatedAt = &now
	}

	fmt.Printf("%d-%d-%d\n", item.CreatedAt.Day(), item.CreatedAt.Month(), item.CreatedAt.Year())

	items[chatID] = append(items[chatID], item)
}

func GetShoppingItems(chatID int64) (*types.ShoppingList, bool) {
	shoppingList, ok := items[chatID]
	return &shoppingList, ok
}
