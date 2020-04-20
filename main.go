package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

func main() {
	ticker := time.NewTicker(1)
	for range ticker.C {
		run()
	}
}

func run() {
	region := "eu-west-2"
	config := aws.NewConfig().WithRegion(region)
	session, _ := session.NewSession(config)
	svc := athena.New(session)
	input := &athena.ListQueryExecutionsInput{}
	// Not iterating on NextToken, because it saves all queries.
	result, _ := svc.ListQueryExecutions(input)

	var runningCounter int
	var queueCounter int
	for _, q := range result.QueryExecutionIds {
		in := &athena.GetQueryExecutionInput{
			QueryExecutionId: q,
		}
		result, _ := svc.GetQueryExecution(in)
		// fmt.Println(result.QueryExecution.Status)

		if *result.QueryExecution.Status.State == "RUNNING" {
			// fmt.Println(result.QueryExecution)
			runningCounter = runningCounter + 1
		}
		if *result.QueryExecution.Status.State == "QUEUED" {
			// fmt.Println(result.QueryExecution)
			queueCounter = queueCounter + 1
		}
	}
	log.Printf("Current qeuries running: %v", runningCounter)
	log.Printf("Current queries queued:  %v", queueCounter)
}
