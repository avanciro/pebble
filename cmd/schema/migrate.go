package schema

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"pebble/config"
	"pebble/types"
	logger "pebble/utils/log"
	parser "pebble/utils/parser/schema"
	"strings"

	schemalex "github.com/schemalex/schemalex/diff"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Schema_Migrate_Command = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate schema",
	Run:   schema_migrate,
}

/**
* Struct to store flags that passed to the command
* and later assign them to the initialized values
 */
type Flags struct {
	Connection string
}

var flags Flags

// INIT
func init() {
	Schema_Migrate_Command.Flags().StringVarP(&flags.Connection, "connection", "c", "", "")
}

func schema_migrate(cmd *cobra.Command, args []string) {

	/**
	* Check if the migrating is running for default connection or
	* or the specified connection name.
	 */
	var Config = config.Config()
	var ConnectionInformation config.Connection
	if len(flags.Connection) >= 1 {
		var is_connection_found bool
		for _, Connection := range Config.Connections {
			if Connection.Alias == flags.Connection {
				ConnectionInformation = Connection
				is_connection_found = true
			}
		}
		if !is_connection_found {
			fmt.Println("Defined connection not found")
			return
		}
	} else {
		var is_connection_found bool
		for _, Connection := range Config.Connections {
			if Connection.Default {
				ConnectionInformation = Connection
				is_connection_found = true
			}
		}
		if !is_connection_found {
			fmt.Println("Default connection not defined")
			return
		}
	}

	/**
	* Establish new database connection to the database
	* using configuration information provided in the
	* main configuration file
	 */
	dialect := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", ConnectionInformation.User, ConnectionInformation.Password, ConnectionInformation.Host, ConnectionInformation.Port, ConnectionInformation.Name)
	db, err := sql.Open(ConnectionInformation.Driver, dialect)
	if err != nil {
		fmt.Println("Database connection failed")
		return
	}
	defer db.Close()

	/**
	* We need to iterate through every file in the
	* migration directory to find every migration answer
	* file to operate.
	 */
	var schemas []string
	files, err := ioutil.ReadDir("./" + Config.Schema.Directory)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".yml") {
			schemas = append(schemas, strings.Replace(file.Name(), ".yml", "", -1))
		}
	}

	/**
	* We need to drop tables that are not in our schema
	* set first to clear out the unwanted tables from
	* the database.
	 */
	query := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME NOT IN ('%s')", ConnectionInformation.Name, strings.Join(schemas, "','"))
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	/**
	* We have the names of the tables that has no schema
	* file present in the migration store. We can safely
	* drop those tables from the database.
	 */
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			log.Fatal(err)
		}
		logger.Println("green", "[TABLE]: ", "DROPPING ( "+table+" )")
		db.Exec("DROP TABLE " + table)
	}

	/**
	* Loop over every schema file we found in the schema directory
	* to do the needful operations in migration schema
	 */
	for _, file := range schemas {
		fmt.Println("[TABLE]: " + file)

		v := viper.New()
		v.SetConfigName(file)
		v.SetConfigType("yml")
		v.AddConfigPath("./" + Config.Schema.Directory)
		err := v.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		var Schema types.Schema
		v.Unmarshal(&Schema)

		/**
		* Check if the table is exists or not and then create
		* the table if it's not exists first.
		 */
		var count int
		query := fmt.Sprintf("SELECT CAST(COUNT(TABLE_NAME) AS UNSIGNED) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", ConnectionInformation.Name, file)
		db.QueryRow(query).Scan(&count)

		/**
		* We need out yml schema parser to parse the provided file
		* and output the SQL CREATE TABLE format.
		 */
		Parser := parser.Schema{}
		Parser.File("./" + Config.Schema.Directory + "/" + file + ".yml")

		/**
		* Here we handle the table not exists state by using the
		* table count from the last sql query and then we create new
		* table in the database according to the migration file.
		 */
		if count == 0 {
			_, err := db.Exec(Parser.Statement())
			if err != nil {
				fmt.Println(err)
			}
		}

		/**
		* Depending on the recent sql query table count details we can
		* decide to modify exsisting table schema according to the
		* changes in the migration file.
		 */
		if count >= 1 {

			result := []string{"table", "ddl"}
			db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", file)).Scan(&result[0], &result[1])

			statement := &bytes.Buffer{}
			err := schemalex.Strings(statement, result[1], Parser.Statement())
			if err != nil {
				log.Fatal(err)
			}

			/**
			* Split whole statement into seperate query and execute one
			* by one to migrate the schema and report errors if anything
			* happen during the process
			 */
			for _, query := range strings.Split(statement.String(), ";") {
				if len(query) >= 1 {
					_, err = db.Exec(query)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

		}

	}

}
