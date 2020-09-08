package cmd

import (
	"fmt"
	"jarvis/api"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().Int("port", 8080, "API port number")
	viper.BindPFlag("api.port", serverCmd.PersistentFlags().Lookup("port"))

	serverCmd.PersistentFlags().String("log_level", "info", "API log level")
	viper.BindPFlag("api.log_level", serverCmd.PersistentFlags().Lookup("log_level"))
	log.SetOutput(os.Stdout)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "REST API server",
	Args:  cobra.ArbitraryArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := setLogLevel()
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logrus.New()
		api := api.InitApi(logger, viper.GetInt("api.port"))

		return api.Run()
	},
}

func setLogLevel() error {
	var level log.Level
	switch m := viper.GetString("api.log_level"); m {
	case "trace":
		level = log.TraceLevel
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel

	default:
		return fmt.Errorf("Unknown log level: %s", m)
	}

	log.SetLevel(level)
	return nil
}
