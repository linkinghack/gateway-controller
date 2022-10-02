package app

import (
	"github.com/linkinghack/gateway-controller/pkg/dbconn"
	"github.com/linkinghack/gateway-controller/pkg/log"
	"github.com/spf13/cobra"
)

var initDbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "initiate db tables or do a migration",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.GetGlobalLogger()

		//  auto create tables
		// TODO: check here every time update DB data model
		err := dbconn.GetDBConn().AssureTables()
		if err != nil {
			logger.
				WithField("err", err.Error()).
				Fatalf("Tables migration failed.")
			return
		}
		logger.Infof("DB auto migrate successfully.")
	},
}
