package types

// Available operations
import (
	"time"
)

type Operation string

const (
	// OperationAdd represents addition
	OperationAdd = Operation("ADD")
)

// UnfinishedOperation is a map of not finished operations
// TODO: This should, probably, be stored in DB for scalability and persistance
type UnfinishedOperation struct {
	Operation Operation
	CreatedAt *time.Time
}

// ShoppingItem represents an item in a shopping list
// TODO: This should, probably, also should be stored in DB for scalability
type ShoppingItem struct {
	Name      string
	IsActive  bool // Indicates that the item is still active
	CreatedBy int  // Telegram User ID
	CreatedAt *time.Time
}

type ShoppingList []*ShoppingItem
