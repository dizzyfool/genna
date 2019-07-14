package validate

import (
	"os"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	pkg = "pkg"
)

// Command gets generator cli command
func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "validate",
		Short: "Validation generator for go-pg models",
		Long:  "",
		Run:   Run,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}

	AddFlags(command)
	return command
}

// Run is callback for command
func Run(cmd *cobra.Command, _ []string) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	if !cmd.HasFlags() {
		if err := cmd.Help(); err != nil {
			logger.Error("help not found", zap.Error(err))
		}
		os.Exit(0)
		return
	}

	conn, options, err := ReadFlags(cmd)
	if err != nil {
		logger.Error("read flags error", zap.Error(err))
		return
	}

	generator := New(conn, logger)

	if err := generator.Generate(options); err != nil {
		logger.Error("generate error", zap.Error(err))
		return
	}
}

// AddFlags adds flags to command
func AddFlags(command *cobra.Command) {
	base.AddBaseFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.StringP(pkg, "p", util.DefaultPackage, "package for model files")
}

// ReadFlags reads flags from user input
func ReadFlags(command *cobra.Command) (conn string, options Options, err error) {
	conn, options.Output, options.Tables, options.FollowFKs, err = base.ReadBaseFlags(command)
	if err != nil {
		return
	}

	flags := command.Flags()

	if options.Package, err = flags.GetString(pkg); err != nil {
		return
	}

	return
}
