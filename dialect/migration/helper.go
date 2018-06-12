package migration

import (
	"bytes"
	"context"
	"database/sql"
	"text/template"
)

type Helper interface {
	Version() int
	GetRevision(ctx context.Context, tx *sql.Tx) (int, error)
	SetRevision(ctx context.Context, tx *sql.Tx, id int) error
	GetVersion(ctx context.Context, tx *sql.Tx) (int, error)
	SetVersion(ctx context.Context, tx *sql.Tx, version int) error
	TableExists(ctx context.Context, tx *sql.Tx) (bool, error)
	RenderSQL(sql string, locals ...map[string]interface{}) string
}

func Render(name string, text string, locals ...map[string]interface{}) string {
	b := &bytes.Buffer{}
	data := map[string]interface{}{}
	for _, local := range locals {
		for key, value := range local {
			data[key] = value
		}
	}
	err := template.Must(template.New(name).Parse(text)).Execute(b, data)
	if err != nil {
		panic(err)
	}
	return b.String()
}
