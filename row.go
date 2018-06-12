package jdb

import (
	"fmt"
	"reflect"
	"time"

	"github.com/silas/jdb/internal/json"
	"github.com/silas/jdb/internal/ptr"
)

const (
	maxKind            = 64
	maxID              = 64
	maxUniqueStringKey = 255
	maxStringKey       = 255
)

var (
	errMustBeStruct = fmt.Errorf("input must be a struct")
)

type row struct {
	Kind            string
	ID              string
	ParentKind      *string
	ParentID        *string
	Data            *string
	UniqueStringKey *string
	StringKey       *string
	NumericKey      *float64
	TimeKey         *time.Time
	CreateTime      *time.Time
	UpdateTime      *time.Time
}

var timeValue = reflect.ValueOf(time.Time{})

func rowScanMeta(src interface{}, ro bool) (*row, error) {
	v := reflect.ValueOf(src)
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return nil, errMustBeStruct
	}

	r := &row{}

	var uniqueStringKeyDefined, stringKeyDefined, numericKeyDefined, timeKeyDefined bool
	if v, ok := src.(DatabaseUniqueStringKey); ok {
		if key, ok := v.DatabaseUniqueStringKey(); ok {
			r.UniqueStringKey = key
			uniqueStringKeyDefined = true
		}
	}
	if v, ok := src.(DatabaseStringKey); ok {
		if key, ok := v.DatabaseStringKey(); ok {
			r.StringKey = key
			stringKeyDefined = true
		}
	}
	if v, ok := src.(DatabaseNumericKey); ok {
		if key, ok := v.DatabaseNumericKey(); ok {
			r.NumericKey = key
			numericKeyDefined = true
		}
	}
	if v, ok := src.(DatabaseTimeKey); ok {
		if key, ok := v.DatabaseTimeKey(); ok {
			r.TimeKey = key
			timeKeyDefined = true
		}
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		name, tagOpts := json.ParseTag(tag)

		value := v.Field(i)
		if !value.IsValid() {
			continue
		}
		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}
		value = reflect.Indirect(value)

		switch name {
		case kindTag:
			if v, ok := value.Interface().(string); ok {
				if v != "" {
					r.Kind = v
				}
			} else {
				return nil, fmt.Errorf("%s is invalid: %v", name, value)
			}
		case idTag:
			if v, ok := value.Interface().(string); ok {
				if v != "" {
					r.ID = v
				}
			} else {
				return nil, fmt.Errorf("%s is invalid: %v", name, value)
			}
		case parentKindTag:
			if value.Kind() == reflect.String {
				s := value.Interface().(string)
				if s != "" {
					r.ParentKind = ptr.String(s)
				}
			} else {
				return nil, fmt.Errorf("%s is invalid: %v", name, value)
			}
		case parentIDTag:
			if value.Kind() == reflect.String {
				s := value.Interface().(string)
				if s != "" {
					r.ParentID = ptr.String(s)
				}
			} else {
				return nil, fmt.Errorf("%s is invalid: %v", name, value)
			}
		case createTimeTag, updateTimeTag:
			if !ro {
				continue
			}
			if value.Kind() == timeValue.Kind() && value.Type() == timeValue.Type() {
				s := value.Interface().(time.Time)
				if !isZero(value.Type(), value) {
					if name == createTimeTag {
						r.CreateTime = &s
					} else {
						r.UpdateTime = &s
					}
				}
			} else {
				return nil, fmt.Errorf("%s is invalid: %v", name, value)
			}
		default:
			omitempty := tagOpts.Contains("omitempty")

			if tagOpts.Contains(uniqueStringKeyTag) {
				if v, ok := value.Interface().(string); ok {
					if !omitempty || !isZero(value.Type(), value) {
						if uniqueStringKeyDefined {
							return nil, fmt.Errorf("has duplicate unique string keys")
						}
						uniqueStringKeyDefined = true

						r.UniqueStringKey = ptr.String(v)
					}
				} else {
					return nil, fmt.Errorf("unique string key is invalid: %v", value)
				}
			}

			if tagOpts.Contains(stringKeyTag) {
				if v, ok := value.Interface().(string); ok {
					if !omitempty || !isZero(value.Type(), value) {
						if stringKeyDefined {
							return nil, fmt.Errorf("has duplicate string keys")
						}
						stringKeyDefined = true

						r.StringKey = ptr.String(v)
					}
				} else {
					return nil, fmt.Errorf("string key is invalid: %v", value)
				}
			}

			if tagOpts.Contains(numericKeyTag) {
				if v, ok := value.Interface().(float64); ok {
					if !omitempty || !isZero(value.Type(), value) {
						if numericKeyDefined {
							return nil, fmt.Errorf("has duplicate numeric keys")
						}
						numericKeyDefined = true

						r.NumericKey = ptr.Float64(v)
					}
				} else {
					return nil, fmt.Errorf("numeric key is invalid: %v", value)
				}
			}

			if tagOpts.Contains(timeKeyTag) {
				if v, ok := value.Interface().(time.Time); ok {
					if !omitempty || !v.IsZero() {
						if timeKeyDefined {
							return nil, fmt.Errorf("has duplicate time keys")
						}
						timeKeyDefined = true

						r.TimeKey = &v
					}
				} else {
					return nil, fmt.Errorf("time key is invalid: %v", value)
				}
			}
		}
	}

	return r, nil
}

func rowScanInput(kind string, src interface{}) (*row, error) {
	r, err := rowScanMeta(src, false)
	if err != nil {
		return nil, err
	}

	if r.Kind != "" && r.Kind != kind {
		return nil, fmt.Errorf("kind mismatch: %s != %s", kind, r.Kind)
	}

	r.Kind = kind

	if r.Kind == "" {
		return nil, fmt.Errorf("kind not defined")
	} else if len(r.ID) > maxKind {
		return nil, fmt.Errorf("kind max length 64 characters: %s (%d)", r.Kind, len(r.Kind))
	}
	if r.ID == "" {
		return nil, fmt.Errorf("id not defined")
	} else if len(r.ID) > maxID {
		return nil, fmt.Errorf("id max length 64 characters: %s (%d)", r.ID, len(r.ID))
	}
	if r.ParentID != nil && r.ParentKind == nil {
		return nil, fmt.Errorf("parent kind not defined")
	}
	if r.ParentKind != nil && len(*r.ParentKind) > maxKind {
		return nil, fmt.Errorf("parent kind max length 64 characters: %s (%d)", *r.ParentKind, len(*r.ParentKind))
	}
	if r.ParentID != nil && len(*r.ParentID) > maxID {
		return nil, fmt.Errorf("parent id max length 64 characters: %s (%d)", *r.ParentID, len(*r.ParentID))
	}
	if r.UniqueStringKey != nil && len(*r.UniqueStringKey) > maxUniqueStringKey {
		return nil, fmt.Errorf("unique string key max length 64 characters: %s (%d)", *r.UniqueStringKey, len(*r.UniqueStringKey))
	}
	if r.StringKey != nil && len(*r.StringKey) > maxStringKey {
		return nil, fmt.Errorf("string key max length 64 characters: %s (%d)", *r.StringKey, len(*r.StringKey))
	}

	if v, err := json.Marshal(src); err != nil {
		return nil, err
	} else if len(v) > 2 {
		r.Data = ptr.String(string(v))
	}

	return r, nil
}

func idScanInput(src interface{}) (string, string, error) {
	r, err := rowScanMeta(src, false)
	if err != nil {
		return "", "", err
	}

	if r.ID == "" {
		err = ErrIDNotFound
	}

	return r.Kind, r.ID, err
}
