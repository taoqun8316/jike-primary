package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
	Utime    time.Time
}
