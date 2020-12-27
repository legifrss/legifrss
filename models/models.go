package models

import "time"

type LegifranceElement struct {
	Title           string
	Id              string
	Nature          string
	Link            string
	Description     string
	Category        []string
	Author          string
	SanitizedAuthor string
	SanitizedNature string
	Date            time.Time
	Content         string
}

type FeedDescription struct {
	TitleSuffix       string
	LinkSuffix        string
	DescriptionSuffix string
}

//-----------------------
type Container struct {
	Id string `json:"id"`
}

type LastNJo struct {
	Containers []Container `json:"containers"`
}

//------------------------
type Summary struct {
	Id       string `json:"id"`
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
	Id        string    `json:"id"`
	Structure Structure `json:"structure"`
	Timestamp int64     `json:"datePubli"`
}

type Item struct {
	Container JOContainer `json:"joCont"`
}

type JOContainerResult struct {
	Items []Item `json:"items"`
}
