package entity

import "time"

type User struct {
	ID        int64
	Username  string
	Password  []byte
	CreatedAt time.Time
}
