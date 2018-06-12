package sqlmock

import (
	"fmt"
	"strings"

	"github.com/silas/jdb/dialect"
)

type mockPath struct {
	parts []string
}

func (p *mockPath) Key(v string) dialect.Path {
	p.parts = append(p.parts, fmt.Sprintf(`.%s`, v))
	return p
}

func (p *mockPath) Index(v int) dialect.Path {
	p.parts = append(p.parts, fmt.Sprintf(`[%d]`, v))
	return p
}

func (p *mockPath) JSONExtract(column string) string {
	path := strings.Join(p.parts, "")
	return fmt.Sprintf("%s->'$%s'", column, path)
}
