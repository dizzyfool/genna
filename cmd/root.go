package cmd

import (
	"os"

	"github.com/dizzyfool/genna/generators/model"
	"github.com/dizzyfool/genna/generators/named"
	"github.com/dizzyfool/genna/generators/search"
	"github.com/dizzyfool/genna/generators/validate"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var root = &cobra.Command{
	Use:   "genna",
	Short: "Genna is model generator for go-pg package",
	Long: `This application is a tool to generate the needed files
to quickly create a models for go-pg https://github.com/go-pg/pg/v9`,
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

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	root.AddCommand(
		model.CreateCommand(logger),
		search.CreateCommand(logger),
		validate.CreateCommand(logger),
		named.CreateCommand(logger),
	)
}

// Execute runs root cmd
func Execute() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
