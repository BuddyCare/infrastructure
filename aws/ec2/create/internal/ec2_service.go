package internal

import (
	"context"
	"fmt"
	"github.com/BuddyCare/infrastructure/ec2/create/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"time"
)

type Ec2Svc struct {
	awsCli    *ec2.Client
	awsConfig aws.Config
	config    config.EC2Config
	ec2Info   info
}

type info struct {
	InstanceID string
	PublicIP   string
	PrivateIP  string
}

var InstanceTypes = map[string]types.InstanceType{
	"t2.medium": types.InstanceTypeT2Medium,
	"t2.small":  types.InstanceTypeT2Small,
}

func NewEc2Svc(ctx context.Context, config config.EC2Config) (*Ec2Svc, error) {
	ec2Svc := &Ec2Svc{
		config: config,
	}
	awsConf, err := awsConfig.LoadDefaultConfig(ctx, func(lc *awsConfig.LoadOptions) error {
		lc.Region = config.Region
		lc.Credentials = ec2Svc
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error loading aws config: %s", err.Error())
	}
	ec2Svc.awsConfig = awsConf
	ec2Svc.awsCli = ec2.NewFromConfig(ec2Svc.awsConfig)
	return ec2Svc, nil
}

func (ec2Svc *Ec2Svc) CreateEc2Instance(ctx context.Context) error {
	typeValue, existsType := InstanceTypes[ec2Svc.config.InstanceType]
	if len(typeValue) == 0 || !existsType {
		return fmt.Errorf("instance type %s is invalid", ec2Svc.config.InstanceType)
	}

	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(ec2Svc.config.ImageId),
		InstanceType: typeValue,
		MinCount:     aws.Int32(ec2Svc.config.MinCount),
		MaxCount:     aws.Int32(ec2Svc.config.MaxCount),
		KeyName:      aws.String(ec2Svc.config.KeyName),
		NetworkInterfaces: []types.InstanceNetworkInterfaceSpecification{
			{
				AssociatePublicIpAddress: aws.Bool(true),
				DeviceIndex:              aws.Int32(0),
			},
		},
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(ec2Svc.config.Name),
					},
				},
			},
		},
	}
	result, err := ec2Svc.awsCli.RunInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("unable to create instance, %v", err)
	}
	instanceId := *result.Instances[0].InstanceId
	describeInstancesOutput, err := ec2Svc.getDescribeInstance(ctx, instanceId)
	if err != nil {
		return err
	}
	instance := describeInstancesOutput.Reservations[0].Instances[0]
	ec2Svc.ec2Info = info{
		InstanceID: *instance.InstanceId,
		PrivateIP:  *instance.PrivateIpAddress,
		PublicIP:   *instance.PublicIpAddress,
	}
	fmt.Printf("Created instance:\n  InstanceID: %s\n  PrivateIP: %s\n  PublicIP: %s\n",
		ec2Svc.ec2Info.InstanceID, ec2Svc.ec2Info.PrivateIP, ec2Svc.ec2Info.PublicIP)
	return nil
}

func (ec2Svc *Ec2Svc) Retrieve(ctx context.Context) (aws.Credentials, error) {
	cred := aws.Credentials{}
	cred.AccessKeyID = ec2Svc.config.Credentials.AccessKeyId
	cred.SecretAccessKey = ec2Svc.config.Credentials.SecretKeyId
	return cred, nil
}

func (ec2Svc *Ec2Svc) getDescribeInstance(ctx context.Context, instanceId string) (ec2.DescribeInstancesOutput, error) {
	describeInstancesOutput := &ec2.DescribeInstancesOutput{}
	var err error
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	}

	for i := 0; i <= 5; i++ {
		describeInstancesOutput, err = ec2Svc.awsCli.DescribeInstances(ctx, describeInstancesInput)
		if err != nil {
			return ec2.DescribeInstancesOutput{}, fmt.Errorf("failed to describe instances: %v", err)
		}
		if describeInstancesOutput.Reservations[0].Instances[0].PublicIpAddress != nil {
			break
		}
		fmt.Printf("waiting for PublicIpAddress assigned....")
		time.Sleep(3 * time.Second)
	}
	return *describeInstancesOutput, nil
}
