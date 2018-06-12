package postgres

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/silas/jdb/dialect"
)

type postgresPath struct {
	parts []string
}

func (p *postgresPath) Key(v string) dialect.Path {
	v = strings.Replace(v, `"`, `\"`, -1)
	p.parts = append(p.parts, fmt.Sprintf(`"%s"`, v))
	return p
}

func (p *postgresPath) Index(v int) dialect.Path {
	p.parts = append(p.parts, strconv.Itoa(v))
	return p
}

func (p *postgresPath) JSONExtract(column string) string {
	path := strings.Join(p.parts, ",")
	path = strings.Replace(path, "'", `''`, -1)
	return fmt.Sprintf("%s#>>'{%s}'", column, path)
}
