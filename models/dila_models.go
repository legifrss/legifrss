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
	Order   string `json:"num"`
}

type SortByOrder []JorfArticle

func (a SortByOrder) Len() int           { return len(a) }
func (a SortByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }

type SortByArticleOrder []JorfContainerSection

func (a SortByArticleOrder) Len() int      { return len(a) }
func (a SortByArticleOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByArticleOrder) Less(i, j int) bool {

	return lowArticleOrderInSection(a[i]) < lowArticleOrderInSection(a[j])
}

func lowArticleOrderInSection(section JorfContainerSection) string {
	if len(section.Articles) == 0 {
		return "-1"
	} else {
		return section.Articles[0].Order
	}
}
