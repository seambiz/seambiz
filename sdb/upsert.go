package sdb
 
import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
)
  
// EscapeString returns mysql escaped string
func EscapeString(sql string) string {
	replacer := strings.NewReplacer("\\", "\\\\",
		"'", "\\'",
		"\\0", "\\\\0",
		"\n", "\\n",
		"\r", "\\r",
		`"`, `\"`,
		"\xef\xbf\xbd", "",
		"\x1a", "\\Z")

	return replacer.Replace(sql)
} 

// UpsertStatement helper for creating upsert statement
type UpsertStatement struct {
	sql                  *SQLStatement
	onduplicatekeyupdate string
	columns              []string
	appended             bool
	recordSet            bool
}

func (u *UpsertStatement) appendOnDuplicateKey() {
	if !u.appended && u.recordSet {
		if u.sql.buffer.Bytes()[u.sql.buffer.Len()-2] == ',' {
			u.sql.buffer.Truncate(u.sql.buffer.Len() - 2)
		}
		u.sql.Append(" ON DUPLICATE KEY UPDATE")
		u.sql.Append(u.onduplicatekeyupdate)
		u.appended = true
	}
}

// String return sql statement
func (u *UpsertStatement) String() string {
	u.appendOnDuplicateKey()
	return u.sql.String()
}

// Query frees the buffer aufter return sql string
func (u *UpsertStatement) Query() string {
	u.appendOnDuplicateKey()
	return u.sql.Query()
}

// InsertInto table name
func (u *UpsertStatement) InsertInto(table string) {
	u.sql = NewSQLStatement()
	u.onduplicatekeyupdate = ""

	u.sql.Append("INSERT INTO")
	u.sql.Append(table)
}

// Columns to be inserted
func (u *UpsertStatement) Columns(cols ...string) {
	u.columns = cols
	u.sql.Append("(")

	for i, col := range cols {
		u.sql.Append("`" + col + "`")
		if i < len(cols)-1 {
			u.sql.Append(",")
		}
	}

	u.sql.Append(") VALUES ")
}

// ColumnsByStruct convinience function
func (u *UpsertStatement) ColumnsByStruct(v interface{}) {
	var cols []string

	for _, f := range structs.New(v).Fields() {
		tag := f.Tag("db")
		switch tag {
		case "-":
			continue
		case "":
			cols = append(cols, f.Name())
		default:
			cols = append(cols, f.Tag("db"))
		}
	}
	u.Columns(cols...)
}

// OnDuplicateKeyUpdate what to do
func (u *UpsertStatement) OnDuplicateKeyUpdate(sqls []string) {
	u.onduplicatekeyupdate = strings.Join(sqls, ",")
}

// Record to be added to the statement
func (u *UpsertStatement) Record(values interface{}) {
	s := structs.New(values)
	s.TagName = "db"
	m := s.Map()
	u.recordSet = true

	u.sql.Append("(")
	for i, col := range u.columns {
		u.sql.Append("'" + EscapeString(fmt.Sprint(m[col])) + "'")
		if i < len(u.columns)-1 {
			u.sql.Append(",")
		}
	}
	u.sql.Append("),")
}
