package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/movaua/gophercises-phone/pkg/db"
	"github.com/movaua/gophercises-phone/pkg/model"
	"github.com/movaua/gophercises-phone/pkg/phone"
)

func main() {
	db, err := db.Open("postgres", "dbname=testdb port=5432 user=user password=test sslmode=disable")
	must(err)
	defer db.Close()

	must(db.Migrate())

	p := &model.Phone{Number: "1234567890"}
	must(db.Insert(p))
	fmt.Printf("%+v\n", p)

	p, err = db.Get(p.ID)
	must(err)
	fmt.Printf("%+v\n", p)

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
		must(db.Insert(&model.Phone{Number: n}))
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("All phones:")
	phones, err := db.All()
	must(err)
	for _, p := range phones {
		fmt.Printf("%+v\n", p)
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("Normalizing and making unique:")
	for _, p := range phones {
		fmt.Printf("Processing: %+v...\n", p)
		n := phone.Norm(p.Number)
		found, err := db.Find(n)
		if errors.Is(err, sql.ErrNoRows) {
			if n == "" {
				must(db.Delete(p.ID))
				fmt.Println("Invalid phone is deleted")
				continue
			}
			must(db.Update(model.Phone{ID: p.ID, Number: n}))
			fmt.Println("Phone is normalized")
			continue
		}
		must(err)
		switch {
		case found.ID != p.ID:
			must(db.Delete(p.ID))
			fmt.Println("Duplicate is deleted")
		case n != p.Number:
			must(db.Update(model.Phone{ID: found.ID, Number: n}))
			fmt.Println("Phone is normalized")
		default:
			fmt.Println("Nothing to do")
		}
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("All phones:")
	phones, err = db.All()
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
