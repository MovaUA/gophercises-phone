package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/movaua/gophercises-phone/pkg/model"
)

// DB provides database access to phones table
type DB struct {
	db *sql.DB
}

// Open opens a database specified by its data source name,
// usually consisting of at least a database name and connection information.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

// Close closes the database
func (db *DB) Close() error {
	return db.db.Close()
}

// Migrate applies migrations to the database
func (db *DB) Migrate() error {
	return createPhonesTable(db.db)
}

func createPhonesTable(db *sql.DB) error {
	statement := `
		create table if not exists public.phones
		(
			id serial,
			number varchar(255)
		)
	`

	_, err := db.Exec(statement)

	return err
}

// Insert inserts a phone and sets its ID if the operation is successful
func (db *DB) Insert(p *model.Phone) error {
	return db.db.
		QueryRow("insert into public.phones (number) values ($1) returning id", p.Number).
		Scan(&p.ID)
}

// Get get Phone by ID
func (db *DB) Get(id int64) (*model.Phone, error) {
	p := model.Phone{ID: id}
	if err := db.db.QueryRow("select number from public.phones where id=$1", id).Scan(&p.Number); err != nil {
		return nil, err
	}
	return &p, nil
}

// All gets all phones
func (db *DB) All() ([]model.Phone, error) {
	rows, err := db.db.Query("select id, number from public.phones")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var phones []model.Phone
	for rows.Next() {
		var p model.Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		phones = append(phones, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, nil
}

// Delete deletes a Phone by ID
func (db *DB) Delete(id int64) error {
	_, err := db.db.Exec("delete from public.phones where id=$1", id)
	return err
}

// Update updates a Phone
func (db *DB) Update(p model.Phone) error {
	_, err := db.db.Exec("update public.phones set number=$2 where id=$1", p.ID, p.Number)
	return err
}

// Find finds a Phone by its Number
func (db *DB) Find(number string) (*model.Phone, error) {
	p := model.Phone{Number: number}
	if err := db.db.QueryRow("select id from public.phones where number=$1", p.Number).Scan(&p.ID); err != nil {
		return nil, err
	}
	return &p, nil
}
