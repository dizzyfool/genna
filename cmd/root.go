package cmd

import (
	"os"

	"github.com/LdDl/bungen/generators/model"
	"github.com/LdDl/bungen/generators/named"
	"github.com/LdDl/bungen/generators/search"
	"github.com/LdDl/bungen/generators/validate"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "bungen",
	Short: "Bungen is model generator for Bun package [Postgres Driver]",
	Long: `This application is a tool to generate the needed files
to quickly create a models for Bun [Postgres driver] https://github.com/uptrace/bun`,
	Version: "0.1.0",
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
