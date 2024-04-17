package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccounts() ([]*Account, error)
	CreateAccount(*Account) error
	DeleteAccount(int64) error
	UpdateAccount(*Account) error
	GetAccountByID(int64) (*Account, error)
	Transfer(fromAccount, toAccount int64, amount float64) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=bankgo sslmode=disable"
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
	return s.createAccountTable()
}

// TODO: Refactor to have seed and drop tables
func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists accounts (
		id serial primary key,
		firstName varchar(50),
		lastName varchar(50),
		number serial,
		balance serial,
		createdAt timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		INSERT INTO accounts 
		(firstName, lastName, number, balance, createdAt)
		VALUES 
		($1, $2, $3, $4, $5)
	`

	res, err := s.db.Query(query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	fmt.Printf("%+v\n", res)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int64) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int64) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) Transfer(fromAccount, toAccount int64, amount float64) error {
	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
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
