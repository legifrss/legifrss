package models

type TwitterJORF struct {
	JORFID       string                        `json:"jorf_id"`
	JORFTitle    string                        `json:"title"`
	StatusID     int64                         `json:"status_id"`
	JORFContents map[string]TwitterJORFContent `json:"contents"`
	Date         int64                         `json:"date"`
	URI          string                        `json:"uri"`
}

type TwitterJORFContent struct {
	JORFContentID string `json:"jorf_content_id"`
	Content       string `json:"content"`
	StatusID      int64  `json:"status_id"`
}
