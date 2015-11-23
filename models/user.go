package models

import (
	"time"
)

type User struct {
    Id               int64
    LastLogin        time.Time
    Email            string `sql:"size:255; not null; unique"`
    FirstName        string `sql:"size:30"`
    LastName         string `sql:"size:30"`
    Password         string `sql:"size:255"`
    IsEnabled        bool   `sql:"default:true"`
    ConfirmPassword  string `sql:"-"`
}
