package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/noborus/hyov-forpg/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var query string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hyov-forpg",
		Short: "Terminal viewer for PostgreSQL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if query == "" {
				cmd.Help()
				return nil
			}

			nopager := isNoPager(cmd)
			return internal.Run(query, nopager)
		},
	}
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
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

// isNoPager determines if the output should bypass the pager based on the command flag and terminal status
func isNoPager(cmd *cobra.Command) bool {
	nopager, _ := cmd.Flags().GetBool("no-pager")
	if !isTerminal(os.Stdout.Fd()) {
		nopager = true
	}
	return nopager
}

// isTerminal returns true if the given file descriptor is a terminal
func isTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
