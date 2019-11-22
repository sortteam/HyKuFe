package horovodjob

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"hykufe-operator/pkg/controller/ssh"
	"time"
)

type AWSController struct {
	session *session.Session
	instances []*ec2.Instance
}

func NewAWSController() (*AWSController) {
	obj := &AWSController{}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2")},
	)
	if err != nil {

	}
	obj.session = sess
	obj.instances = make([]*ec2.Instance, 10)

	return obj
}

func (ac *AWSController) CreateEC2Instance(instanceType string, replicas int64) ([]*ec2.Instance, error) {
	// Create EC2 service client
	svc := ec2.New(ac.session)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-00379ec40a3e30f87"),
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(replicas),
		MaxCount:     aws.Int64(replicas),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Mananged-HyKuFe-Operator"),
						Value: aws.String("True"),
					},
				},
			},
		},
		KeyName:	aws.String("SoRT"),
		SecurityGroupIds: []*string{
			aws.String("sg-086fe62b78c5edfa5"),
		},

		// TODO: Security Group, Subnet 추
	})
	if err != nil {
		return nil, err
	}


	targetInstances := []*string{}
	for _, instance := range runResult.Instances {
		ac.instances = append(ac.instances, instance)
		targetInstances = append(targetInstances, instance.InstanceId)
	}


	runningCount := int64(0)

	for runningCount != replicas {
		runningCount = 0
		output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds:         targetInstances,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range output.Reservations {
			for _, status := range r.Instances {
				if *status.State.Name == ec2.InstanceStateNameRunning {
					runningCount++

				}
			}
		}
		println(fmt.Sprintf("running : %d, replicas : %d", runningCount, replicas))

		time.Sleep(time.Second * 3)
	}

	sshClient := &ssh.SshClient{}
	if err := sshClient.NewSshClient("221.148.248.140", 7777); err != nil {
		return nil, fmt.Errorf("Failed to SSH Handshake : %v", err)
	}

	ac.getInstanceStatus(targetInstances)

	ipAggreation := ""
	for _, i := range ac.instances {
		ipAggreation += *i.PublicIpAddress
		ipAggreation += " "
	}

	println("mhg ip is " + ipAggreation)
	if err := sshClient.CommandExecution(fmt.Sprintf("~/onprem-kubespray/add_node.sh %s", ipAggreation)); err != nil {
		return nil, err
	}


	return runResult.Instances, nil
}

func (ac *AWSController) DeleteEC2Instance(instanceID string, ) error {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(false),
	}

	// Create EC2 service client
	svc := ec2.New(ac.session)

	_, err := svc.StopInstances(input)
	
	if err != nil {
		return err
	}

	tmpSlice := make([]*ec2.Instance, 10)
	for _, i := range ac.instances {
		if *i.InstanceId != instanceID {
			tmpSlice = append(tmpSlice, i)
		}
	}
	ac.instances = tmpSlice

	return nil
}

func (ac *AWSController) getInstanceStatus(targetInstances []*string) error {
	// Create EC2 service client
	svc := ec2.New(ac.session)

	output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds:         targetInstances,
	})

	if err != nil {
		return err
	}
	for _, i := range output.Reservations {
		ac.instances = i.Instances
	}

	return nil
}