package main

import (
	"context"
	"log"

	temporalconverters "edu-converters-go-code/exercises/custom-converter/practice"

	"go.temporal.io/sdk/client"
	// TODO Part B: Import "go.temporal.io/sdk/temporal"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		// Set DataConverter here to ensure that workflow inputs and results are
		// encoded as required.
		// TODO Part A: Set a `DataConverter` key to use the `DataConverter` from `data_converter.go`.
		// This overrides the stock behavior — otherwise, the default data converter will be used.
		// TODO Part B: Set a `FailureConverter` key to use an instance of
		// `temporal.NewDefaultFailureConverter` with a single argument,
		// `temporal.DefaultFailureConverterOptions{}`, and in the options array, set
		// `EncodeCommonAttributes: true`.
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
