package format

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"bytes"
)

func instance_profile_import(r *InstanceDiff, cond map[string]string) string {

	name, ok := cond["name"];
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	svc := iam.New(session.New())
	input := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(cond["name"]),
	}
	result, err := svc.GetInstanceProfile(input)

	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no aws_iam_instance_profile named " + name)
		buffer.WriteString("\n")
		return buffer.String();
		//log.Fatal(err.Error())
	}

	buffer.WriteString("terraform import  ")
	buffer.WriteString(r.Addr.String() + "  ")
	buffer.WriteString(*(result.InstanceProfile.InstanceProfileId) + "\n\n");
	return buffer.String()

}