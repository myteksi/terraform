package format

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"bytes"
	"github.com/aws/aws-sdk-go/service/elb"
)

func elb_import(r *InstanceDiff, cond map[string]string) string {
	name, ok := cond["name"];
	if !ok {
		return ""
	}
	names := []*string{}
	names = append(names, &name)
	elbsvc := elb.New(session.New())
	describeELB := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: names,
	}
	var buffer bytes.Buffer
	resp, err := elbsvc.DescribeLoadBalancers(describeELB)
	if err != nil {
		buffer.WriteString("No Import: There is no elb named " + name)
		buffer.WriteString("\n")
		return buffer.String();
		//log.Fatal(err.Error())
	}

	if len(resp.LoadBalancerDescriptions) == 1 {
		buffer.WriteString("terraform import  ")
		buffer.WriteString(r.Addr.String() + "  ")
		buffer.WriteString(*(resp.LoadBalancerDescriptions[0].LoadBalancerName) + "\n\n");
		return buffer.String()
	}

	buffer.WriteString("Multiple ELBs found\n")
	for _, res := range resp.LoadBalancerDescriptions {
		buffer.WriteString( "> "+*res.LoadBalancerName + "\n")
	}
	buffer.WriteString("\n")
	return buffer.String()
}