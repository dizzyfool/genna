package cmd

import (
	"os"

	"github.com/dizzyfool/genna/generators/model"
	"github.com/dizzyfool/genna/generators/named"
	"github.com/dizzyfool/genna/generators/search"
	"github.com/dizzyfool/genna/generators/validate"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "genna",
	Short: "Genna is model generator for go-pg package",
	Long: `This application is a tool to generate the needed files
to quickly create a models for go-pg https://github.com/go-pg/pg/v9`,
	Version: "1.1.7",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic("help not found")
		}
	},
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
}

func init() {
	root.AddCommand(
		model.CreateCommand(),
		search.CreateCommand(),
		validate.CreateCommand(),
		named.CreateCommand(),
	)
}

// Execute runs root cmd
func Execute() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
