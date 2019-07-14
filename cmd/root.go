package cmd

import (
	"os"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/generators/search"
	"github.com/dizzyfool/genna/generators/validate"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "genna",
	Short: "Genna is model generator for go-pg package",
	Long: `This application is a tool to generate the needed files
to quickly create a models for go-pg https://github.com/go-pg/pg`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.HasSubCommands() {
			if err := cmd.Help(); err != nil {
				panic("help not found")
			}
			os.Exit(0)
		}
	},
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
}

// Execute runs root cmd
func Execute() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(
		base.Command("base", "Basic go-pg model generator"),
		search.Command(),
		validate.Command(),
	)
}
