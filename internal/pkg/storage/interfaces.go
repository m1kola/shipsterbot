package storage

import "github.com/m1kola/shipsterbot/internal/pkg/models"

// Generates mocks for tests
//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

// DataStorageInterface represents a struct that handles storage logic
type DataStorageInterface interface {
	AddUnfinishedCommand(command models.UnfinishedCommand) error
	GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, error)
	DeleteUnfinishedCommand(chatID int64, userID int) error

	AddShoppingItemIntoShoppingList(item models.ShoppingItem) error
	GetShoppingItem(itemID int64) (*models.ShoppingItem, error)
	DeleteShoppingItem(itemID int64) error
	GetShoppingItems(chatID int64) ([]*models.ShoppingItem, error)
	DeleteAllShoppingItems(chatID int64) error
}
