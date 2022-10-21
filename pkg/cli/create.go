package cli

import (
	"log"

	"github.com/gen1us2k/kubeprov/pkg/config"
	"github.com/gen1us2k/kubeprov/pkg/eks"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.InitAndParse()
		if err != nil {
			log.Fatal(err)
		}
		eks, err := eks.NewEKSClient(conf)
		if err != nil {
			log.Fatal(err)
		}
		if err := eks.ProvisionCluster(); err != nil {
			log.Fatal(err)
		}
	},
}
