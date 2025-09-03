package main

import (
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
	viper.ReadInConfig()
	viper.BindPFlag("db.connection", rootCmd.Flags().Lookup("connection"))
	viper.SetDefault("db.connection", "")

	_ = rootCmd.Execute()
}
