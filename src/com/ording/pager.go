package ording

type Pager struct {
	Total int                      `json:"total"`
	Rows  []map[string]interface{} `json:"rows"`
	Text  string
}
