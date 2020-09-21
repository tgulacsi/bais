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
		var inputFile *os.File
		var outputFile *os.File
		if len(args) == 0 {
			inputFile = os.Stdin
			outputFile = os.Stdout
		}
		if len(args) >= 1 {
			inf, err := os.Open(args[0])
			if err != nil {
				panic(err)
			}
			inputFile = inf
		}
		if len(args) == 2 {
			outf, err := os.Open(args[0])
			if err != nil {
				panic(err)
			}
			outputFile = outf
		}
		content, err := ioutil.ReadAll(inputFile)
		if err != nil {
			panic(err)
		}
		decoded, err := bais.Decode(string(content))
		if err != nil {
			panic(err)
		}
		_, err = outputFile.Write(decoded)
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
