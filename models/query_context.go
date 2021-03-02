package models

// QueryContext is used to search through the database
type QueryContext struct {
	Keyword string
	Author  string
	Nature  string
}
