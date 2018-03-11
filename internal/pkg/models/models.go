package models

import "time"

// UnfinishedCommand represents unfinished operations
// for multi step user interactions.
// For example, when a user wasnts to add an item into his shopping list,
// he needs to send the "add" command and then answer bot's question
// (send the name of item), so we need to remember what we have asked user for.
type UnfinishedCommand struct {
	Command   string
	ChatID    int64
	CreatedBy int
	CreatedAt *time.Time
}

// ShoppingItem represents an item in a shopping list
type ShoppingItem struct {
	ID        int64
	Name      string
	ChatID    int64
	CreatedBy int
	CreatedAt *time.Time
}
