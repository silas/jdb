# jdb

jdb is a SQL query builder for the JSON features of MySQL, PostgreSQL and SQLite.

This package is currently in development and the API is not stable.

## Usage

``` go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/silas/jdb"
	_ "github.com/silas/jdb/dialect/sqlite3"
)

type Name struct {
	GivenName  string   `jdb:",omitempty"`
	FamilyName string   `jdb:",omitempty"`
	Aliases    []string `jdb:",omitempty"`
}

type User struct {
    Kind       string    `jdb:"-kind"`
    ID         string    `jdb:"-id"`
    Email      string    `jdb:",uniquestringkey"`
    Name       Name      `jdb:",omitempty"`
    Age        int       `jdb:",omitempty"`
    CreateTime time.Time `jdb:"-createtime"`
    UpdateTime time.Time `jdb:"-updatetime"`
}

func (u User) DatabaseNumericKey() (*float64, bool) {
	if u.Age > 0 {
		age := float64(u.Age)
		return &age, true
	} else {
		return nil, true
	}
}

func main() {
	ctx := context.Background()

	db, err := jdb.Open("sqlite3", "test.db?cache=shared")
	if err != nil {
		panic(err)
	}

	err = db.Migrate(ctx)
	if err != nil {
		panic(err)
	}

	err = db.Tx(ctx, func(tx *jdb.Tx) error {
		user := User{
			ID:    "1",
			Email: "jane@example.com",
			Name: Name{
				GivenName:  "Jane",
				FamilyName: "Doe",
				Aliases:    []string{"Janie", "Roe"},
			},
			Age: 34,
		}

		err = db.Query("user").Insert(user).Exec(ctx, tx)
		if err != nil {
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		panic(err)
	}

	err = db.Tx(ctx, func(tx *jdb.Tx) error {
		var user User
		err = db.Query("user").Get("1").Select().First(ctx, tx, &user)
		if err != nil {
			return err
		}

		fmt.Println(user)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
```

Run with required sqlite3 flags:

``` sh
$ go run \
  -tags "libsqlite3 $(uname -s | tr '[:upper:]' '[:lower:]') json1" \
  main.go
```

This will result in a row that looks something like the following:

``` sql
sqlite> SELECT * FROM jdb WHERE kind = 'user' AND id = '1';
             kind = user
               id = 1
      parent_kind = ¤
        parent_id = ¤
unique_string_key = jane@example.com
       string_key = ¤
      numeric_key = 34.0
         time_key = ¤
             data = {"Email":"jane@example.com","Name":{"GivenName":"Jane","FamilyName":"Doe","Aliases":["Janie","Roe"]},"Age":34}
      create_time = 2018-06-12 12:41:37.035
      update_time = 2018-06-12 12:41:37.035
```

With the following schema (in SQLite):

```
sqlite> .schema
CREATE TABLE jdb (
  kind VARCHAR(64),
  id VARCHAR(64),
  parent_kind VARCHAR(64),
  parent_id VARCHAR(64),
  unique_string_key VARCHAR(255),
  string_key VARCHAR(255),
  numeric_key REAL,
  time_key DATETIME,
  data JSON,
  create_time DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
  update_time DATETIME NOT NULL DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
  PRIMARY KEY (kind, id),
  FOREIGN KEY (parent_kind, parent_id) REFERENCES jdb (kind, id)
);
CREATE INDEX jdb_r2 ON jdb (create_time);
CREATE INDEX jdb_r3 ON jdb (update_time);
CREATE UNIQUE INDEX jdb_r4 ON jdb (kind, unique_string_key);
CREATE INDEX jdb_r5 ON jdb (kind, string_key);
CREATE INDEX jdb_r6 ON jdb (kind, numeric_key);
CREATE INDEX jdb_r7 ON jdb (kind, time_key);
CREATE INDEX jdb_r8 ON jdb (kind, create_time);
CREATE INDEX jdb_r9 ON jdb (kind, update_time);
CREATE INDEX jdb_r10 ON jdb (kind, id, parent_kind);
CREATE INDEX jdb_r11 ON jdb (parent_kind, parent_id, kind, unique_string_key);
CREATE INDEX jdb_r12 ON jdb (parent_kind, parent_id, kind, string_key);
CREATE INDEX jdb_r13 ON jdb (parent_kind, parent_id, kind, numeric_key);
CREATE INDEX jdb_r14 ON jdb (parent_kind, parent_id, kind, time_key);
CREATE INDEX jdb_r15 ON jdb (parent_kind, parent_id, kind, create_time);
CREATE INDEX jdb_r16 ON jdb (parent_kind, parent_id, kind, update_time);
```

## License

This work is licensed under the MIT License (see the LICENSE and NOTICE files).
