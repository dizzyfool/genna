package base

import (
	"os"

	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	Conn      = "conn"
	Output    = "output"
	Tables    = "tables"
	FollowFKs = "follow-fk"

	pkg          = "pkg"
	keepPK       = "keep-pk"
	noDiscard    = "no-discard"
	noAlias      = "no-alias"
	withSearch   = "with-search"
	strictSearch = "strict-search"
	softDelete   = "soft-delete"
	validator    = "validator"
)

func Command(name, description string) *cobra.Command {
	command := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  "",
		Run:   Run,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
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
	AddBaseFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.StringP(pkg, "p", util.DefaultPackage, "package for model files")

	flags.Bool(keepPK, false, "keep primary key name as is (by default it should be converted to 'ID')")
	flags.String(softDelete, "", "field for soft_delete tag\n")

	flags.Bool(noAlias, false, `do not set 'alias' tag to "t"`)
	flags.Bool(noDiscard, false, "do not use 'discard_unknown_columns' tag\n")

	flags.BoolP(withSearch, "s", false, "generate search filters")
	flags.Bool(strictSearch, false, "use exact type (with pointer) in search filters\n")

	flags.Bool(validator, false, "generate validator functions")
}

func AddBaseFlags(command *cobra.Command) {
	flags := command.Flags()

	flags.StringP(Conn, "c", "", "connection string to your postgres database")
	if err := command.MarkFlagRequired(Conn); err != nil {
		panic(err)
	}

	flags.StringP(Output, "o", "", "output file name")
	if err := command.MarkFlagRequired(Output); err != nil {
		panic(err)
	}

	flags.StringSliceP(Tables, "t", []string{"public.*"}, "table names for model generation separated by comma\nuse 'schema_name.*' to generate model for every table in model")
	flags.BoolP(FollowFKs, "f", false, "generate models for foreign keys, even if it not listed in Tables")

	return
}

func ReadFlags(command *cobra.Command) (conn string, options Options, err error) {
	conn, options.Output, options.Tables, options.FollowFKs, err = ReadBaseFlags(command)
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

	if options.SoftDelete, err = flags.GetString(softDelete); err != nil {
		return
	}

	if options.NoDiscard, err = flags.GetBool(noDiscard); err != nil {
		return
	}

	if options.NoAlias, err = flags.GetBool(noAlias); err != nil {
		return
	}

	return
}

func ReadBaseFlags(command *cobra.Command) (conn, output string, tables []string, followFKs bool, err error) {
	flags := command.Flags()

	if conn, err = flags.GetString(Conn); err != nil {
		return
	}

	if output, err = flags.GetString(Output); err != nil {
		return
	}

	if tables, err = flags.GetStringSlice(Tables); err != nil {
		return
	}

	if followFKs, err = flags.GetBool(FollowFKs); err != nil {
		return
	}

	return
}
