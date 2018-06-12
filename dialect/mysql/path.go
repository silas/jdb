package mysql

import (
	"fmt"
	"strings"

	"github.com/silas/jdb/dialect"
)

type mysqlPath struct {
	parts []string
}

func (p *mysqlPath) Key(v string) dialect.Path {
	v = strings.Replace(v, `"`, `\\"`, -1)
	p.parts = append(p.parts, fmt.Sprintf(`."%s"`, v))
	return p
}

func (p *mysqlPath) Index(v int) dialect.Path {
	p.parts = append(p.parts, fmt.Sprintf(`[%d]`, v))
	return p
}

func (p *mysqlPath) JSONExtract(column string) string {
	path := strings.Join(p.parts, "")
	path = strings.Replace(path, "'", `''`, -1)
	return fmt.Sprintf("json_unquote(json_extract(%s, '$%s'))", column, path)
}
