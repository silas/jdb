package dialect

type OrderField interface {
	OrderField() string
	OrderDesc() bool
}

type Path interface {
	Key(v string) Path
	Index(v int) Path
	JSONExtract(column string) string
}
