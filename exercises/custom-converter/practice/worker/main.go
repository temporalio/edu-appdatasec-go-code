package main

import (
	"log"

	temporalconverters "edu-converters-go-code/exercises/custom-converter/practice"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	// TODO Part B: Import "go.temporal.io/sdk/temporal"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		// Set DataConverter here so that workflow and activity inputs/results will
		// be compressed as required.
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

	w := worker.New(c, "converters", worker.Options{})

	w.RegisterWorkflow(temporalconverters.Workflow)
	w.RegisterActivity(temporalconverters.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
