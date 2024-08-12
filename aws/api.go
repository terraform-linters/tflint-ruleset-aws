package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// Client is an interface for API client.
// This is primarily used for mock clients.
type Client interface {
	DescribeSecurityGroups() (map[string]bool, error)
	DescribeSubnets() (map[string]bool, error)
	DescribeDBSubnetGroups() (map[string]bool, error)
	DescribeOptionGroups() (map[string]bool, error)
	DescribeDBParameterGroups() (map[string]bool, error)
	DescribeCacheParameterGroups() (map[string]bool, error)
	DescribeCacheSubnetGroups() (map[string]bool, error)
	DescribeInstances() (map[string]bool, error)
	DescribeImages(*ec2.DescribeImagesInput) (map[string]bool, error)
	ListInstanceProfiles() (map[string]bool, error)
	DescribeKeyPairs() (map[string]bool, error)
	DescribeEgressOnlyInternetGateways() (map[string]bool, error)
	DescribeInternetGateways() (map[string]bool, error)
	DescribeNatGateways() (map[string]bool, error)
	DescribeNetworkInterfaces() (map[string]bool, error)
	DescribeRouteTables() (map[string]bool, error)
	DescribeVpcPeeringConnections() (map[string]bool, error)
}

// DescribeSecurityGroups is a wrapper of DescribeSecurityGroups
func (c *AwsClient) DescribeSecurityGroups() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeSecurityGroups(context.Background(), &ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return ret, err
	}
	for _, sg := range resp.SecurityGroups {
		ret[*sg.GroupId] = true
	}
	return ret, err
}

// DescribeSubnets is a wrapper of DescribeSubnets
func (c *AwsClient) DescribeSubnets() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{})
	if err != nil {
		return ret, err
	}
	for _, subnet := range resp.Subnets {
		ret[*subnet.SubnetId] = true
	}
	return ret, err
}

// DescribeDBSubnetGroups is a wrapper of DescribeDBSubnetGroups
func (c *AwsClient) DescribeDBSubnetGroups() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.RDS.DescribeDBSubnetGroups(context.Background(), &rds.DescribeDBSubnetGroupsInput{})
	if err != nil {
		return ret, err
	}
	for _, subnetGroup := range resp.DBSubnetGroups {
		ret[*subnetGroup.DBSubnetGroupName] = true
	}
	return ret, err
}

// DescribeOptionGroups is a wrapper of DescribeOptionGroups
func (c *AwsClient) DescribeOptionGroups() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.RDS.DescribeOptionGroups(context.Background(), &rds.DescribeOptionGroupsInput{})
	if err != nil {
		return ret, err
	}
	for _, optionGroup := range resp.OptionGroupsList {
		ret[*optionGroup.OptionGroupName] = true
	}
	return ret, err
}

// DescribeDBParameterGroups is a wrapper of DescribeDBParameterGroups
func (c *AwsClient) DescribeDBParameterGroups() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.RDS.DescribeDBParameterGroups(context.Background(), &rds.DescribeDBParameterGroupsInput{})
	if err != nil {
		return ret, err
	}
	for _, parameterGroup := range resp.DBParameterGroups {
		ret[*parameterGroup.DBParameterGroupName] = true
	}
	return ret, err
}

// DescribeCacheParameterGroups is a wrapper of DescribeCacheParameterGroups
func (c *AwsClient) DescribeCacheParameterGroups() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.ElastiCache.DescribeCacheParameterGroups(context.Background(), &elasticache.DescribeCacheParameterGroupsInput{})
	if err != nil {
		return ret, err
	}
	for _, parameterGroup := range resp.CacheParameterGroups {
		ret[*parameterGroup.CacheParameterGroupName] = true
	}
	return ret, err
}

// DescribeCacheSubnetGroups is a wrapper of DescribeCacheSubnetGroups
func (c *AwsClient) DescribeCacheSubnetGroups() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.ElastiCache.DescribeCacheSubnetGroups(context.Background(), &elasticache.DescribeCacheSubnetGroupsInput{})
	if err != nil {
		return ret, err
	}
	for _, subnetGroup := range resp.CacheSubnetGroups {
		ret[*subnetGroup.CacheSubnetGroupName] = true
	}
	return ret, err
}

