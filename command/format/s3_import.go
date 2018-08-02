package format

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3_import(r *InstanceDiff, cond map[string]string) string {

	bucket, ok := cond["bucket"]
	if !ok {
		return ""
	}

	s3svc := s3.New(session.New())
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	var buffer bytes.Buffer

	if err != nil {
		buffer.WriteString("No Import: There is no s3 bucket named " + bucket)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	list := make([]string, 0)

	for _, res := range result.Buckets {
		if *(res.Name) == bucket {
			list = append(list, *res.Name)
		}
	}

	if len(list) == 0 {
		buffer.WriteString("No Import: There is no S3 with bucket name " + bucket)
		buffer.WriteString("\n")
		return buffer.String()
	}

	if err != nil {
		buffer.WriteString("No Import: There is no s3 bucket named " + bucket)
		buffer.WriteString("\n")
		return buffer.String()
		//log.Fatal(err.Error())
	}

	if len(list) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(list[0] + "\n\n")
		return buffer.String()
	}

	buffer.WriteString("Multiple s3 found\n")
	for _, res := range list {
		buffer.WriteString("# terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(res + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()

}
