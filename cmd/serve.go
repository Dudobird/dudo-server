package cmd

import (
	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/routers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the dudo serve",
	Run: func(cmd *cobra.Command, args []string) {
		app := core.NewApp(cfgFile)
		router, err := routers.LoadRouters()
		if err != nil {
			return
		}
		app.Router = router
		app.Run()
	},
}
