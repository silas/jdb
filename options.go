package jdb

type Option interface {
	jdbOption()
}

type option struct{}

func (o option) jdbOption() {}

type optionTable struct {
	option
	table string
}

func Table(table string) Option {
	return optionTable{table: table}
}

type optionReadOnly struct {
	option
	readOnly bool
}

func ReadOnly(readOnly bool) Option {
	return optionReadOnly{readOnly: readOnly}
}
