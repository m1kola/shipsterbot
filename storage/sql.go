package storage

import (
	"database/sql"
	"log"

	"github.com/m1kola/shipsterbot/models"
)

// SQLStorage implements the DataStorageInterface
// to store data persistently in an SQL RDBMS
type SQLStorage struct {
	db *sql.DB
}

// NewSQLStorage initialises a new NewSQLStorage instance
func NewSQLStorage(db *sql.DB) *SQLStorage {
	storage := SQLStorage{db: db}
	return &storage
}

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s *SQLStorage) AddUnfinishedCommand(chatID int64, userID int, command models.Command) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO
			unfinished_commands(command, chat_id, created_by)
		VALUES ($1, $2, $3)`,
		command, chatID, userID)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// GetUnfinishedCommand returns an unfinished operaiont from the storage
func (s *SQLStorage) GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, bool) {
	command := models.UnfinishedCommand{}
	row := s.db.QueryRow(
		`SELECT
			command, chat_id, created_by, created_at
		FROM unfinished_commands
		WHERE
			chat_id = $1
			AND created_by = $2`,
		chatID, userID)

	err := row.Scan(
		&command.Command,
		&command.ChatID,
		&command.CreatedBy,
		&command.CreatedAt)

	return &command, err == nil
}

// DeleteUnfinishedCommand deletes an unfinished operaiont from the storage
func (s *SQLStorage) DeleteUnfinishedCommand(chatID int64, userID int) {
	_, err := s.db.Exec(
		`DELETE FROM
			unfinished_commands
		WHERE
			chat_id = $1
			AND created_by = $2`,
		chatID, userID)

	if err != nil {
		log.Fatal(err)
	}
}

// AddShoppingItemIntoShoppingList adds a shoping item into a shipping list
// of a specific chat
func (s *SQLStorage) AddShoppingItemIntoShoppingList(chatID int64, item *models.ShoppingItem) {
	_, err := s.db.Exec(
		`INSERT INTO
			shopping_items (name, chat_id, created_by)
		VALUES ($1, $2, $3)`,
		item.Name, item.ChatID, item.CreatedBy)

	if err != nil {
		log.Fatal(err)
	}
}

// GetShoppingItems returns a shopping list for a specific chat
func (s *SQLStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, bool) {
	var itemsList []*models.ShoppingItem

	rows, err := s.db.Query(
		`SELECT
			id, name, chat_id, created_by, created_at
		FROM shopping_items
		WHERE
			chat_id = $1`,
		chatID)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		item := models.ShoppingItem{}
		err = rows.Scan(
			&item.ID, &item.Name, &item.ChatID,
			&item.CreatedBy, &item.CreatedAt)

		if err != nil {
			log.Fatal(err)
		}

		itemsList = append(itemsList, &item)

	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return itemsList, true
}

// GetShoppingItem returns a shopping item by id from a specific chat
func (s *SQLStorage) GetShoppingItem(chatID int64, itemID int64) (*models.ShoppingItem, bool) {
	item := models.ShoppingItem{}
	row := s.db.QueryRow(
		`SELECT
			id, name, chat_id, created_by, created_at
		FROM shopping_items
		WHERE
			id = $1`,
		itemID)

	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.ChatID,
		&item.CreatedBy,
		&item.CreatedAt)

	return &item, err == nil
}

// DeleteShoppingItem deletes a shipping item from a shipping lits
// for a specific chat
func (s *SQLStorage) DeleteShoppingItem(chatID int64, itemID int64) {
	_, err := s.db.Exec(
		`DELETE FROM
			shopping_items
		WHERE
			id = $1`,
		itemID)

	if err != nil {
		log.Fatal(err)
	}
}

// DeleteAllShoppingItems deletes all shopping items for a specific chat
func (s *SQLStorage) DeleteAllShoppingItems(chatID int64) {
	_, err := s.db.Exec(
		`DELETE FROM
			shopping_items
		WHERE
			chat_id = $1`,
		chatID)

	if err != nil {
		log.Fatal(err)
	}
}
