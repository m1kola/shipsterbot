package storage

import (
	"database/sql"

	"github.com/m1kola/shipsterbot/internal/pkg/models"
)

// SQLStorage implements the DataStorageInterface
// to store data persistently in an SQL RDBMS
type SQLStorage struct {
	db *sql.DB
}

// NewSQLStorage initialises a new NewSQLStorage instance
func NewSQLStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{db: db}
}

// AddUnfinishedCommand inserts an unfinished operaiont into the storage
func (s *SQLStorage) AddUnfinishedCommand(command models.UnfinishedCommand) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete a previous unfinshed command (if any) in a transaction
	_, err = tx.Exec(
		`DELETE FROM
			unfinished_commands
		WHERE
			chat_id = $1
			AND created_by = $2`,
		command.ChatID, command.CreatedBy)
	if err != nil {
		return err
	}

	// Add a new unfinshed command
	_, err = tx.Exec(
		`INSERT INTO
			unfinished_commands(command, chat_id, created_by)
		VALUES ($1, $2, $3)`,
		command.Command, command.ChatID, command.CreatedBy)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// GetUnfinishedCommand returns an unfinished operaiont from the storage
func (s *SQLStorage) GetUnfinishedCommand(chatID int64, userID int) (*models.UnfinishedCommand, error) {
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

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &command, err
}

// DeleteUnfinishedCommand deletes an unfinished operaiont from the storage
func (s *SQLStorage) DeleteUnfinishedCommand(chatID int64, userID int) error {
	_, err := s.db.Exec(
		`DELETE FROM
			unfinished_commands
		WHERE
			chat_id = $1
			AND created_by = $2`,
		chatID, userID)

	return err
}

// AddShoppingItemIntoShoppingList adds a shoping item into a shipping list
// of a specific chat
func (s *SQLStorage) AddShoppingItemIntoShoppingList(item models.ShoppingItem) error {
	_, err := s.db.Exec(
		`INSERT INTO
			shopping_items (name, chat_id, created_by)
		VALUES ($1, $2, $3)`,
		item.Name, item.ChatID, item.CreatedBy)

	return err
}

// GetShoppingItems returns a shopping list for a specific chat
func (s *SQLStorage) GetShoppingItems(chatID int64) ([]*models.ShoppingItem, error) {
	var itemsList []*models.ShoppingItem

	rows, err := s.db.Query(
		`SELECT
			id, name, chat_id, created_by, created_at
		FROM shopping_items
		WHERE
			chat_id = $1`,
		chatID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := models.ShoppingItem{}
		err = rows.Scan(
			&item.ID, &item.Name, &item.ChatID,
			&item.CreatedBy, &item.CreatedAt)

		if err != nil {
			return nil, err
		}

		itemsList = append(itemsList, &item)

	}

	err = rows.Err()
	return itemsList, err
}

// GetShoppingItem returns a shopping item by id from a specific chat
func (s *SQLStorage) GetShoppingItem(itemID int64) (*models.ShoppingItem, error) {
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

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &item, err
}

// DeleteShoppingItem deletes a shipping item from a shipping lits
// for a specific chat
func (s *SQLStorage) DeleteShoppingItem(itemID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM
			shopping_items
		WHERE
			id = $1`,
		itemID)

	return err
}

// DeleteAllShoppingItems deletes all shopping items for a specific chat
func (s *SQLStorage) DeleteAllShoppingItems(chatID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM
			shopping_items
		WHERE
			chat_id = $1`,
		chatID)

	return err
}
