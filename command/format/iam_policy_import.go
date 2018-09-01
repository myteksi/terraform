package format

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func iam_policy_import(r *InstanceDiff, cond map[string]string) string {
	name, ok := cond["name"]

	if !ok {
		return ""
	}

	svc := iam.New(session.New())
	var buffer bytes.Buffer

	params := iam.ListPoliciesInput{}

	params.SetOnlyAttached(true)

	var policyArn string

	svc.ListPoliciesPages(&params,
		func(page *iam.ListPoliciesOutput, lastPage bool) bool {
			for _, policy := range page.Policies {
				if *policy.PolicyName == name {
					policyArn = *policy.Arn
				}
			}
			if lastPage {
				return false
			}
			return true
		})

	if policyArn == "" {
		buffer.WriteString("No Import: There is no iam policy named " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}

	buffer.WriteString("terraform import  ")
	buffer.WriteString(r.Addr.String() + "  ")
	buffer.WriteString(policyArn + "\n\n")
	return buffer.String()
}
