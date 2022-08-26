package eks

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
)

type (
	EKSClient struct {
		ec2 *ec2.EC2
		eks *eks.EKS
		iam *iam.IAM
	}
)

func NewEKSClient() (*EKSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return nil, err
	}
	return &EKSClient{
		ec2: ec2.New(sess),
		eks: eks.New(sess),
		iam: iam.New(sess),
	}, nil
}
func (e *EKSClient) CreateRole(name string) (*iam.Role, error) {
	var managedPolicyArns = []string{
		"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
	}
	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String("{\"Version\": \"2012-10-17\",\"Statement\": [{\"Effect\": \"Allow\",\"Principal\": {\"Service\": \"ec2.amazonaws.com\"},\"Action\": \"sts:AssumeRole\"}]}"),
		Description:              aws.String("Role description"),
		RoleName:                 aws.String("managed-nodegroup-scope"),
	}
	role, err := e.iam.CreateRole(params)
	if err != nil {
		return nil, err
	}

	for _, policy := range managedPolicyArns {
		_, err := e.iam.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: aws.String(policy),
			RoleName:  role.Role.RoleName,
		})
		if err != nil {
			return nil, err
		}
	}
	return role.Role, nil
}

func (e *EKSClient) DescribeRole(name string) (*iam.Role, error) {
	role, err := e.iam.GetRole(&iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return nil, err
	}
	return role.Role, nil
}

func (e *EKSClient) CreateCluster() error {
	vpcID, err := e.GetDefaultVPC()
	if err != nil {
		return err
	}
	subnets, err := e.GetAllSubnets(vpcID)
	if err != nil {
		return err
	}
	secGroups, err := e.GetAllSecurityGroups()
	if err != nil {
		return err
	}
	result, err := e.eks.CreateCluster(&eks.CreateClusterInput{
		ClientRequestToken: aws.String("1d2129a1-3d38-460a-9756-e5b91fddb951"),
		Name:               aws.String("library-created"),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: secGroups,
			SubnetIds:        subnets,
		},
		RoleArn: aws.String("rolearn"),
	})
	fmt.Println(result)
	return err
}
func (e *EKSClient) GetDefaultVPC() (*string, error) {
	vpcResponse, err := e.ec2.DescribeVpcs(new(ec2.DescribeVpcsInput))
	if err != nil {
		return nil, err
	}
	var vpcID *string
	for _, vpc := range vpcResponse.Vpcs {
		if *vpc.IsDefault {
			vpcID = vpc.VpcId
		}
	}
	return vpcID, nil
}
func (e *EKSClient) GetAllSubnets(vpcID *string) ([]*string, error) {
	netResponse, err := e.ec2.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					vpcID,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var subnets []*string
	for _, net := range netResponse.Subnets {
		subnets = append(subnets, net.SubnetId)
	}
	return subnets, nil
}
func (e *EKSClient) GetAllSecurityGroups() ([]*string, error) {
	secResponse, err := e.ec2.DescribeSecurityGroups(new(ec2.DescribeSecurityGroupsInput))
	if err != nil {
		return nil, err
	}
	var secGroups []*string
	for _, group := range secResponse.SecurityGroups {
		secGroups = append(secGroups, group.GroupId)
	}
	return secGroups, nil
}

func (e *EKSClient) WaitClusterUntilAvailable(name string) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			cluster, err := e.eks.DescribeCluster(&eks.DescribeClusterInput{
				Name: aws.String(name),
			})
			if err != nil {
				return err
			}
			fmt.Println(*cluster.Cluster.Status)
			if *cluster.Cluster.Status == "ACTIVE" {
				return nil
			}
		}
	}
}

func (e *EKSClient) CreateNodeGroup(role *iam.Role) error {
	vpcID, err := e.GetDefaultVPC()
	if err != nil {
		return err
	}
	subnets, err := e.GetAllSubnets(vpcID)
	if err != nil {
		return err
	}
	//amiID := "ami-068d000e86e5d6a81"
	nodeGroup, err := e.eks.CreateNodegroup(&eks.CreateNodegroupInput{
		ClusterName:        aws.String("library-created"),
		AmiType:            aws.String("AL2_x86_64"),
		ClientRequestToken: aws.String("asdomasodmaodma"),
		InstanceTypes: []*string{
			aws.String("m5.large"),
		},
		NodegroupName: aws.String("library-created"),
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(3),
			MaxSize:     aws.Int64(3),
			MinSize:     aws.Int64(2),
		},
		Subnets:  subnets,
		NodeRole: role.Arn,
	})
	if err != nil {
		return err
	}
	fmt.Println(nodeGroup)
	return nil
}

func (e *EKSClient) DeleteNodeGroup(role *iam.Role) error {
	_, err := e.eks.DeleteNodegroup(&eks.DeleteNodegroupInput{
		ClusterName:   aws.String("library-created"),
		NodegroupName: aws.String("library-created"),
	})
	return err
}

func (e *EKSClient) DeleteCluster() error {
	_, err := e.eks.DeleteCluster(&eks.DeleteClusterInput{
		Name: aws.String("library-created"),
	})
	return err
}
