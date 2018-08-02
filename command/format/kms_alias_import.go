package format

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func kms_alias_import(r *InstanceDiff, cond map[string]string) string {

	name, ok := cond["name"]
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	svc := kms.New(session.New())
	input := &kms.ListAliasesInput{}

	list := make([]string, 0)

	err := svc.ListAliasesPages(input, func(result *kms.ListAliasesOutput, b bool) bool {
		for _, res := range result.Aliases {
			if *(res.AliasName) == name {
				list = append(list, *res.AliasName)
			}
		}
		return !b
	})

	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no kms alias named " + name)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	if len(list) == 0 {
		buffer.WriteString("No Import: There is no kms alias with name prefix " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}
	if len(list) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(list[0] + "\n\n")
		return buffer.String()
	}

	buffer.WriteString("Multiple kms alias found\n")
	for _, res := range list {
		buffer.WriteString("# terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(res + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()

}
