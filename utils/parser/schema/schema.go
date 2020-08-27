package schema

import (
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"strings"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

type (
	Schema struct {
		Name		string
		Path		string
		Structure	Structure
	}

	Structure struct {
		Table		Table		`yaml:"table"`
		Columns		[]Column	`yaml:"columns"`
	}

	Table struct {
		Engine		string	`yaml:"engine"`
		Charset		string	`yaml:"charset"`
		Collation	string	`yaml:"collation"`
	}

	Column struct {
		Name		string		`yaml:"name"`
		Type		string		`yaml:"type"`
		Primary		bool		`yaml:"primary"`
		Unique		bool		`yaml:"unique"`
		Nullable	bool		`yaml:"nullable"`
		Increment	bool		`yaml:"increment"`
	}
)


func (schema *Schema) File(path string) {

	if _, err := os.Stat(path);
	os.IsNotExist(err) { log.Fatal(err) }

	/*
	We have a working schema file and now we can
	set known properties of the schema using the
	migration file
	*/
	schema.Path, _ = filepath.Abs(path)
	schema.Name = strings.Replace(filepath.Base(schema.Path), ".yml", "", 1)

	/*
	We need to read the file and append the table
	properties to the Schema struct to continue
	the migration file parse.
	*/
	buffer, err := ioutil.ReadFile(schema.Path)
	if err != nil { log.Fatal(err) }

	structure := Structure {}
	yaml.Unmarshal(buffer, &structure)
	schema.Structure = structure

}


func (schema *Schema) Statement() string {
	return fmt.Sprintf("CREATE TABLE `%s` (%s) ENGINE=%s DEFAULT CHARSET=%s DEFAULT COLLATE=%s", schema.Name, schema.GenerateColumnBaseStatement(), schema.Structure.Table.Engine, schema.Structure.Table.Charset, schema.Structure.Table.Collation)
}



/*
Generate sql statment lines for each column we have
in our migration answer file to append to the final
statement.
*/
func (schema *Schema) GenerateColumnBaseStatement() string {

	var statement string
	for _, column := range schema.Structure.Columns {

		// BASE
		sql := fmt.Sprintf("`%s` %s", column.Name, column.Type)

		// NULLABLE
		if column.Nullable == true { sql = sql + " NULL" }
		if column.Nullable == false { sql = sql + " NOT NULL" }

		// UNIQUE
		if column.Unique == true { sql = sql + " UNIQUE" }

		// AUTO INCREMENT
		if column.Increment == true { sql = sql + " AUTO_INCREMENT" }

		// APPEND
		statement = statement + sql + ", "

	}
	return strings.TrimSuffix(statement, ", ")

}


func (schema *Schema) GenerateColumnKeyStatement() string {
	return fmt.Sprintf("")
}