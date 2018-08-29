package format

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func iam_role_policy_import(r *InstanceDiff, cond map[string]string) string {
	name, ok := cond["name"]

	if !ok {
		return ""
	}

	svc := iam.New(session.New())
	input := &iam.GetRolePolicyInput{
		PolicyName: aws.String(cond["name"]),
		RoleName:   aws.String(cond["role"]),
	}
	rolePolicy, err := svc.GetRolePolicy(input)

	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no iam role policy named " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}

	buffer.WriteString("terraform import  ")
	buffer.WriteString(r.Addr.String() + "  ")
	buffer.WriteString(*rolePolicy.RoleName + ":" + *rolePolicy.PolicyName + "\n\n")
	return buffer.String()
}
