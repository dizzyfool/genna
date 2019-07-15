package cmd

import (
	"os"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/generators/model"
	"github.com/dizzyfool/genna/generators/model_named"
	"github.com/dizzyfool/genna/generators/search"
	"github.com/dizzyfool/genna/generators/validate"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	root.AddCommand(
		CreateCommand("model", "Basic go-pg model generator", model.New(logger)),
		CreateCommand("search", "Search generator for go-pg models", search.New(logger)),
		CreateCommand("validation", "Validation generator for go-pg models", validate.New(logger)),
		CreateCommand("model-named", "Basic go-pg model generator with named structures", model_named.New(logger)),
	)
}

// Execute runs root cmd
func Execute() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// CreateCommand creates cobra command
func CreateCommand(name, description string, generator base.Gen) *cobra.Command {
	command := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  "",
		Run: func(command *cobra.Command, args []string) {
			logger := generator.Logger()

			if !command.HasFlags() {
				if err := command.Help(); err != nil {
					logger.Error("help not found", zap.Error(err))
				}
				os.Exit(0)
				return
			}

			if err := generator.ReadFlags(command); err != nil {
				logger.Error("read flags error", zap.Error(err))
				return
			}

			if err := generator.Generate(); err != nil {
				logger.Error("generate error", zap.Error(err))
				return
			}
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}

	generator.AddFlags(command)

	return command
}
