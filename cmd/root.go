// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/dizzyfool/genna/database"
	"github.com/dizzyfool/genna/generator"
	"github.com/dizzyfool/genna/model"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	conn      = "conn"
	out       = "out"
	impPath   = "import"
	pkg       = "pkg"
	schemaPkg = "schema-pkg"
	multiFile = "multi-file"
	tables    = "tables"
	view      = "view"
	followFK  = "follow-fk"
	keepPK    = "keep-pk"
	noDiscard = "no-discard"
	noAlias   = "no-alias"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "genna",
	Short: "Genna is model generator for go-pg package",
	Long: `This application is a tool to generate the needed files
to quickly create a models for go-pg https://github.com/go-pg/pg`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.HasFlags() {
			if err := cmd.Help(); err != nil {
				panic("help not found")
			}
			os.Exit(0)
		}

		flags := cmd.Flags()

		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.Encoding = "console"
		logger, err := config.Build()
		if err != nil {
			panic(err)
		}

		url, err := flags.GetString(conn)
		if err != nil {
			panic(err)
		}

		db, err := database.NewDatabase(url, logger)
		if err != nil {
			panic(err)
		}

		store := database.NewStore(db)
		options, err := flagsToOptions(flags)
		if err != nil {
			panic(err)
		}

		fmt.Println("getting info from database ...")
		tables, err := store.Tables(model.Schemas(options.Tables))
		if err != nil {
			panic(err)
		}

		fmt.Println("running generator ...")
		result, err := generator.NewGenerator(options, logger).Process(tables)
		if err != nil {
			panic(err)
		}

		fmt.Printf("generated %d models from %d tables in total. saved in %d packages in %d files\n",
			result.GeneratedModels, result.TotalTables, result.GeneratedPackages, result.GeneratedFiles,
		)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()

	flags.SortFlags = false

	flags.StringP(conn, "c", "", "connection string to your postgres database")
	if err := rootCmd.MarkFlagRequired(conn); err != nil {
		panic(err)
	}

	flags.StringP(out, "o", "", "output directory for generated files")
	if err := rootCmd.MarkFlagRequired(out); err != nil {
		panic(err)
	}

	flags.StringSliceP(tables, "t", []string{"public.*"}, "table names for model generation separeted by comma\nuse 'schema_name.*' to generate model for every table in model")
	flags.StringP(impPath, "i", "", "prefix for imports for generated models")

	flags.StringP(pkg, "p", model.DefaultPackage, "package for model files. ignored with --multi-file param")
	flags.BoolP(schemaPkg, "s", false, "generate every schema as separate package")
	flags.BoolP(multiFile, "m", false, `generate one file for package
schema-pkg | multi-file | result
true       | false      | each generated package will contain one file
true       | true       | each generated package will contain several files, one per model
false      | true       | one package for all models separated to different files
false      | false      | one big file for all models
`)

	flags.BoolP(view, "v", false, "use view for selects e.g. getUsers for users table")
	flags.BoolP(followFK, "f", false, "generate models for foreign keys, even if it not listed in tables\n")

	flags.Bool(keepPK, false, "keep primary key name as is (by default it should be converted to 'ID') \n")
	flags.Bool(noAlias, false, `do set 'alias' tag to "t"`)
	flags.Bool(noDiscard, false, "do not use 'discard_unknown_columns' tag")
}

func flagsToOptions(flags *pflag.FlagSet) (generator.Options, error) {
	var err error

	options := generator.Options{}

	if options.Output, err = flags.GetString(out); err != nil {
		return options, err
	}

	if options.ImportPath, err = flags.GetString(impPath); err != nil {
		return options, err
	}

	if options.Package, err = flags.GetString(pkg); err != nil {
		return options, err
	}

	if options.Tables, err = flags.GetStringSlice(tables); err != nil {
		return options, err
	}

	if options.SchemaPackage, err = flags.GetBool(schemaPkg); err != nil {
		return options, err
	}

	if options.MultiFile, err = flags.GetBool(multiFile); err != nil {
		return options, err
	}

	if options.View, err = flags.GetBool(view); err != nil {
		return options, err
	}

	if options.FollowFKs, err = flags.GetBool(followFK); err != nil {
		return options, err
	}

	if options.KeepPK, err = flags.GetBool(keepPK); err != nil {
		return options, err
	}

	if options.NoDiscard, err = flags.GetBool(noDiscard); err != nil {
		return options, err
	}

	if options.NoAlias, err = flags.GetBool(noAlias); err != nil {
		return options, err
	}

	return options, nil
}
