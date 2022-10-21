package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kubeprov",
	Short: "A simple managed kubernetes cluster provisioner",
}

func init() {
	rootCmd.PersistentFlags().StringP("region", "r", "AWS region", "")
	rootCmd.PersistentFlags().StringP("cluster_name", "c", "Cluster name", "")
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("cluster_name", rootCmd.PersistentFlags().Lookup("cluster_name"))
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
