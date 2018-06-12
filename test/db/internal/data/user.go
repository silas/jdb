package data

import (
	"strings"
	"time"
)

type Name struct {
	GivenName  string   `jdb:",omitempty"`
	FamilyName string   `jdb:",omitempty"`
	Aliases    []string `jdb:",omitempty"`
}

type User struct {
	Kind        string    `jdb:"-kind"`
	ID          string    `jdb:"-id"`
	Email       string    `jdb:",uniquestringkey"`
	Name        Name      `jdb:",omitempty"`
	Age         int       `jdb:",omitempty"`
	RefreshTime time.Time `jdb:",omitempty"`
	CreateTime  time.Time `jdb:"-createtime"`
	UpdateTime  time.Time `jdb:"-updatetime"`
}

type UserMeta struct {
	ID         string    `jdb:"-id"`
	Email      string    `jdb:",uniquestringkey"`
	CreateTime time.Time `jdb:"-createtime"`
}

func (u User) DatabaseStringKey() (*string, bool) {
	var domain *string
	if strings.Contains(u.Email, "@") {
		p := strings.Split(u.Email, "@")
		v := strings.ToLower(p[len(p)-1])
		domain = &v
	}
	return domain, true
}

func (u User) DatabaseNumericKey() (*float64, bool) {
	var length = float64(len(u.Email))
	return &length, true
}

func (u User) DatabaseTimeKey() (*time.Time, bool) {
	if !u.RefreshTime.IsZero() {
		return &u.RefreshTime, true
	}
	return nil, true
}

type NamePtr struct {
	GivenName  *string    `jdb:",omitempty"`
	FamilyName *string    `jdb:",omitempty"`
	Aliases    *[]*string `jdb:",omitempty"`
}

type UserPtr struct {
	Kind        *string    `jdb:"-kind"`
	ID          *string    `jdb:"-id"`
	Email       *string    `jdb:",uniquestringkey"`
	Name        *NamePtr   `jdb:",omitempty"`
	Age         *int       `jdb:",omitempty"`
	RefreshTime *time.Time `jdb:",omitempty"`
	CreateTime  *time.Time `jdb:"-createtime"`
	UpdateTime  *time.Time `jdb:"-updatetime"`
}

type UserMetaPtr struct {
	ID         *string    `jdb:"-id"`
	Email      *string    `jdb:",uniquestringkey"`
	CreateTime *time.Time `jdb:"-createtime"`
}

func (u UserPtr) DatabaseStringKey() (*string, bool) {
	var domain *string
	if u.Email != nil && strings.Contains(*u.Email, "@") {
		p := strings.Split(*u.Email, "@")
		v := strings.ToLower(p[len(p)-1])
		domain = &v
	}
	return domain, true
}

func (u UserPtr) DatabaseNumericKey() (*float64, bool) {
	var length float64 = 0
	if u.Email != nil {
		length = float64(len(*u.Email))
	}
	return &length, true
}

func (u UserPtr) DatabaseTimeKey() (*time.Time, bool) {
	if u.RefreshTime != nil && !u.RefreshTime.IsZero() {
		return u.RefreshTime, true
	}
	return nil, true
}
