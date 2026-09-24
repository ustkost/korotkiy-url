package model

import "time"

type Click struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	Timestamp time.Time `json:"timestamp"`
	Referrer  string    `json:"referrer"`
}
