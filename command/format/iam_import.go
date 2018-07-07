package format

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func iam_role_import(r *InstanceDiff, cond map[string]string) string {

	name, ok := cond["name"]
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	svc := iam.New(session.New())
	input := &iam.GetRoleInput{
		RoleName: aws.String(cond["name"]),
	}
	result, err := svc.GetRole(input)

	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no iam named " + name)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	buffer.WriteString("terraform import  ")
	buffer.WriteString(r.Addr.String() + "  ")
	buffer.WriteString(*(result.Role.RoleName) + "\n\n")
	return buffer.String()

}
