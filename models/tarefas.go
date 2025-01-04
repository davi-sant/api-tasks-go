package models

type Tasks struct {
	ID        int    `json:"id"`
	Title     string `json:"titulo"`
	Descricao string `json:"descricao"`
	Status    bool   `json:"status"`
}
