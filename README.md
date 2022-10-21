# KubeProv provision of managed kubernetes clusters easily

## Roadmap 

- [x] Design configuration layer
- [x] Design `kubeprov.ClusterCreate` interface
- [x] Add support of EKS
- [ ] Add support of GKE
- [ ] Add support of LKE
- [ ] Add support of AKE
- [ ] Handle different node configurations/regions/subnets/VPCs
- [ ] Generate kubeconfig of created cluster

## Installing

```
go get github.com/gen1us2k/kubeprov
```

## Creating cluster on EKS

### Using Cli interface

```
kubeprov create --cluster_name example_cluster --region eu-central-1
```

Wait up to 10-15 minutes and it'll create a cluster for you

Here's the sample code to spin up a new EKS cluster

### Using Go code

```
package main

import (
	"log"
	
	"github.com/gen1us2k/kubeprov/pkg/eks"
	"github.com/gen1us2k/kubeprov/pkg/config"
)

func main() {
	e, err := eks.NewEKSClient(&config.Config{Region:"eu-central-1", ClusterName:"example_cluster"})
	if err != nil {
		log.Fatal(err)
	}
	if err := e.ProvisionCluster(); err != nil {
		log.Fatal(err)
	}
}

```

### Getting kubeconfig 

```
aws eks update-kubeconfig --region eu-central-1 --name example_cluster
kubectl get nodes 
```

## Deleting cluster on EKS
### Using CLI

```
kubeprov delete --cluster_name example_cluster --region eu-central-1

```

### Using Go Code

```
package main

import (
	"log"
	
	"github.com/gen1us2k/kubeprov/pkg/eks"
	"github.com/gen1us2k/kubeprov/pkg/config"
)

func main() {
	e, err := eks.NewEKSClient(&config.Config{Region:"eu-central-1", ClusterName:"example_cluster"})
	if err != nil {
		log.Fatal(err)
	}
	if err := e.UnrovisionCluster(); err != nil {
		log.Fatal(err)
	}
}

```
