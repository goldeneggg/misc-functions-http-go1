// Package entity provides entity
package entity

type Workstatus struct {
	ID     string `json:"id"`
	YM     int64  `json:"ym"`
	Status string `json:"status"`
	Desc   string `json:"desc"`
}
