package models

type TwitterJORF struct {
	StatusID     int64            `json:"status_id"`
	JORFContents map[string]int64 `json:"contents"`
}
