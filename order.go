package jdb

type Order struct {
	field WhereField
	desc  bool
}

func (o Order) OrderField() string {
	return o.field.toWhereField()
}

func (o Order) OrderDesc() bool {
	return o.desc
}
