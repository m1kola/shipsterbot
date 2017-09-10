package main

// Available operations
import (
	"time"
)

const (
	operationAdd = "ADD"
)

// Map of not finished operations
// TODO: This should, probably, be stored in DB for scalability and persistance
type unfinishedOperation struct {
	Operation string
	Time      time.Time
}

// Note that it will be fun, if an user starts opperation in a private chat
// and finishes in a group chat
type unfinishedOperationsByUserID map[int]*unfinishedOperation

// TODO: This should, probably, also should be stored in DB for scalability
type shoppingItem struct {
	Name      string
	IsActive  bool // Indicates that the item is still active
	CreatedBy int  // Telegram User ID
	CreatedAt time.Time
}
type shoppingItemsByChatID map[int64][]*shoppingItem
