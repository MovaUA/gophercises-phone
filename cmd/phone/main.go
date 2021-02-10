package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	ph "github.com/movaua/gophercises-phone/pkg/phone"
)

func main() {
	db, err := sql.Open("postgres", "dbname=testdb port=5432 user=user password=test sslmode=disable")
	must(err)
	defer db.Close()

	must(db.Ping())

	must(createPhonesTable(db))

	id, err := insertPhone(db, "1234567890")
	must(err)
	fmt.Printf("id=%d\n", id)

	number, err := getPhone(db, id)
	must(err)
	fmt.Printf("phone=%q\n", number)

	numbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"(123)456-7892",
		"there is no phone here",
	}
	for _, n := range numbers {
		_, err = insertPhone(db, n)
		must(err)
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("All phones:")
	phones, err := allPhones(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("%+v\n", p)
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("Normalizing and making unique:")
	for _, p := range phones {
		fmt.Printf("Processing: %+v...\n", p)
		n := ph.Norm(p.number)
		id, err := findPhone(db, n)
		if errors.Is(err, sql.ErrNoRows) {
			if n == "" {
				must(deletePhone(db, p.id))
				fmt.Println("Invalid phone is deleted")
				continue
			}
			must(updatePhone(db, p.id, n))
			fmt.Println("Phone is normalized")
			continue
		}
		must(err)
		switch {
		case id != p.id:
			must(deletePhone(db, p.id))
			fmt.Println("Duplicate is deleted")
		case n != p.number:
			must(updatePhone(db, id, n))
			fmt.Println("Phone is normalized")
		default:
			fmt.Println("Nothing to do")
		}
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("All phones:")
	phones, err = allPhones(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("%+v\n", p)
	}
}

func must(err error) {
	if err == nil {
		return
	}
	log.Panicln(err)
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

func insertPhone(db *sql.DB, number string) (int64, error) {
	var id int64
	if err := db.QueryRow("insert into public.phones (number) values ($1) returning id", number).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func getPhone(db *sql.DB, id int64) (string, error) {
	var number string
	err := db.QueryRow("select number from public.phones where id=$1", id).Scan(&number)
	return number, err
}

type phone struct {
	id     int64
	number string
}

func allPhones(db *sql.DB) ([]phone, error) {
	rows, err := db.Query("select id, number from public.phones order by id asc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var phones []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		phones = append(phones, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, nil
}

func deletePhone(db *sql.DB, id int64) error {
	_, err := db.Exec("delete from public.phones where id=$1", id)
	return err
}

func updatePhone(db *sql.DB, id int64, number string) error {
	_, err := db.Exec("update public.phones set number=$2 where id=$1", id, number)
	return err
}

func findPhone(db *sql.DB, number string) (int64, error) {
	var id int64
	row := db.QueryRow("select id from public.phones where number=$1", number)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
