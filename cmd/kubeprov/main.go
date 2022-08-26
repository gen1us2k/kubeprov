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
	err = eks.DeleteCluster()
	if err != nil {
		log.Fatal(err)
	}
	//role, err := eks.DescribeRole("managed-nodegroup-scope")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = eks.DeleteNodeGroup(role)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// err = eks.CreateCluster()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
