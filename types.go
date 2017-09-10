package main

// Available operations
import (
	"time"
)

const (
	operationAdd = "ADD"
)

// Map of not finished operations
// TODO: This should, probably stored in DB for scalability
type unfinishedOperation struct {
	Operation string
	Time      time.Time
}
type unfinishedOperationsByUserID map[int]*unfinishedOperation
