package models

type JorfContainerResult struct {
	Id       string                 `json:"cid"`
	Title    string                 `json:"title"`
	Sections []JorfContainerSection `json:"sections"`
	Articles []JorfArticle          `json:"articles"`
}

type JorfContainerSection struct {
	Title    string                 `json:"title"`
	Articles []JorfArticle          `json:"articles"`
	Sections []JorfContainerSection `json:"sections"`
}

type JorfArticle struct {
	Content string `json:"content"`
	// This is a string. It represents a number ¯\_(ツ)_/¯
	Order string `json:"num"`
}
