package types

type (
	Column struct {
		Name      string `mapstructure:"name"`
		Type      string `mapstructure:"type"`
		Nullable  bool   `mapstructure:"nullable"`
		Primary   bool   `mapstructure:"primary"`
		Increment bool   `mapstructure:"increment"`
		Collation string `mapstructure:"collation"`
	}

	Table struct {
		Engine    string `mapstructure:"engine"`
		Charset   string `mapstructure:"charset"`
		Collation string `mapstructure:"collation"`
	}

	Schema struct {
		Table   Table    `mapstructure:"table"`
		Columns []Column `mapstructure:"columns"`
	}
)
