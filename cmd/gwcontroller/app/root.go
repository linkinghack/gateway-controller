package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkinghack/gateway-controller/config"
	"github.com/linkinghack/gateway-controller/pkg/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var GitTag string
var BuildTime string

// serverRootCmd represents the base command when called without any subcommands
var serverRootCmd = &cobra.Command{
	Use:     "gwcontroller",
	Short:   "Generic Gateway Controller",
	Version: fmt.Sprintf("%s|%s", GitTag, BuildTime),
	Long: fmt.Sprintf(`
A Generic Gateway Controller providing simple & efficient APIs for configuring HTTP/TCP/TLS SNI reverse proxy rules 
of underlying gateway software like Envoy.
Version: %s / %s
©️LINKINGHACK · 2022
`, GitTag, BuildTime),
}

// Execute shoule only be called by main()
func Execute() {
	err := serverRootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// adds all child commands to the root command and sets flags appropriately.
func init() {
	// Define flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for this application (serverRootCmd).
	serverRootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Specify location of config file (Default to ./gwcontrol-conf.yaml|json|toml)")

	// Sub-commands
	serverRootCmd.AddCommand(serveCmd)
	serverRootCmd.AddCommand(printConfTemplateCmd)
	serverRootCmd.AddCommand(initDbCmd)

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		pwd, err := filepath.Abs(".")
		cobra.CheckErr(err)

		// Search config in home directory with name "provision-conf.yaml" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("gwcontrol-conf.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stdout, "use config:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintf(os.Stderr, "Warn: config file not found. err=%s\n", err.Error())
	}

	// Merge cmd-line flags to global configs
	viper.BindPFlags(serverRootCmd.Flags())
	viper.BindPFlags(serveCmd.Flags())
	viper.BindPFlags(initDbCmd.Flags())

	// Parse configs to global config object
	refreshConf := func() {
		mainServerConf := config.GlobalConfig{}
		err := viper.Unmarshal(&mainServerConf)
		cobra.CheckErr(err)

		// Save global config
		config.InitGlobalConfig(&mainServerConf)
	}

	// Load configs immediately
	refreshConf()
	log.InitGlobalLogger()

	// Add conf reload watcher
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("Config file change detected: %s of %s", in.Op.String(), in.Name)
		refreshConf()
	})

}
