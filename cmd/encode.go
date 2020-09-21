package cmd

import (
	"github.com/jarrodhroberson/bais/bais"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "bais encode",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
		inputFile := args[0]
		outputFile := args[1]
		content, err := ioutil.ReadFile(inputFile)
		if err != nil {
			panic(err)
		}
		encoded := bais.Encode(&content, viper.GetBool("allow-control-characters"))
		err = ioutil.WriteFile(outputFile, []byte(encoded), 0644)
		if err != nil {
			panic(err)
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	encodeCmd.Flags().BoolP("allow-control-characters", "a", false, "Allow control characters to be passed through as is.")
}
