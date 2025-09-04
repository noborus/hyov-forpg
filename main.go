package main

import (
	"fmt"
	"os"

	"github.com/noborus/hyov-forpg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var query string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hyov-forpg",
		Short: "Terminal viewer for PostgreSQL",
		Run: func(cmd *cobra.Command, args []string) {
			nopager, _ := cmd.Flags().GetBool("no-pager")
			if query == "" {
				cmd.Help()
				return
			}
			internal.Run(query, nopager)
		},
	}

	rootCmd.Flags().StringP("connection", "c", "", "Database connection string")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "SQL query to execute")
	rootCmd.Flags().BoolP("no-pager", "n", false, "Disable pager (output to stdout)")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/hyov-forpg")
	if err := viper.ReadInConfig(); err != nil {
		// Notify user if config file is not found or cannot be read
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error if not required
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			return
		}
	}
	viper.BindPFlag("db.connection", rootCmd.Flags().Lookup("connection"))
	viper.SetDefault("db.connection", "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	}
}