// DescribeInstances is a wrapper of DescribeInstances
func (c *AwsClient) DescribeInstances() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return ret, err
	}
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			ret[*instance.InstanceId] = true
		}
	}
	return ret, err
}

// DescribeImages is a wrapper of DescribeImages
func (c *AwsClient) DescribeImages(in *ec2.DescribeImagesInput) (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeImages(context.Background(), in)
	if err != nil {
		return ret, err
	}
	for _, image := range resp.Images {
		ret[*image.ImageId] = true
	}
	return ret, err
}

// ListInstanceProfiles is a wrapper of ListInstanceProfiles
func (c *AwsClient) ListInstanceProfiles() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.IAM.ListInstanceProfiles(context.Background(), &iam.ListInstanceProfilesInput{})
	if err != nil {
		return ret, err
	}
	for _, iamProfile := range resp.InstanceProfiles {
		ret[*iamProfile.InstanceProfileName] = true
	}
	return ret, err
}

// DescribeKeyPairs is a wrapper of DescribeKeyPairs
func (c *AwsClient) DescribeKeyPairs() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeKeyPairs(context.Background(), &ec2.DescribeKeyPairsInput{})
	if err != nil {
		return ret, err
	}
	for _, keyPair := range resp.KeyPairs {
		ret[*keyPair.KeyName] = true
	}
	return ret, err
}

// DescribeEgressOnlyInternetGateways is wrapper of DescribeEgressOnlyInternetGateways
func (c *AwsClient) DescribeEgressOnlyInternetGateways() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeEgressOnlyInternetGateways(context.Background(), &ec2.DescribeEgressOnlyInternetGatewaysInput{})
	if err != nil {
		return ret, err
	}
	for _, egateway := range resp.EgressOnlyInternetGateways {
		ret[*egateway.EgressOnlyInternetGatewayId] = true
	}
	return ret, err
}

// DescribeInternetGateways is a wrapper of DescribeInternetGateways
func (c *AwsClient) DescribeInternetGateways() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeInternetGateways(context.Background(), &ec2.DescribeInternetGatewaysInput{})
	if err != nil {
		return ret, err
	}
	for _, gateway := range resp.InternetGateways {
		ret[*gateway.InternetGatewayId] = true
	}
	return ret, err
}

// DescribeNatGateways is a wrapper of DescribeNatGateways
func (c *AwsClient) DescribeNatGateways() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeNatGateways(context.Background(), &ec2.DescribeNatGatewaysInput{})
	if err != nil {
		return ret, err
	}
	for _, ngateway := range resp.NatGateways {
		ret[*ngateway.NatGatewayId] = true
	}
	return ret, err
}

// DescribeNetworkInterfaces is a wrapper of DescribeNetworkInterfaces
func (c *AwsClient) DescribeNetworkInterfaces() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeNetworkInterfaces(context.Background(), &ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		return ret, err
	}
	for _, networkInterface := range resp.NetworkInterfaces {
		ret[*networkInterface.NetworkInterfaceId] = true
	}
	return ret, err
}

// DescribeRouteTables is a wrapper of DescribeRouteTables
func (c *AwsClient) DescribeRouteTables() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeRouteTables(context.Background(), &ec2.DescribeRouteTablesInput{})
	if err != nil {
		return ret, err
	}
	for _, routeTable := range resp.RouteTables {
		ret[*routeTable.RouteTableId] = true
	}
	return ret, err
}

// DescribeVpcPeeringConnections is a wrapper of DescribeVpcPeeringConnections
func (c *AwsClient) DescribeVpcPeeringConnections() (map[string]bool, error) {
	ret := map[string]bool{}
	resp, err := c.EC2.DescribeVpcPeeringConnections(context.Background(), &ec2.DescribeVpcPeeringConnectionsInput{})
	if err != nil {
		return ret, err
	}
	for _, vpcPeeringConnection := range resp.VpcPeeringConnections {
		ret[*vpcPeeringConnection.VpcPeeringConnectionId] = true
	}
	return ret, err
}
