package format

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

func redis_import(r *InstanceDiff, cond map[string]string) string {

	name, ok := cond["replication_group_id"]
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	svc := elasticache.New(session.New())
	input := &elasticache.DescribeCacheClustersInput{
		CacheClusterId: aws.String("my-mem-cluster"),
	}

	result, err := svc.DescribeCacheClusters(input)
	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no redis with replication_group_id " + name)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	if len(result.CacheClusters) == 0 {
		buffer.WriteString("No Import: There is no redis with name replication_group_id " + name)
		buffer.WriteString("\n")
		return buffer.String()
	}
	if len(result.CacheClusters) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(*(result.CacheClusters[0].CacheClusterId) + "\n\n")
		return buffer.String()
	}

	buffer.WriteString("Multiple RDS found\n")
	for _, res := range result.CacheClusters {
		buffer.WriteString("# terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString( *(res.CacheClusterId) + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}
