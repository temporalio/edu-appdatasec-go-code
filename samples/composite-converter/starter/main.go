package main

import (
	"context"
	"log"

	compositeconverter "edu-converters-go-code/samples/composite-converter"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func main() {
	dataConverter := converter.NewCompositeDataConverter(
		converter.NewNilPayloadConverter(),
		converter.NewByteSlicePayloadConverter(),
		converter.NewProtoJSONPayloadConverter(),
		converter.NewProtoPayloadConverter(),
		compositeconverter.NewCustomPayloadConverter(),
		converter.NewJSONPayloadConverter(),
	)

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		// Set DataConverter here.
		DataConverter: dataConverter,
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
		compositeconverter.Workflow,
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
