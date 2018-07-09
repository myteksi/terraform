package format

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func asg_import(r *InstanceDiff, cond map[string]string) string {
	name, ok := cond["name"]
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

	result, err := elbsvc.DescribeAutoScalingGroups(input)
	if err != nil {
		buffer.WriteString("No Import: There is no asg named " + name)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	if len(result.AutoScalingGroups) == 0 {
		buffer.WriteString("No Import: There is no asg named " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}

	if len(result.AutoScalingGroups) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(*(result.AutoScalingGroups[0].AutoScalingGroupName) + "\n\n")
		return buffer.String()
	}

	for _, res := range result.AutoScalingGroups {
		buffer.WriteString("> " + *res.AutoScalingGroupName + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}
