package sqlite3

import (
	"fmt"
	"strings"

	"github.com/silas/jdb/dialect"
)

type sqlite3Path struct {
	parts []string
}

func (p *sqlite3Path) Key(v string) dialect.Path {
	v = strings.Replace(v, `"`, `\"`, -1)
	p.parts = append(p.parts, fmt.Sprintf(`."%s"`, v))
	return p
}

func (p *sqlite3Path) Index(v int) dialect.Path {
	p.parts = append(p.parts, fmt.Sprintf(`[%d]`, v))
	return p
}

func (p *sqlite3Path) JSONExtract(column string) string {
	path := strings.Join(p.parts, "")
	path = strings.Replace(path, "'", `''`, -1)
	return fmt.Sprintf("json_extract(%s, '$%s')", column, path)
}
