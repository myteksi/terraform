package format

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"bytes"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"strings"
)

func lc_import(r *InstanceDiff, cond map[string]string) string {

	name, ok := cond["name_prefix"];
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	svc := autoscaling.New(session.New())
	input := &autoscaling.DescribeLaunchConfigurationsInput{
		LaunchConfigurationNames: []*string{
		},
	}

	list := make([]string, 0)
	err := svc.DescribeLaunchConfigurationsPages(input, func (result *autoscaling.DescribeLaunchConfigurationsOutput, lastPage bool) bool {
		for _, res := range result.LaunchConfigurations {

			if strings.HasPrefix(*(res.LaunchConfigurationName), name) {
				list = append(list, *res.LaunchConfigurationName)
			}

		}
		return !lastPage
	})
	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no lc named " + name)
		buffer.WriteString("\n")
		return buffer.String();
		//log.Fatal(err.Error())
	}

	if len(list) == 0 {
		buffer.WriteString("No Import: There is no LC with name prefix " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}
	if len(list) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(list[0]+ "\n\n");
		return buffer.String()
	}

	buffer.WriteString("Multiple LC found\n")
	for _, res := range list {
		buffer.WriteString( "> "+ res + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}