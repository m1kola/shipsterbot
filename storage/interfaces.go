package storage

import "github.com/m1kola/telegram_shipsterbot/models"

// DataStorageInterface represents a struct that handles storage logic
type DataStorageInterface interface {
	AddUnfinishedCommand(UserID int, command models.Command)
	GetUnfinishedCommand(UserID int) (*models.UnfinishedCommand, bool)
	DeleteUnfinishedCommand(UserID int)
	AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem)
	GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool)
	GetShoppingItem(chatID int64, itemID int64) (*models.ShoppingItem, bool)
	DeleteShoppingItem(chatID int64, itemID int64)
	DeleteAllShoppingItems(chatID int64)
}
