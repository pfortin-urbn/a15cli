package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"flag"
)

var action = flag.String("action", "", "Action to perform (create,depth, purge, send (test), receive")
var queue = flag.String("queue", "", "Queue to work with")
var local = flag.Bool("local", true, "Local (ElasticMQ) or Remote (AmazonSQS")

var queueName string

func init() {
	fmt.Println("args len:", len(os.Args))
	if(len(os.Args) != 3 && len(os.Args) != 5 && len(os.Args) != 7) {
		fmt.Println("Usage: ./sqsQueueUtil -action <action> -queue <queueName> [-local <true|false>]")
		os.Exit(0)
	}
	flag.Parse()

	fmt.Println("local:", *local)

	//if *local == true {
	//	svc = sqs.New(session.New(), &aws.Config{Endpoint: aws.String("http://localhost:9324"), Region: aws.String("us-east-1")})
	//	url = "http://localhost:9324/queue/" + *queue
	//	queueName = *queue
	//} else {
	svc = sqs.New(session.New(), &aws.Config{ Region: aws.String("us-east-1") })
	url = "https://sqs.us-east-1.amazonaws.com/478989820108/" + *queue
	queueName = *queue
	//}
	fmt.Println("svc:", svc, "url", url, "queueName", queueName)
}

var svc 			*sqs.SQS
var url 			string
var attrib 			string

/*
 *
 *	Program entry point creates connection to SQS and Mongo then pool SQS for messages
 *
 */
func main() {
	switch *action {
	case "create":
		createSQSQueue()
		break
	case "depth":
		getSQSQueueDepth()
		break
	case "purge":
		purgeQueue()
		break;
	case "send":
		sendMessage()
		break;
	case "receive":
		receiveMessage()
		break;
	default:
		fmt.Println("Unrecognized action - try again!")
	}
}

func createSQSQueue() {
	params := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName), // Required
	}
	resp, err := svc.CreateQueue(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}


func getSQSQueueDepth() {

	attrib = "ApproximateNumberOfMessages"
	sendParams := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(url), // Required
		AttributeNames: []*string{
			&attrib, // Required
		},
	}
	resp2, sendErr := svc.GetQueueAttributes(sendParams)
	if sendErr != nil {
		fmt.Println("Depth: " + sendErr.Error())
		return
	}
	fmt.Println(resp2)
}

func sendMessage() {
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String("Testing 1,2,3,..."), // Required
		QueueUrl:     aws.String(url), // Required
	}
	resp, err := svc.SendMessage(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)

}


func receiveMessage() {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(url), // Required
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}

func purgeQueue() {
	params := &sqs.PurgeQueueInput{
		QueueUrl: aws.String(url), // Required
	}
	resp, err := svc.PurgeQueue(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}
