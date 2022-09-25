package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().StringVarP(&(conf.Host), "host", "h", "127.0.0.1", "Server hostname")
	rootCmd.Flags().IntVarP(&(conf.Port), "port", "p", 6379, "Server port")
	rootCmd.Flags().StringVar(&(conf.Password), "password", "", "Password to use when connecting to the server.")
	rootCmd.Flags().Bool("help", false, "Print usage")

}

type Config struct {
	Host     string
	Port     int
	Password string
}

var conf Config

func (c *Config) Prompt() string {
	return fmt.Sprintf("%s:%d> ", c.Host, c.Port)
}

var rootCmd = &cobra.Command{
	Use:   "rcli",
	Short: "redis cli",
	Long:  `redis cli in go`,
	Run: func(cmd *cobra.Command, args []string) {
		initRedis(conf)
		repl()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
