package jdb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/silas/jdb/internal/json"
)

type Rows struct {
	*sql.Rows

	columns []SelectField
}

func newRows(rows *sql.Rows, columns []SelectField) *Rows {
	return &Rows{rows, columns}
}

func (rs *Rows) Close() error {
	return rs.Rows.Close()
}

func (rs *Rows) Err() error {
	return rs.Rows.Err()
}

func (rs *Rows) Next() bool {
	return rs.Rows.Next()
}

func (rs *Rows) scan(dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}
	s := dest.Elem()
	if len(rs.columns) > 1 && s.Kind() != reflect.Struct {
		return errors.New("dest must be a struct")
	}
	if s.Kind() == reflect.Struct {
		return rs.scanColumns(dest)
	} else {
		return rs.Rows.Scan(dest.Interface())
	}
}

func (rs *Rows) scanColumns(dest reflect.Value) error {
	if dest.IsNil() {
		return errors.New("dest must be non-nil")
	}

	s := dest.Elem()
	s.Set(reflect.Zero(s.Type()))

	var kind, id, parentKind, parentID, data *string
	var createTime, updateTime *time.Time

	var columns []interface{}
	for _, c := range rs.columns {
		switch c {
		case kindField:
			columns = append(columns, &kind)
		case idField:
			columns = append(columns, &id)
		case parentKindField:
			columns = append(columns, &parentKind)
		case parentIdField:
			columns = append(columns, &parentID)
		case dataField:
			columns = append(columns, &data)
		case createTimeField:
			columns = append(columns, &createTime)
		case updateTimeField:
			columns = append(columns, &updateTime)
		}
	}

	err := rs.Rows.Scan(columns...)
	if err != nil {
		return err
	}

	if data != nil && *data != "" {
		err = json.Unmarshal([]byte(*data), dest.Interface())
		if err != nil {
			return err
		}
	}

	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		value := s.Field(i)
		if !value.IsValid() {
			continue
		}
		if !value.CanSet() {
			continue
		}

		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		name, _ := json.ParseTag(tag)

		switch name {
		case idTag:
			if id == nil {
				continue
			}
			var tv reflect.Value
			if value.Kind() == reflect.Ptr {
				tv = reflect.ValueOf(id)
			} else {
				tv = reflect.ValueOf(*id)
			}
			if value.Kind() != tv.Kind() {
				return fmt.Errorf("%s must be a string", idTag)
			}
			value.Set(tv)
		case kindTag:
			if kind == nil {
				continue
			}
			var tv reflect.Value
			if value.Kind() == reflect.Ptr {
				tv = reflect.ValueOf(kind)
			} else {
				tv = reflect.ValueOf(*kind)
			}
			if value.Kind() != tv.Kind() {
				return fmt.Errorf("%s must be a string", kindTag)
			}
			value.Set(tv)
		case parentKindTag:
			if parentKind == nil {
				continue
			}
			var tv reflect.Value
			if value.Kind() == reflect.Ptr {
				tv = reflect.ValueOf(parentKind)
			} else if parentKind != nil {
				tv = reflect.ValueOf(*parentKind)
			} else {
				tv = stringValue
			}
			if value.Kind() != tv.Kind() {
				return fmt.Errorf("%s must be a string", parentKindTag)
			}
			if parentKind != nil {
				value.Set(tv)
			}
		case parentIDTag:
			if parentID == nil {
				continue
			}
			var tv reflect.Value
			if value.Kind() == reflect.Ptr {
				tv = reflect.ValueOf(parentID)
			} else if parentID != nil {
				tv = reflect.ValueOf(*parentID)
			} else {
				tv = stringValue
			}
			if value.Kind() != tv.Kind() {
				return fmt.Errorf("%s must be a string", parentIDTag)
			}
			if parentID != nil {
				value.Set(tv)
			}
		case createTimeTag:
			if createTime == nil {
				continue
			}
			var tv reflect.Value
			if value.Kind() == reflect.Ptr {
				tv = reflect.ValueOf(createTime)
			} else {
				tv = reflect.ValueOf(*createTime)
			}
			if value.Kind() != tv.Kind() {
				return fmt.Errorf("%s must be a time.Time value", createTimeTag)
			}
			value.Set(tv)
		case updateTimeTag:
			if updateTime == nil {
				continue
			}
			var tv reflect.Value
			if value.Kind() == reflect.Ptr {
				tv = reflect.ValueOf(updateTime)
			} else {
				tv = reflect.ValueOf(*updateTime)
			}
			if value.Kind() != tv.Kind() {
				return fmt.Errorf("%s must be a time.Time value", updateTimeTag)
			}
			value.Set(tv)
		}
	}

	return nil
}

func (rs *Rows) Scan(dest interface{}) error {
	return rs.scan(reflect.ValueOf(dest))
}

func (rs *Rows) scanAll(dest reflect.Value) error {
	if dest.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}
	v := dest.Elem()
	if v.Kind() != reflect.Slice {
		return errors.New("dest must be a slice")
	}
	if v.Len() != 0 {
		v.Set(reflect.MakeSlice(dest.Type().Elem(), 0, dest.Elem().Cap()))
	}
	t := v.Type().Elem()
	ptr := t.Kind() == reflect.Ptr
	if ptr {
		t = t.Elem()
	}

	for rs.Next() {
		e := reflect.New(t)
		if err := rs.scan(e); err != nil {
			return err
		}
		if ptr {
			v.Set(reflect.Append(v, e))
		} else {
			v.Set(reflect.Append(v, e.Elem()))
		}
	}

	return rs.Err()
}

func (rs *Rows) ScanAll(dest interface{}) error {
	return rs.scanAll(reflect.ValueOf(dest))
}
