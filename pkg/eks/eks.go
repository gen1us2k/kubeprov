package eks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gen1us2k/kubeprov/pkg/config"
	"github.com/google/uuid"
)

type (
	EKSClient struct {
		ec2  *ec2.EC2
		eks  *eks.EKS
		iam  *iam.IAM
		conf *config.Config
	}
)

func NewEKSClient(c *config.Config) (*EKSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
	if err != nil {
		return nil, err
	}
	return &EKSClient{
		ec2:  ec2.New(sess),
		eks:  eks.New(sess),
		iam:  iam.New(sess),
		conf: c,
	}, nil
}
func (e *EKSClient) CreateRole() (*iam.Role, error) {
	var managedPolicyArns = []string{
		"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
		"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
	}
	params := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String("{\"Version\": \"2012-10-17\",\"Statement\": [{\"Effect\": \"Allow\",\"Principal\": {\"Service\": \"ec2.amazonaws.com\"},\"Action\": \"sts:AssumeRole\"}, {\"Effect\": \"Allow\",\"Principal\": {\"Service\": \"eks.amazonaws.com\"},\"Action\": \"sts:AssumeRole\"}]}"),
		Description:              aws.String("Kubeprov library role"),
		RoleName:                 aws.String(e.conf.RoleName),
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
func (e *EKSClient) DeleteRole() error {
	policies, err := e.iam.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(e.conf.RoleName),
	})
	if err != nil {
		return err
	}
	for _, policy := range policies.AttachedPolicies {
		_, err := e.iam.DetachRolePolicy(&iam.DetachRolePolicyInput{
			RoleName:  aws.String(e.conf.RoleName),
			PolicyArn: policy.PolicyArn,
		})
		if err != nil {
			return err
		}
	}
	params := &iam.DeleteRoleInput{
		RoleName: aws.String(e.conf.RoleName),
	}
	_, err = e.iam.DeleteRole(params)
	return err
}

func (e *EKSClient) DescribeRole(name string) (*iam.Role, error) {
	role, err := e.iam.GetRole(&iam.GetRoleInput{RoleName: aws.String(name)})
	if err != nil {
		return nil, err
	}
	return role.Role, nil
}

func (e *EKSClient) CreateCluster(roleArn *string) error {
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
	reqToken, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	result, err := e.eks.CreateCluster(&eks.CreateClusterInput{
		ClientRequestToken: aws.String(reqToken.String()),
		Name:               aws.String(e.conf.ClusterName),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: secGroups,
			SubnetIds:        subnets,
		},
		RoleArn: roleArn,
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

func (e *EKSClient) WaitClusterUntilAvailable() error {
	return e.eks.WaitUntilClusterActive(&eks.DescribeClusterInput{
		Name: aws.String(e.conf.ClusterName),
	})
}
func (e *EKSClient) WaitClusterUntilDeleted() error {
	return e.eks.WaitUntilClusterDeleted(&eks.DescribeClusterInput{
		Name: aws.String(e.conf.ClusterName),
	})
}
func (e *EKSClient) WaitUntilNodegroupActive() error {
	return e.eks.WaitUntilNodegroupActive(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(e.conf.ClusterName),
		NodegroupName: aws.String(e.conf.NodegroupName()),
	})
}
func (e *EKSClient) WaitUntilNodegroupDeleted() error {
	return e.eks.WaitUntilNodegroupDeleted(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(e.conf.ClusterName),
		NodegroupName: aws.String(e.conf.NodegroupName()),
	})
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
	reqToken, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	_, err = e.eks.CreateNodegroup(&eks.CreateNodegroupInput{
		ClusterName:        aws.String(e.conf.ClusterName),
		AmiType:            aws.String(e.conf.AMIType),
		ClientRequestToken: aws.String(reqToken.String()),
		InstanceTypes: []*string{
			aws.String(e.conf.InstanceType),
		},
		NodegroupName: aws.String(e.conf.NodegroupName()),
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(e.conf.DesiredState),
			MaxSize:     aws.Int64(e.conf.MaxSize),
			MinSize:     aws.Int64(e.conf.MinSize),
		},
		Subnets:  subnets,
		NodeRole: role.Arn,
	})
	if err != nil {
		return err
	}
	return nil
}

func (e *EKSClient) DeleteNodeGroup() error {
	_, err := e.eks.DeleteNodegroup(&eks.DeleteNodegroupInput{
		ClusterName:   aws.String(e.conf.ClusterName),
		NodegroupName: aws.String(e.conf.NodegroupName()),
	})
	return err
}

func (e *EKSClient) DeleteCluster() error {
	_, err := e.eks.DeleteCluster(&eks.DeleteClusterInput{
		Name: aws.String(e.conf.ClusterName),
	})
	return err
}

func (e *EKSClient) ProvisionCluster() error {
	role, err := e.CreateRole()
	if err != nil {
		return err
	}
	err = e.CreateCluster(role.Arn)
	if err != nil {
		return err
	}
	err = e.WaitClusterUntilAvailable()
	if err != nil {
		return err
	}
	err = e.CreateNodeGroup(role)
	if err != nil {
		return err
	}
	err = e.WaitUntilNodegroupActive()
	if err != nil {
		return err
	}
	return nil
}

func (e *EKSClient) UnprovisionCluster() error {
	err := e.DeleteNodeGroup()
	if err != nil {
		return err
	}
	err = e.WaitUntilNodegroupDeleted()
	if err != nil {
		return err
	}
	err = e.DeleteCluster()
	if err != nil {
		return err
	}
	err = e.WaitClusterUntilDeleted()
	if err != nil {
		return err
	}
	err = e.DeleteRole()
	if err != nil {
		return err
	}
	return nil
}
