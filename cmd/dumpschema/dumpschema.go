package main

import (
	"database/sql"
	"flag"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	args := struct {
		dbName     string
		schemaName string
		columnName string
	}{}
	flag.StringVar(&args.dbName, "db", "postgres", "The name of the database to dump table names and column names of.")
	flag.StringVar(&args.schemaName, "schema", "information_schema", "The name of the schema to dump table names and column names of.")
	flag.StringVar(&args.columnName, "column", "^*$", "A regular expression to match based on column name in order print tables names.")
	flag.Parse()

	var err error
	if db, err = sql.Open("postgres", fmt.Sprintf("dbname=%s sslmode=disable", args.dbName)); err != nil {
		panic(err)
	}
	defer db.Close()
	dumpSchema(db, args.schemaName, regexp.MustCompile(args.columnName))
}

func dumpSchema(db *sql.DB, schemaName string, colRe *regexp.Regexp) {
	t, err := tableNames(db, schemaName)
	if err != nil {
		panic(err)
	}
	for _, n := range t {
		c, err := columnNames(db, n)
		if err != nil {
			panic(err)
		}
		cc := make([]string, 0, len(c))
		for i, m := range c {
			if colRe.MatchString(m) {
				cc = append(cc, c[i])
			}
		}
		if len(cc) > 0 {

			fmt.Printf("\n\n%v\n", n)
			for _, m := range cc {
				fmt.Printf("\t\t%v\n", m)
			}
		}
	}
}

func tableNames(db *sql.DB, schemaName string) ([]string, error) {
	q := `
select 
 	table_name
from
	information_schema.tables
where
	table_type = 'BASE TABLE'
		and
	table_schema = $1
order by
	table_name
`
	return queryForVarCharColumn(db, q, schemaName)
}

func columnNames(db *sql.DB, tableName string) ([]string, error) {
	q := `
select 
 	column_name
from
	information_schema.columns
where
	table_name = $1
order by
	column_name
`
	return queryForVarCharColumn(db, q, tableName)
}

func queryForVarCharColumn(db *sql.DB, query string, arg string) ([]string, error) {
	t := make([]string, 0)
	rows, err := db.Query(query, arg)
	if err != nil {
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return t, err
		}
		t = append(t, s)
	}
	return t, rows.Err()
}
