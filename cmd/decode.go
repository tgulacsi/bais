package cmd

import (
	"github.com/jarrodhroberson/bais/bais"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
		inputFile := args[0]
		outputFile := args[1]
		content, err := ioutil.ReadFile(inputFile)
		if err != nil {
			panic(err)
		}
		decoded, err := bais.Decode(string(content))
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(outputFile, decoded, 0644)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
