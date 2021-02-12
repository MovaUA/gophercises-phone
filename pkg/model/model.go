// Package model defines models for the database tables
package model

// Phone represents a record in phones table
type Phone struct {
	// Unique ID
	ID int64
	// Phone number string
	Number string
}
