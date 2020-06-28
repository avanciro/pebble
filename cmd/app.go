package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "pebble",
	Short: "Pebble is a simple DB schema management system.",
	Run: func(cmd *cobra.Command, args []string) {},
}


func Dispatch() {

	/*
		Add commands to the command stack of the
		Cobra instance
	*/
	Command.AddCommand(Version_Command)

	Command.AddCommand(Schema_Command)
	Schema_Command.AddCommand(Schema_Sub_Command)


	if err := Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}