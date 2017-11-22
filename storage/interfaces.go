package storage

import "github.com/m1kola/shipsterbot/models"

// DataStorageInterface represents a struct that handles storage logic
type DataStorageInterface interface {
	AddUnfinishedCommand(command models.UnfinishedCommand)
	GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, bool)
	DeleteUnfinishedCommand(chatID int64, userID int)

	AddShoppingItemIntoShoppingList(item models.ShoppingItem)
	// TODO: Decide if we can remoe chatID. It doesn't seem useful
	GetShoppingItem(chatID, itemID int64) (*models.ShoppingItem, bool)
	GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool)
	// TODO: Decide if we can remoe chatID. It doesn't seem useful
	DeleteShoppingItem(chatID, itemID int64)
	DeleteAllShoppingItems(chatID int64)
}
