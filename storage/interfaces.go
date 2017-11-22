package storage

import "github.com/m1kola/shipsterbot/models"

// DataStorageInterface represents a struct that handles storage logic
// TODO: Check if makes sense to pass pointers
type DataStorageInterface interface {
	// TODO: why not to accept UnfinishedCommand instance?
	AddUnfinishedCommand(chatID int64, userID int, command models.Command)
	GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, bool)
	DeleteUnfinishedCommand(chatID int64, userID int)
	// TODO: Decide if it makes any sense to pass chatID (we have item.ChatID now)
	AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem)
	GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool)
	GetShoppingItem(chatID, itemID int64) (*models.ShoppingItem, bool)
	DeleteShoppingItem(chatID, itemID int64)
	DeleteAllShoppingItems(chatID int64)
}
