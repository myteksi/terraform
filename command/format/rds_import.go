package format

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func rds_import(r *InstanceDiff, cond map[string]string) string {
	fmt.Println(cond)
	name, ok := cond["identifier"]
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	svc := rds.New(session.New())
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(cond["identifier"]),
	}

	result, err := svc.DescribeDBInstances(input)
	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no rds with identifier " + name)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	if len(result.DBInstances) == 0 {
		buffer.WriteString("No Import: There is no RDS with name prefix " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}
	if len(result.DBInstances) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(*(result.DBInstances[0].DBInstanceIdentifier) + "\n\n")
		return buffer.String()
	}

	buffer.WriteString("Multiple RDS found\n")
	for _, res := range result.DBInstances {
		buffer.WriteString("> " + *(res.DBClusterIdentifier) + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}
