package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int64) error
	UpdateAccount(*Account) error
	GetAccountById(int64) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id serial PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    number serial,
    balance NUMERIC(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		insert into account
		(first_name, last_name, number, balance, created_at)
		values
		($1, $2, $3, $4, $5)
	`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int64) error {
	_, err := s.db.Query("delete from account where id = $1", id)

	return err
}

func (s *PostgresStore) GetAccountById(id int64) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return s.ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)

	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	return account, err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`select * from account;`)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		acc, err := s.ScanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}
