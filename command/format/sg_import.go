package format

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

func sg_suggest_filter(r *InstanceDiff, cond map[string]string) []*ec2.Filter {

	listFilter := []*ec2.Filter{}


	if _, ok := cond["name"]; ok {
		aFilter := &ec2.Filter{
			Name: aws.String("group-name"),
			Values: []*string{
				aws.String(cond["name"]),
			},
		}
		listFilter = append(listFilter, aFilter)
		return listFilter
	}

	if _, ok := cond["name_prefix"]; ok {
		aFilter := &ec2.Filter{
			Name: aws.String("group-name"),
			Values: []*string{
				aws.String(cond["name_prefix"] + "*"),
			},
		}
		listFilter = append(listFilter, aFilter)
	}

	for k, v := range cond {
		if !strings.HasPrefix(k, "tag:") {
			continue
		}
		aFilter := &ec2.Filter{
			Name: aws.String(k),
			Values: []*string{
				aws.String(v),
			},
		}
		listFilter = append(listFilter, aFilter)
	}

	return listFilter

}

func sg_imports(r *InstanceDiff, cond map[string]string) string {
	ec2svc := ec2.New(session.New())
	listFilter := sg_suggest_filter(r, cond)
	describeSecurityGroupInput := &ec2.DescribeSecurityGroupsInput{
		Filters: listFilter,
	}
	var buffer bytes.Buffer

	if len(listFilter) == 0 {
		buffer.WriteString("No Import: No security groups found \n\n")
		return buffer.String()
	}

	resp, err := ec2svc.DescribeSecurityGroups(describeSecurityGroupInput)
	if err != nil {
		buffer.WriteString("No Import: No security groups found \n\n")
		return buffer.String()
	}

	if len(resp.SecurityGroups) < 1 {
		buffer.WriteString("No Import: No security groups found with name, name_prefix  or tags\n\n")
		return buffer.String()
	}

	if len(resp.SecurityGroups) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(*(resp.SecurityGroups[0].GroupId) + "\n\n")
		return buffer.String()
	}

	buffer.WriteString("Multiple Security Groups found\n")
	for _, res := range resp.SecurityGroups {
		buffer.WriteString("> " + *res.GroupId + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}
