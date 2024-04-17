package main

import (
	"context"
	"log"

	temporalconverters "edu-converters-go-code/exercises/codec-server/practice"
	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		// Set DataConverter here to ensure that workflow inputs and results are
		// encoded as required.
		DataConverter: temporalconverters.DataConverter,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "converters_workflowID",
		TaskQueue: "converters",
	}

	we, err := c.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		temporalconverters.Workflow,
		"Plain text input",
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
