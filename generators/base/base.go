package base

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"os"

	"github.com/dizzyfool/genna/lib"
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
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

	// Go-PG version to use
	GoPgVer = "gopg"
)

// Gen is interface for all generators
type Gen interface {
	AddFlags(command *cobra.Command)
	ReadFlags(command *cobra.Command) error

	Generate() error
}

// Packer is a function that compile entities to package
type Packer func(entities []model.Entity) (interface{}, error)

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

	// go-pg version
	GoPgVer int
}

// Def sets default options if empty
func (o *Options) Def() {
	if len(o.Tables) == 0 {
		o.Tables = []string{util.Join(util.PublicSchema, "*")}
	}

	if o.GoPgVer == 0 {
		o.GoPgVer = 10
	}
}

// Generator is base generator used in other generators
type Generator struct {
	genna.Genna
}

// NewGenerator creates generator
func NewGenerator(url string) Generator {
	return Generator{
		Genna: genna.New(url, nil),
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

	flags.IntP(GoPgVer, "g", 10, "specify go-pg version (8, 9 and 10 are supported)\n")

	return
}

// ReadFlags reads basic flags from command
func ReadFlags(command *cobra.Command) (conn, output string, tables []string, followFKs bool, gopgVer int, err error) {
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

	if gopgVer, err = flags.GetInt(GoPgVer); err != nil {
		return
	}

	if gopgVer < 8 && gopgVer > 10 {
		err = errors.Errorf("go-pg version %d not supported", gopgVer)
		return
	}

	return
}

// Generate runs whole generation process
func (g Generator) Generate(tables []string, followFKs, useSQLNulls bool, output, tmpl string, packer Packer, goPGVer int) error {
	entities, err := g.Read(tables, followFKs, useSQLNulls, goPGVer)
	if err != nil {
		return fmt.Errorf("read database error: %w", err)
	}

	return g.GenerateFromEntities(entities, output, tmpl, packer)
}

func (g Generator) GenerateFromEntities(entities []model.Entity, output, tmpl string, packer Packer) error {
	parsed, err := template.New("base").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("parsing template error: %w", err)
	}

	pack, err := packer(entities)
	if err != nil {
		return fmt.Errorf("packing data error: %w", err)
	}

	var buffer bytes.Buffer
	if err := parsed.ExecuteTemplate(&buffer, "base", pack); err != nil {
		return fmt.Errorf("processing model template error: %w", err)
	}

	saved, err := util.FmtAndSave(buffer.Bytes(), output)
	if err != nil {
		if !saved {
			return fmt.Errorf("saving file error: %w", err)
		}
		log.Printf("formatting file %s error: %s", output, err)
	}

	log.Printf("successfully generated %d models", len(entities))

	return nil
}

// CreateCommand creates cobra command
func CreateCommand(name, description string, generator Gen) *cobra.Command {
	command := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  "",
		Run: func(command *cobra.Command, args []string) {
			if !command.HasFlags() {
				if err := command.Help(); err != nil {
					log.Printf("help not found, error: %s", err)
				}
				os.Exit(0)
				return
			}

			if err := generator.ReadFlags(command); err != nil {
				log.Printf("read flags error: %s", err)
				return
			}

			if err := generator.Generate(); err != nil {
				log.Printf("generate error: %s", err)
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
