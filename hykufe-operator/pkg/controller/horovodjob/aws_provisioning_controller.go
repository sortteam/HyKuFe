package horovodjob
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSController struct {
	session *session.Session
	instances []*ec2.Instance
}

func NewAWSController() (*AWSController, error) {
	obj := &AWSController{}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2")},
	)
	if err != nil {
		return nil, err
	}
	obj.session = sess
	obj.instances = make([]*ec2.Instance, 10)

	return obj, nil
}

func (ac *AWSController) CreateEC2Instance(instanceType string, ) (*ec2.Instance, error) {
	// Create EC2 service client
	svc := ec2.New(ac.session)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-00379ec40a3e30f87"),
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: nil,
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Mananged-HyKuFe-Operator"),
						Value: aws.String("True"),
					},
				},
			},
		},
		// TODO: Security Group, Subnet 추
	})
	if err != nil {
		return nil, err
	}

	instanceRef := runResult.Instances[0]
	ac.instances = append(ac.instances, instanceRef)

	return instanceRef, nil
}
func (ac *AWSController) DeleteEC2Instance(instanceID string, ) error {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		DryRun: aws.Bool(true),
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