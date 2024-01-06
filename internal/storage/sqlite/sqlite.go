package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	op := "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS events (
    						event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    						name VARCHAR(100) NOT NULL UNIQUE,
    						tickets INT NOT NULL,
    						price INT NOT NULL,
    						start_date DATETIME NOT NULL);`)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS users_orders (
    						user_order_id INTEGER PRIMARY KEY AUTOINCREMENT,
    						name VARCHAR(100) NOT NULL,
    						event_id INT NOT NULL,
							tickets INT NOT NULL,
							order_date DATETIME NOT NULL,
							FOREIGN KEY (event_id)
								REFERENCES events (event_id)
									ON DELETE CASCADE
									ON UPDATE NO ACTION);`)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}
	defer stmt.Close()

	return &Storage{db: db}, nil
}

func (storage *Storage) SaveUserTickets(name string, eventId, ticketsQty int) error {
	op := "storage.sqlite.AddUserTickets"

	stmt, err := storage.db.Prepare(`UPDATE events
									SET tickets = tickets - ?`)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec(ticketsQty)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	stmt, err = storage.db.Prepare(`INSERT INTO users_orders(name, event_id, tickets, order_date) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec(name, eventId, ticketsQty, time.Now())
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (storage *Storage) SaveEvent(name string, tickets, price int, startDate time.Time) error {
	op := "storage.sqlite.AddEvent"

	stmt, err := storage.db.Prepare(`INSERT INTO events (name, tickets, price, start_date) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec(name, tickets, price, startDate)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (storage *Storage) DeleteOrder(orderId int) error {
	op := "storage.sqlite.DeleteOrder"

	row := storage.db.QueryRow(`SELECT tickets FROM users_orders WHERE user_order_id = ?`, orderId)

	var tickets int
	err := row.Scan(&tickets)

	stmt, err := storage.db.Prepare(`DELETE FROM users_orders WHERE user_order_id = ?`)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec(orderId)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	stmt, err = storage.db.Prepare(`UPDATE events
									SET tickets = tickets + ?`)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	_, err = stmt.Exec(tickets)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}
