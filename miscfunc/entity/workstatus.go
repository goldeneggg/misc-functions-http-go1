// Package entity provides entity
package entity

type Workstatus struct {
	ID      string            `json:"id"`
	Content WorkstatusContent `json:"content"`
}

type WorkstatusContent struct {
	YM                int64  `json:"ym"`
	BufferDaysPerWeek int64  `json:"buffer_days_per_week"`
	Desc              string `json:"desc"`
}

type DescWorkstatus struct {
	TableName string   `json:"table_name"`
	Attrs     []string `json:"attrs"`
	Status    string   `json:"status"`
}
