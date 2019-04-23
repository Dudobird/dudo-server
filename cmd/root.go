package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "config file for dudo appliction")
}

var rootCmd = &cobra.Command{
	Use:   "dudo",
	Short: "dudo is a saas production for file storage and share",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dudo help to show all the command")
	},
}

// Execute for cmd root
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panicf("Program exit with error : %s\n", err)
	}
}
