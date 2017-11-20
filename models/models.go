package models

import "time"

// Command represents an abstract action
type Command string

const (
	// CommandAddShoppingItem represents adding into a shopping cart
	CommandAddShoppingItem = Command("ADD_SHOPPING_ITEM")
)

// UnfinishedCommand represents unfinished operations
// for multi step user interactions.
// For example, when a user wasnts to add an item into his shopping list,
// he needs to send the "add" command and then answer bot's question
// (send the name of item), so we need to remember what we have asked user for.
//
// Note: it will be fun, if an user starts opperation in a private chat
// and finishes in a group chat. So, probably, we should
// store unfinished commands per chatID
type UnfinishedCommand struct {
	Command   Command
	CreatedBy int // Telegram User ID
	CreatedAt *time.Time
}

// ShoppingItem represents an item in a shopping list
type ShoppingItem struct {
	ID        int64
	Name      string
	IsActive  bool // Indicates that the item is still active
	CreatedBy int  // Telegram User ID
	CreatedAt *time.Time
}
