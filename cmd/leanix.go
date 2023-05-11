/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/interlynk-io/sbomasm/pkg/assemble"
	"github.com/interlynk-io/sbomasm/pkg/leanix"
	"github.com/interlynk-io/sbomasm/pkg/logger"
	"github.com/spf13/cobra"
)

// leanixCmd represents the leanix command
var leanixCmd = &cobra.Command{
	Use:   "leanix",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitDebugLogger()
		ctx := logger.WithLogger(context.Background())

		token, err := leanix.Authenticate()
		if err != nil {
			panic(err)
		}

		fmt.Println("Authentication token:", token)

		resp, _ := leanix.GetProduct(token, "Interlynk Service")

		files, _ := leanix.WriteTempCDXFiles(resp)
		fmt.Println("Wrote", len(files), "CDX files", files)

		aParams := assemble.NewParams()
		aParams.Name = resp.Data.Products[0].Name
		aParams.Version = "1.0.0"
		aParams.Input = files
		aParams.Ctx = &ctx
		aParams.ConfigPath = "./leanix-config.yaml"
		aParams.Output = "leanix-combined.sbom.json"
		assemble.Assemble(aParams)
	},
}

func init() {
	rootCmd.AddCommand(leanixCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// leanixCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// leanixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
