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
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/dizzyfool/genna/database"
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

		// TODO dig container
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		url, err := flags.GetString("conn")
		if err != nil {
			panic(err)
		}

		db, err := database.NewDatabase(url, logger)
		if err != nil {
			panic(err)
		}

		store := database.NewStore(db)

		// TODO move away from here
		tables, err := flags.GetStringSlice("tables")
		if err != nil {
			panic(err)
		}

		info, err := store.Tables(database.Schemas(tables))
		if err != nil {
			panic(err)
		}

		if b, err := json.MarshalIndent(info, "", "\t"); err == nil {
			logger.Info(string(b))
		}
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

	flags.StringP("conn", "c", "", "connection string")
	if err := rootCmd.MarkFlagRequired("conn"); err != nil {
		panic(err)
	}

	flags.StringSliceP("tables", "t", []string{"public.*"}, "table names for model generation separeted by comma\nuse 'schema_name.*' to generate model for every table in model")

	flags.StringP("package", "p", "model", "package for model files")

	flags.BoolP("view", "v", false, "use view for selects e.g. getUsers for users table")
	flags.BoolP("fk", "f", false, "generate models for foreign keys, even if it not listed in tables\n")

	// TODO implement that!
	//flags.StringSliceP("arrays", "a", []string{}, "pg json fields should be parsed as arrays separated by comma")
	//flags.StringSliceP("maps", "m", []string{}, "pg json fields should be parsed as maps separated by comma")
	//flags.StringToStringP("structs", "s", nil, "pg json fields should be parsed as structs, e.g. -s location=LocationStruct,settings=Settings")

	//flags.BoolP("hooks", "k", false, "generate hooks to fill foreign keys after insert/update\nwarning: may not work with recursive relations")
	//flags.Bool("keep-pk", false, "keep primary key name as is (by default it should be converted to 'ID')\nwarning: may break some go-pg features like many-to-many table relations")
	flags.Bool("no-discard", false, "do not use 'discard_unknown_columns' tag\nwarning: may break incomplete models")
}
