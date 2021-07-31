package sqlrs

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// Resultset contains a slice of the columns from a query, and the rows as a slice of label/value pairs.
type Resultset struct {
	Cols []string
	Rows []map[string]string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New() *Resultset {
	rs := new(Resultset)
	rs.Rows = make([]map[string]string, 0)
	rs.Cols = make([]string, 0)
	return rs
}

// QueryRows Executes query on db. Populates rs.Cols with column names, and Rows with a slice of rows of map[string]string.
func (rs *Resultset) QueryRows(db *sql.DB, query string) error {
	dbrows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer dbrows.Close()

	cols, err := dbrows.Columns()
	if err != nil {
		return err
	}
	for _, col := range cols {
		rs.Cols = append(rs.Cols, col)
	}
	vals := make([]interface{}, len(rs.Cols))
	for i := 0; i < len(rs.Cols); i++ {
		vals[i] = new(interface{})
	}

	for dbrows.Next() {
		err = dbrows.Scan(vals...)
		if err != nil {
			return err
		}
		row := make(map[string]string)
		for i := 0; i < len(vals); i++ {
			value := vals[i].(*interface{})
			switch v := (*value).(type) {
			case nil:
				row[rs.Cols[i]] = "NULL"
			case bool:
				if v {
					row[rs.Cols[i]] = "true"
				} else {
					row[rs.Cols[i]] = "false"
				}
			case []byte:
				row[rs.Cols[i]] = string(v)
			case time.Time:
				row[rs.Cols[i]] = v.Format("2006-01-02 15:04:05.999")
			case int64:
				row[rs.Cols[i]] = fmt.Sprintf("%d", v)
			case float64:
				row[rs.Cols[i]] = fmt.Sprintf("%f", v)
			default:
				row[rs.Cols[i]] = v.(string)
			}
		}
		rs.Rows = append(rs.Rows, row)
	}

	return dbrows.Err()
}

func (rs *Resultset) QueryRow(db *sql.DB, query string) (err error) {
	err = rs.QueryRows(db, query)
	if len(rs.Rows) == 0 {
		err = errors.New("sql: no Rows available")
	}
	if len(rs.Rows) > 1 {
		err = errors.New("sql: too many Rows")
	}
	return err
}

func (rs *Resultset) QueryRowString(db *sql.DB, query string) (result string, err error) {
	err = rs.QueryRow(db, query)
	if err != nil {
		return
	}
	result = rs.Rows[0][rs.Cols[0]]
	return
}

// Vprint vertically prints the resultset.
func (rs *Resultset) Vprint() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	for _, row := range rs.Rows {
		for _, col := range rs.Cols {
			out := strings.Replace(row[col], "\t", " ", -1)
			out = strings.Replace(out, "%", "%%", -1)
			fmt.Fprintf(w, "%s:\t%s\n", col, out)
		}
		fmt.Fprintln(w, "")
	}
	w.Flush()
}

// Hprint horizontally prints the resultset.
func (rs *Resultset) Hprint() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 1, ' ', 0)

	header := ""
	for _, col := range rs.Cols {
		header = header + fmt.Sprintf("%s\t", col)
	}
	fmt.Fprintln(w, header)

	for _, row := range rs.Rows {
		out := ""
		for _, col := range rs.Cols {
			val := strings.Replace(row[col], "\t", " ", -1)
			out = out + fmt.Sprintf("%s\t", val)
		}
		fmt.Fprintln(w, out)
	}
	w.Flush()
}
