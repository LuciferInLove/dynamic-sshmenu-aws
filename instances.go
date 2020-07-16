package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSliceOfInstances(tagKey string, tagValue string, displayName string) ([]instance, error) {
	var (
		i                int = 1
		sliceOfInstances []instance
		instanceName     string
	)

	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ec2Service := ec2.New(awsSession)
	filterParams := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:" + tagKey),
				Values: []*string{aws.String(tagValue)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	describeInstancesResponce, err := ec2Service.DescribeInstances(filterParams)

	if err != nil {
		return sliceOfInstances, err
	}

	// Appending instances to sliceOfInstances
	for idx := range describeInstancesResponce.Reservations {
		for _, inst := range describeInstancesResponce.Reservations[idx].Instances {
			// Getting instance name from Name tag
			for _, tag := range inst.Tags {
				if *tag.Key == displayName {
					instanceName = *tag.Value
				}
			}
			currentInstance := instance{
				Number: i,
				IP:     *inst.PrivateIpAddress,
				Name:   instanceName,
				Zone:   *inst.Placement.AvailabilityZone,
			}
			sliceOfInstances = append(sliceOfInstances, currentInstance)
			i++
		}
	}

	return sliceOfInstances, nil
}
