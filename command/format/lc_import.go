package format

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func lc_import(r *InstanceDiff, cond map[string]string) string {
	name, ok := cond["name_prefix"]
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	elbsvc := autoscaling.New(session.New())

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(name),
		},
	}

	var buffer bytes.Buffer

	result, _ := elbsvc.DescribeAutoScalingGroups(input)

	LaunchConfigName := *(result.AutoScalingGroups[0].LaunchConfigurationName)

	buffer.WriteString("terraform import  ")
	buffer.WriteString(r.Addr.String() + "  ")
	buffer.WriteString(LaunchConfigName + "\n\n\n")

	return buffer.String()
}
