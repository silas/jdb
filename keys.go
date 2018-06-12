package jdb

import "time"

type DatabaseUniqueStringKey interface {
	DatabaseUniqueStringKey() (*string, bool)
}

type DatabaseStringKey interface {
	DatabaseStringKey() (*string, bool)
}

type DatabaseNumericKey interface {
	DatabaseNumericKey() (*float64, bool)
}

type DatabaseTimeKey interface {
	DatabaseTimeKey() (*time.Time, bool)
}
