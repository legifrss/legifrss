package models

type TwitterJORF struct {
	JORFID       string                        `json:"jorf_id"`
	StatusID     int64                         `json:"status_id"`
	JORFContents map[string]TwitterJORFContent `json:"contents"`
	Date         int64                         `json:"date"`
}

type TwitterJORFContent struct {
	JORFContentID string `json:"jorf_content_id"`
	StatusID      int64  `json:"status_id"`
}
