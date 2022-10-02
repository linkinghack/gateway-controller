package app

import (
	"github.com/linkinghack/gateway-controller/config"
	"github.com/spf13/cobra"
)

var printConfTemplateCmd = &cobra.Command{
	Use:   "printConf",
	Short: "Print a config file template",
	Run: func(cmd *cobra.Command, args []string) {
		config.PrintConfigTemplate(format)
	},
}
var format string

func init() {
	printConfTemplateCmd.Flags().StringVarP(&format, "format", "f", "yaml", "Serialization mode of config: yaml|json|toml")
}
