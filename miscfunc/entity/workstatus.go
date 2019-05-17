// Package entity provides entity
package entity

type Workstatus struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type DescWorkstatus struct {
	TableName string `json:"table_name"`
	Attrs     string `json:"attrs"`
	Status    string `json:"status"`
}
