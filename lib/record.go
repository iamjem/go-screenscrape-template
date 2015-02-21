package screenscrape

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Record struct {
	Id       int64     `db:"id"`
	Source   string    `db:"source"`
	Latest   string    `db:"latest"`
	Modified time.Time `db:"modified"`
}

func (p *Record) PreInsert(s gorp.SqlExecutor) error {
	p.Modified = time.Now().UTC()
	return nil
}

func (p *Record) PreUpdate(s gorp.SqlExecutor) error {
	p.Modified = time.Now()
	return nil
}
