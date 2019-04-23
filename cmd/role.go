package cmd

import (
	"os"

	"github.com/Dudobird/dudo-server/core"
	"github.com/Dudobird/dudo-server/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	email    string
	password string
)

func init() {
	rootCmd.AddCommand(roleadmin)
	roleadmin.Flags().StringVarP(&email, "email", "u", "", "email for admin user")
	roleadmin.Flags().StringVarP(&password, "password", "p", "", "password for admin user")
}

var roleadmin = &cobra.Command{
	Use:   "role",
	Short: "create a new admin account",
	Run: func(cmd *cobra.Command, args []string) {
		_ = core.NewApp(cfgFile)
		if email == "" || password == "" {
			log.Errorln("email and password can not be nil")
			os.Exit(1)
		}
		err := models.InsertAdminUser(email, password)
		if err != nil {
			log.Errorln("insert admin account fail with error:" + err.Error())
			os.Exit(1)
		}
		log.Infof("insert admin account success")
	},
}
