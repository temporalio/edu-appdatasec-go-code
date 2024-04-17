package main

import (
	"log"

	compositeconverter "edu-converters-go-code/samples/composite-converter"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "converters", worker.Options{})

	w.RegisterWorkflow(compositeconverter.Workflow)
	w.RegisterActivity(compositeconverter.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
