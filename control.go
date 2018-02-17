package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

/** K8LSecurityGroupName string name of the k8l-security-group */
const K8LSecurityGroupName = "k8l-sg"
const K8LInstanceTagName = "K8L_GUNPOWDER"
const K8LInstanceTagValue = "yes"

func main() {
	amiPerRegion := map[string]string{
		"eu-central-1": "ami-5055cd3f",
		"eu-west-1":    "ami-1b791862",
		"eu-west-3":    "ami-c1cf79bc",
	}

	var command, region, instanceID string
	var enableVerbose bool

	flag.StringVar(&command, "c", "", "Command (list, create, destroy, ensure-sg).")
	flag.StringVar(&region, "r", os.Getenv("AWS_REGION"), "AWS-Region")
	flag.StringVar(&instanceID, "i", "", "Instance-ID")
	flag.BoolVar(&enableVerbose, "v", false, "Verbose")
	flag.Parse()

	if enableVerbose {
		os.Setenv("__VERBOSE", "1")
	}

	if region == "" {
		fmt.Println("No region specified.")
		os.Exit(1)
	}

	verbose("Configured AWS-Region:" + region)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	switch command {

	case "list":
		instances := getAllInstances(sess)
		fmt.Println(instances.GoString())

	case "create":
		securityGroupID := ensureSecurityGroup(sess, K8LSecurityGroupName)
		publicIP := createInstance(sess, amiPerRegion[region], securityGroupID, K8LSecurityGroupName)
		fmt.Println(publicIP)

	case "destroy":
		var instanceIDs []string

		if instanceID != "" {
			instanceIDs = []string{instanceID}
		} else {
			verbose("Destroying all instances with tag:" + K8LInstanceTagName)
			instanceIDs = getAllInstanceIDsForTag(sess, K8LInstanceTagName)
		}

		for i := range instanceIDs {
			instanceID = instanceIDs[i]
			destroyInstance(sess, instanceID)
		}

	default:
		fmt.Println("No command specified.")
		os.Exit(1)
	}
}

func ensureSecurityGroup(sess *session.Session, groupName string) string {

	groupFound := true

	groups, err := getSecurityGroups(sess, []string{groupName})

	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "InvalidGroup.NotFound":
			groupFound = false
			verbose("Security-Group not found...creating one")
		default:
			handleError(err)
		}
	}

	if !groupFound {
		group := createSecurityGroup(sess, groupName)
		fmt.Println(group)
		verbose("Security-Group created")

		time.Sleep(time.Duration(20) * time.Second)

		attachSecurityGroupRules(sess, groupName)
		verbose("Security-Group-Rules attached")

		return aws.StringValue(group.GroupId)
	}

	return aws.StringValue(groups.SecurityGroups[0].GroupId)
}

func getSecurityGroups(sess *session.Session, groupNames []string) (*ec2.DescribeSecurityGroupsOutput, error) {
	svc := ec2.New(sess)

	result, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: aws.StringSlice(groupNames),
	})

	return result, err
}

func createSecurityGroup(sess *session.Session, name string) *ec2.CreateSecurityGroupOutput {
	input := &ec2.CreateSecurityGroupInput{
		Description: &name,
		GroupName:   &name,
	}
	ec2 := ec2.New(sess)

	out, err := ec2.CreateSecurityGroup(input)
	if err != nil {
		handleError(err)
	}

	return out
}

func attachSecurityGroupRules(sess *session.Session, groupName string) {
	svc := ec2.New(sess)

	_, err := svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(groupName),
		IpPermissions: []*ec2.IpPermission{
			(&ec2.IpPermission{}).
				SetIpProtocol("tcp").
				SetFromPort(0).
				SetToPort(65535).
				SetIpRanges([]*ec2.IpRange{
					{CidrIp: aws.String("0.0.0.0/0")},
				}),
		},
	})

	if err != nil {
		handleError(err)
	}
}

func getAllInstances(sess *session.Session) *ec2.DescribeInstancesOutput {
	input := &ec2.DescribeInstancesInput{}
	ec2 := ec2.New(sess)

	fmt.Println(reflect.TypeOf(ec2))

	out, err := ec2.DescribeInstances(input)
	if err != nil {
		handleError(err)
	}

	return out
}

func createInstance(sess *session.Session,
	ami string,
	securityGroupID string,
	securityGroupName string,
) string {

	svc := ec2.New(sess)

	verbose("Creating instance, using ami:" + ami + ", security-group:" + securityGroupID)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String(ami),
		InstanceType:     aws.String("t2.micro"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		SecurityGroupIds: aws.StringSlice([]string{securityGroupID}),
		SecurityGroups:   aws.StringSlice([]string{securityGroupName}),
		KeyName:          aws.String("home"),
	})

	if err != nil {
		handleError(err)
	}

	instanceID := *runResult.Instances[0].InstanceId
	verbose("Created instance:" + instanceID)

	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(K8LInstanceTagName),
				Value: aws.String(K8LInstanceTagValue),
			},
		},
	})
	if err != nil {
		handleError(err)
	}

	verbose("Successfully tagged instance")
	verbose("Waiting for instance to be ready...")

	err = svc.WaitUntilInstanceStatusOk(&ec2.DescribeInstanceStatusInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	})
	if err != nil {
		handleError(err)
	}

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	})

	if err != nil {
		handleError(err)
	}

	publicIP := aws.StringValue(result.Reservations[0].Instances[0].PublicIpAddress)
	verbose("Instance:" + instanceID + " with Public-IP:" + publicIP + "is ready")

	return publicIP
}

func getAllInstanceIDsForTag(sess *session.Session, tag string) []string {
	svc := ec2.New(sess)

	verbose("Getting all instances for tag:" + tag)

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:" + tag),
				Values: []*string{
					aws.String("yes"),
				},
			},
		},
	})

	if err != nil {
		handleError(err)
	}

	instanceIDs := make([]string, 0, 16)

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, aws.StringValue(instance.InstanceId))
		}
	}

	verbose("Found following instances:" + strings.Join(instanceIDs, ","))

	return instanceIDs
}

func destroyInstance(sess *session.Session, instanceID string) {
	svc := ec2.New(sess)

	verbose("Destroying instance:" + instanceID)

	_, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	})

	if err != nil {
		handleError(err)
	}

	verbose("Instance:" + instanceID + " destroyed")
}

func handleError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "InvalidGroupId.Malformed":
			fallthrough
		case "InvalidGroup.NotFound":
			exitErrorf("%s.", aerr.Message())
		}
	}

	panic(err)
	exitErrorf("Error: %v", err)
	os.Exit(1)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func verbose(message string) {
	if os.Getenv("__VERBOSE") == "1" {
		t := time.Now()
		fmt.Println(t.Format("2006-01-02 15:04:05"), " - ", message)
	}
}
