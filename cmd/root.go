package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "proxy",
		Short: "A single file reverse proxy",
		Long: `A single file reverse proxy
Example:
  proxy --from 0.0.0.0:8080 --to https://www.google.com`,
		Run: func(cmd *cobra.Command, args []string) {

			from, err := cmd.Flags().GetString("from")
			if err != nil {
				log.Fatalln(err)
			}

			to, err := cmd.Flags().GetString("to")
			if err != nil {
				log.Fatalln(err)
			}

			target, err := url.Parse(to)
			if err != nil {
				log.Fatalln(err)
			}

			targets := []*middleware.ProxyTarget{
				{
					URL: target,
				},
			}

			e := echo.New()
			e.HideBanner = true

			e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

			e.Logger.Fatal(e.Start(from))
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("from", "", "localhost:18080", "Address to listen")
	rootCmd.PersistentFlags().StringP("to", "", "", "Proxy to")
	rootCmd.MarkPersistentFlagRequired("to")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
