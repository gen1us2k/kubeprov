# KubeProv provision managed kubernetes clusters easily

The goal of this project is to create cli/go interfaces to create managed k8s clusters and automate process for end user.

Ideal usecase looks like this

1. Get API access credentials for desired cloud provider
2. Specify region where k8s cluster needs to be configured
3. Call `CreateCluster()` and get kubeconfig in return


## Roadmap 

- [ ] Design configuration layer
- [x] Design `kubeprov.ClusterCreate` interface
- [x] Add support of EKS
- [ ] Add support of GKE
- [ ] Add support of LKE
- [ ] Add support of AKE
- [ ] Handle different node configurations/regions/subnets/VPCs

## Creating cluster on EKS

Here's the sample code to spin up a new EKS cluster

```
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
	
	role, err := eks.CreateRole("managed-nodegroup-scope")
	if err != nil {
		log.Fatal(err)
	}
	err = eks.CreateCluster(role.Arn)
	if err != nil {
		log.Fatal(err)
	}
	err = eks.WaitClusterUntilAvailable("library-created")
	if err != nil {
		log.Fatal(err)
	}
	err = eks.CreateNodeGroup(role)
	if err != nil {
		log.Fatal(err)
	}
}


```

Getting kubconfig 

```
aws eks update-kubeconfig --region eu-central-1 --name library-created
kubectl get nodes 
```

## Deleting cluster on EKS

```
package main

import (
	"log"
  "time"

	"github.com/gen1us2k/kubeprov/pkg/eks"
)

func main() {
	eks, err := eks.NewEKSClient()
	if err != nil {
		log.Fatal(err)
	}
	err = eks.DeleteNodeGroup(nil)
	if err != nil {
		log.Fatal(err)
	}
  time.Sleep(5*time.Minute)
	err = eks.DeleteCluster()
	if err != nil {
		log.Fatal(err)
	}
	err = eks.DeleteRole("managed-nodegroup-scope")
	if err != nil {
		log.Fatal(err)
	}
}
```
