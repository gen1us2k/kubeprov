# KubeProv provision managed kubernetes clusters easily

The goal of this project is to create cli/go interfaces to create managed k8s clusters and automate process for end user.

Ideal usecase looks like this

1. Get API access credentials for desired cloud provider
2. Specify region where k8s cluster needs to be configured
3. Call `CreateCluster()` and get kubeconfig in return


## Roadmap 

- [ ] Design configuration layer
- [ ] Design `kubeprov.ClusterCreate` interface
- [ ] Add support of EKS
- [ ] Add support of GKE
- [ ] Add support of LKE
- [ ] Add support of AKE
- [ ] Handle different node configurations/regions/subnets/VPCs
