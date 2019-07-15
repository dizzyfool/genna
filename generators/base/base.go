package base

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dizzyfool/genna/lib"
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const (
	// Conn is connection string (-c) basic flag
	Conn = "conn"

	// Output is output filename (-o) basic flag
	Output = "output"

	// Tables is basic flag (-t) for tables to generate
	Tables = "tables"

	// FollowFKs is basic flag (-f) for generate foreign keys models for selected tables
	FollowFKs = "follow-fk"
)

// Packer is a function that compile entities to package
type Packer func(entities []model.Entity) interface{}

// Options is common options for all generators
type Options struct {
	// URL connection string
	URL string

	// Output file path
	Output string

	// List of Tables to generate
	// Default []string{"public.*"}
	Tables []string

	// Generate model for foreign keys,
	// even if Tables not listed in Tables param
	// will not generate fks if schema not listed
	FollowFKs bool
}

// Def sets default options if empty
func (o *Options) Def() {
	if len(o.Tables) == 0 {
		o.Tables = []string{util.Join(util.PublicSchema, "*")}
	}
}

// Generator is base generator used in other generators
type Generator struct {
	genna.Genna
}

// NewGenerator creates generator
func NewGenerator(url string, log *zap.Logger) Generator {
	return Generator{
		Genna: genna.New(url, log),
	}
}

// AddFlags adds basic flags to command
func AddFlags(command *cobra.Command) {
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

// ReadFlags reads basic flags from command
func ReadFlags(command *cobra.Command) (conn, output string, tables []string, followFKs bool, err error) {
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

// Generate runs whole generation process
func (g Generator) Generate(tables []string, followFKs, useSQLNulls bool, output, tmpl string, packer Packer) error {
	entities, err := g.Read(tables, followFKs, useSQLNulls)
	if err != nil {
		return xerrors.Errorf("read database error: %w", err)
	}

	parsed, err := template.New("base").Parse(tmpl)
	if err != nil {
		return xerrors.Errorf("parsing template error: %w", err)
	}

	pack := packer(entities)

	var buffer bytes.Buffer
	if err := parsed.ExecuteTemplate(&buffer, "base", pack); err != nil {
		return xerrors.Errorf("processing model template error: %w", err)
	}

	saved, err := util.FmtAndSave(buffer.Bytes(), output)
	if err != nil {
		if !saved {
			return xerrors.Errorf("saving file error: %w", err)
		}
		g.Logger.Error("formatting file error", zap.Error(err), zap.String("file", output))
	}

	g.Logger.Info(fmt.Sprintf("succesfully generated %d models\n", len(entities)))

	return nil
}

// Gen is interface for all generators
type Gen interface {
	Logger() *zap.Logger

	AddFlags(command *cobra.Command)
	ReadFlags(command *cobra.Command) error

	Generate() error
}
