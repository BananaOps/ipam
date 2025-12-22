package aws

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// AWSConfig represents AWS configuration
type AWSConfig struct {
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

// Client represents an AWS client
type Client struct {
	ec2Client *ec2.Client
	config    AWSConfig
}

// VPCInfo represents VPC information
type VPCInfo struct {
	ID        string
	CIDR      string
	Name      string
	Region    string
	IsDefault bool
	Tags      map[string]string
}

// SubnetInfo represents subnet information
type SubnetInfo struct {
	ID               string
	CIDR             string
	Name             string
	VPCId            string
	AvailabilityZone string
	Region           string
	IsPublic         bool
	Tags             map[string]string
}

// NewClient creates a new AWS client
func NewClient(ctx context.Context, awsConfig AWSConfig) (*Client, error) {
	var cfg aws.Config
	var err error

	if awsConfig.AccessKeyID != "" && awsConfig.SecretAccessKey != "" {
		// Use static credentials (not recommended for production)
		log.Printf("Using static credentials for AWS authentication in region: %s", awsConfig.Region)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(awsConfig.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				awsConfig.AccessKeyID,
				awsConfig.SecretAccessKey,
				"",
			)),
		)
	} else {
		// Use default credential chain (environment variables, IAM roles, IRSA, etc.)
		// This automatically handles IRSA in EKS environments, IAM roles on EC2, etc.
		log.Printf("Using default credential chain for AWS authentication in region: %s", awsConfig.Region)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(awsConfig.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Client{
		ec2Client: ec2.NewFromConfig(cfg),
		config:    awsConfig,
	}, nil
}

// ListVPCs retrieves all VPCs in the configured region
func (c *Client) ListVPCs(ctx context.Context) ([]VPCInfo, error) {
	input := &ec2.DescribeVpcsInput{}

	result, err := c.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe VPCs: %w", err)
	}

	var vpcs []VPCInfo
	for _, vpc := range result.Vpcs {
		vpcInfo := VPCInfo{
			ID:        aws.ToString(vpc.VpcId),
			CIDR:      aws.ToString(vpc.CidrBlock),
			Region:    c.config.Region,
			IsDefault: aws.ToBool(vpc.IsDefault),
			Tags:      make(map[string]string),
		}

		// Extract name from tags
		for _, tag := range vpc.Tags {
			if aws.ToString(tag.Key) == "Name" {
				vpcInfo.Name = aws.ToString(tag.Value)
			}
			vpcInfo.Tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}

		// If no name tag, use VPC ID as name
		if vpcInfo.Name == "" {
			vpcInfo.Name = vpcInfo.ID
		}

		vpcs = append(vpcs, vpcInfo)
	}

	return vpcs, nil
}

// ListSubnets retrieves all subnets in the configured region
func (c *Client) ListSubnets(ctx context.Context) ([]SubnetInfo, error) {
	input := &ec2.DescribeSubnetsInput{}

	result, err := c.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe subnets: %w", err)
	}

	var subnets []SubnetInfo
	for _, subnet := range result.Subnets {
		subnetInfo := SubnetInfo{
			ID:               aws.ToString(subnet.SubnetId),
			CIDR:             aws.ToString(subnet.CidrBlock),
			VPCId:            aws.ToString(subnet.VpcId),
			AvailabilityZone: aws.ToString(subnet.AvailabilityZone),
			Region:           c.config.Region,
			IsPublic:         aws.ToBool(subnet.MapPublicIpOnLaunch),
			Tags:             make(map[string]string),
		}

		// Extract name from tags
		for _, tag := range subnet.Tags {
			if aws.ToString(tag.Key) == "Name" {
				subnetInfo.Name = aws.ToString(tag.Value)
			}
			subnetInfo.Tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}

		// If no name tag, use subnet ID as name
		if subnetInfo.Name == "" {
			subnetInfo.Name = subnetInfo.ID
		}

		subnets = append(subnets, subnetInfo)
	}

	return subnets, nil
}

// GetSubnetUtilization calculates subnet utilization based on available IPs
func (c *Client) GetSubnetUtilization(ctx context.Context, subnetID string) (float64, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	}

	result, err := c.ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("failed to describe subnet %s: %w", subnetID, err)
	}

	if len(result.Subnets) == 0 {
		return 0, fmt.Errorf("subnet %s not found", subnetID)
	}

	subnet := result.Subnets[0]
	availableIPs := aws.ToInt32(subnet.AvailableIpAddressCount)

	// Calculate total IPs from CIDR
	cidr := aws.ToString(subnet.CidrBlock)
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CIDR %s: %w", cidr, err)
	}

	// Calculate total IPs (subtract 5 for AWS reserved IPs)
	prefixLen, _ := ipNet.Mask.Size()
	totalIPs := (1 << (32 - prefixLen)) - 5

	if totalIPs <= 0 {
		return 0, nil
	}

	usedIPs := totalIPs - int(availableIPs)
	utilization := (float64(usedIPs) / float64(totalIPs)) * 100

	return utilization, nil
}

// ValidateCredentials tests the AWS credentials and permissions
func (c *Client) ValidateCredentials(ctx context.Context) error {
	// Try to describe VPCs to validate credentials
	input := &ec2.DescribeVpcsInput{
		MaxResults: aws.Int32(5), // AWS requires minimum 5
	}

	_, err := c.ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to validate AWS credentials: %w", err)
	}

	return nil
}

// GetRegion returns the configured region
func (c *Client) GetRegion() string {
	return c.config.Region
}
