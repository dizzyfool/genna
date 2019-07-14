package search

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/util"
)

const (
	pkg     = "pkg"
	keepPK  = "keep-pk"
	noAlias = "no-alias"
	relaxed = "relaxed"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "search",
		Short: "Search generator for go-pg model",
		Long:  "",
		Run:   Run,
	}

	AddFlags(command)
	return command
}

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

func AddFlags(command *cobra.Command) {
	base.AddBaseFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.StringP(pkg, "p", util.DefaultPackage, "package for model files")

	flags.Bool(keepPK, false, "keep primary key name as is (by default it should be converted to 'ID')")

	flags.Bool(noAlias, false, `do not set 'alias' tag to "t"`)

	flags.Bool(relaxed, false, "use interface{} type in search filters\n")
}

func ReadFlags(command *cobra.Command) (conn string, options Options, err error) {
	conn, options.Output, options.Tables, options.FollowFKs, err = base.ReadBaseFlags(command)
	if err != nil {
		return
	}

	flags := command.Flags()

	if options.Package, err = flags.GetString(pkg); err != nil {
		return
	}

	if options.KeepPK, err = flags.GetBool(keepPK); err != nil {
		return
	}

	if options.NoAlias, err = flags.GetBool(noAlias); err != nil {
		return
	}

	if options.Relaxed, err = flags.GetBool(relaxed); err != nil {
		return
	}

	return
}
