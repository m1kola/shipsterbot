package models

import "time"

// Command represents an abstract action
type Command string

const (
	// CommandAddShoppingItem represents adding into a shopping cart
	CommandAddShoppingItem = Command("ADD_SHOPPING_ITEM")
	// TODO: Decide if we need to have all bot commands here
)

// UnfinishedCommand is a map of not finished operations
// for multi step user interactions.
// For example, when a user wasnts to add an item into his shopping list,
// he needs to send the "add" command and then answer bot's question
// (send the name of item), so we need to remember what we have asked user for.
type UnfinishedCommand struct {
	Command   Command
	CreatedBy int // Telegram User ID
	CreatedAt *time.Time
}

// ShoppingItem represents an item in a shopping list
type ShoppingItem struct {
	Name      string
	IsActive  bool // Indicates that the item is still active
	CreatedBy int  // Telegram User ID
	CreatedAt *time.Time
}
