package main

import (
	"log"

	"github.com/gen1us2k/kubeprov/pkg/eks"
)

func main() {
	eks, err := eks.NewEKSClient()
	if err != nil {
		log.Fatal(err)
	}
	// err = eks.DeleteNodeGroup(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err = eks.DeleteCluster()
	if err != nil {
		log.Fatal(err)
	}
	err = eks.DeleteRole("managed-nodegroup-scope")
	if err != nil {
		log.Fatal(err)
	}
	// role, err := eks.CreateRole("managed-nodegroup-scope")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = eks.CreateCluster(role.Arn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = eks.WaitClusterUntilAvailable("library-created")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = eks.CreateNodeGroup(role)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
