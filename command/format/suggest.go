
package format

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

type SuggestAWSImport struct {
	Ec2svc *ec2.EC2
}

func SuggestImport(r *InstanceDiff) string {

	filters := make(map[string]string)

	for _, attr := range r.Attributes {

		v := attr.NewValue

		switch {
		case v == "" && attr.NewComputed:

		case attr.Sensitive:

		default:
			if attr.Path == "tags.%" {

			} else if strings.HasPrefix(attr.Path, "tags.") {
				filters[strings.Replace(attr.Path,"tags.", "tags:",1)] = v
			} else if strings.HasSuffix(attr.Path, ":#"){

			} else {
				filters[attr.Path] = v
			}
		}
	}

	switch r.Addr.Type {
	case "aws_security_group":
		return sg_imports(r, filters)
	case "aws_elb":
		return elb_import(r, filters)
	case "aws_autoscaling_group":
		return asg_import(r, filters)
	case "aws_launch_configuration":
		return lc_import(r, filters)
	case "aws_db_instance":
		return rds_import(r, filters)
	case "aws_elasticache_replication_group":
		return redis_import(r, filters)
	case "aws_iam_role":
		return iam_role_import(r, filters)
	case "aws_iam_instance_profile":
		return instance_profile_import(r, filters);
	case "aws_s3_bucket":
		return s3_import(r, filters);
	case "aws_kms_alias":
		return kms_alias_import(r, filters)
		
		
	}
	return ""
}

/*
func sg_suggest_filter(r *InstanceDiff, cond map[string]string) []*ec2.Filter{

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
				aws.String(cond["name_prefix"]+"*"),
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
	resp, err := ec2svc.DescribeSecurityGroups(describeSecurityGroupInput);
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	if len(resp.SecurityGroups) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(*(resp.SecurityGroups[0].GroupId) + "\n\n");
		return buffer.String()
	}

	buffer.WriteString("Multiple Security Groups found\n")
	for _, res := range resp.SecurityGroups {
		buffer.WriteString( "> "+*res.GroupId + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}*/