package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Item struct {
	Login string  `json:"login"`
	Money int     `json:"money"`
	Score float64 `json:"score"`
}

type Storage struct {
	Db *sqlx.DB
}

func NewStorage() (*Storage, error) {
	db, err := sqlx.Open("postgres", "host=localhost port=5432 user=postgres password=12345 dbname=postgres sslmode=disable")

	if err != nil {
		return nil, err
		fmt.Print("Что-то не c аргументами подключения к БД")
	}

	if err = db.Ping(); err != nil {
		return nil, err
		fmt.Print("Что-то не так с подключением к БД (2)")
	}

	return &Storage{db}, nil
}

func (s *Storage) CreateTables() error {
	var createTablesQuery string = `
CREATE TABLE IF NOT EXISTS items (
    money INT NOT NULL,
    login TEXT NOT NULL,
    score FLOAT NOT NULL
);`
	_, err := s.Db.Exec(createTablesQuery)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) PostStorageFunc(item *Item) error {
	query := `INSERT INTO items (login, money, score) VALUES ($1, $2, $3)`
	_, err := s.Db.Exec(query, item.Login, item.Money, item.Score)
	if err != nil {
		fmt.Println("При попытке записать произошла ошибка:", err)
		return err
	}
	return nil
}

func (s *Storage) GetAllItems() ([]Item, error) {
	err2 := s.Db.Ping()
	if err2 != nil {
		fmt.Println("проблема с подключением ")
	}
	rows, err := s.Db.Query("SELECT login, money, score FROM items")
	fmt.Println(rows)
	defer rows.Close()
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса SELECT:", err)
		return nil, err
	}
	var items []Item
	for rows.Next() {
		item := Item{}
		err := rows.Scan(&item.Login, &item.Money, &item.Score)
		if err != nil {
			fmt.Println(items)
			return nil, err
		}
		items = append(items, item)
	}
	fmt.Println(items)
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
