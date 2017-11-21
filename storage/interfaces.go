package storage

import "github.com/m1kola/shipsterbot/models"

// TODO: Decide if it makes sense to move we move shoppingItemsMap and chatShoppingListsMap to models

// shoppingItemsMap stores shopping items by their id
type shoppingItemsMap map[int64]*models.ShoppingItem

// chatShoppingListsMap stores shopping lisets by chat id
type chatShoppingListsMap map[int64]*shoppingItemsMap

// DataStorageInterface represents a struct that handles storage logic
type DataStorageInterface interface {
	AddUnfinishedCommand(chatID int64, userID int, command models.Command)
	GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, bool)
	DeleteUnfinishedCommand(chatID int64, userID int)
	AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem)
	GetShoppingItems(chatID int64) (*shoppingItemsMap, bool)
	GetShoppingItem(chatID, itemID int64) (*models.ShoppingItem, bool)
	DeleteShoppingItem(chatID, itemID int64)
	DeleteAllShoppingItems(chatID int64)
}
