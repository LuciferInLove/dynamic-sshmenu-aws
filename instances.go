package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSliceOfInstances(tags string, displayName string) ([]instance, error) {
	var (
		i                int = 1
		sliceOfInstances []instance
		instanceName     string
		filterMapsList   []*ec2.Filter
	)

	// Generating a list of instances filters
	if tags != "" {
		tagsList := strings.Split(tags, ";")
		for _, tag := range tagsList {
			keyValue := strings.Split(tag, ":")
			if len(keyValue) != 2 {
				return sliceOfInstances, fmt.Errorf("WrongTagDefinition")
			}

			values := strings.Split(keyValue[1], ",")
			var valuesList []*string
			for _, value := range values {
				valuesList = append(valuesList, aws.String(value))
			}

			filterMap := ec2.Filter{
				Name:   aws.String("tag:" + keyValue[0]),
				Values: valuesList,
			}
			filterMapsList = append(filterMapsList, &filterMap)
		}
	}

	filterMapsList = append(filterMapsList, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("running")},
	})

	// Getting instances
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ec2Service := ec2.New(awsSession)
	filterParams := &ec2.DescribeInstancesInput{
		Filters: filterMapsList,
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
