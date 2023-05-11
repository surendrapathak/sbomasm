/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/interlynk-io/sbomasm/pkg/assemble"
	"github.com/interlynk-io/sbomasm/pkg/leanix"
	"github.com/interlynk-io/sbomasm/pkg/logger"
	"github.com/spf13/cobra"
)

// leanixCmd represents the leanix command
var leanixCmd = &cobra.Command{
	Use:   "leanix",
	Short: "Assembled products for leanix",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logger.InitDebugLogger()
		} else {
			logger.InitProdLogger()
		}
		ctx := logger.WithLogger(context.Background())

		lParams, _ := extractArgs(cmd, args)
		lParams.Ctx = &ctx
		lParams.Version = "1.0.0"

		files := initializeLeanIXProduct(lParams.Name)
		lParams.Input = files

		err := assemble.Assemble(lParams)
		if err != nil {
			panic(err)
		}
	},
}

func initializeLeanIXProduct(name string) []string {
	token, err := leanix.Authenticate()
	if err != nil {
		panic(err)
	}
	resp, err := leanix.GetProduct(token, name)
	if err != nil {
		panic(err)
	}

	files, err := leanix.WriteTempCDXFiles(resp)
	if err != nil {
		panic(err)
	}

	return files
}

func init() {
	rootCmd.AddCommand(leanixCmd)
	leanixCmd.Flags().StringP("output", "o", "", "path to assembled sbom, defaults to stdout")
	leanixCmd.Flags().StringP("configPath", "c", "", "path to config file")

	leanixCmd.Flags().StringP("name", "n", "", "Name of product on leanix dashboard")
	leanixCmd.MarkFlagRequired("name")

	leanixCmd.Flags().BoolP("flatMerge", "f", false, "flat merge")
	leanixCmd.Flags().BoolP("hierMerge", "m", true, "hierarchical merge")

	leanixCmd.Flags().BoolP("xml", "x", false, "output in xml format")
	leanixCmd.Flags().BoolP("json", "j", true, "output in json format")

	leanixCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")
}
