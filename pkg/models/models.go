package models

import "time"

type LegifranceElement struct {
	Title            string   `json:"title"`
	ID               string   `json:"id"`
	Nature           string   `json:"nature"`
	Link             string   `json:"link"`
	Description      string   `json:"description"`
	Category         []string `json:"category"`
	Author           string   `json:"author"`
	SanitizedAuthor  string
	SanitizedNature  string
	Date             time.Time `json:"date"`
	Content          string    `json:"content"`
	TwitterPublished bool      `json:"twitter_published"`
}

type FeedDescription struct {
	TitleSuffix       string
	LinkSuffix        string
	DescriptionSuffix string
}

type Container struct {
	ID string `json:"id"`
}

type LastNJo struct {
	Containers []Container `json:"containers"`
}

type Summary struct {
	ID       string `json:"id"`
	Title    string `json:"titre"`
	Nature   string `json:"nature"`
	Minister string `json:"ministere"`
	Emitter  string `json:"emetteur"`
}

type HierarchyStep struct {
	Title     string          `json:"titre"`
	Level     int             `json:"niv"`
	Step      []HierarchyStep `json:"tms"`
	Summaries []Summary       `json:"liensTxt"`
}

type Structure struct {
	Contents []HierarchyStep `json:"tms"`
}

type JOContainer struct {
	ID        string    `json:"id"`
	Structure Structure `json:"structure"`
	Timestamp int64     `json:"datePubli"`
}

type Item struct {
	Container JOContainer `json:"joCont"`
}

type JOContainerResult struct {
	Items []Item `json:"items"`
}
