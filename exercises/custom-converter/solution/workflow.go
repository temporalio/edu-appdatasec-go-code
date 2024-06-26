package temporalconverters

import (
	"context"
	"errors"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a standard workflow definition.
// Note that the Workflow and Activity don't need to care that
// their inputs/results are being encoded.
func Workflow(ctx workflow.Context, input string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Converter workflow started", "input", input)

	var result string

	err := workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
	if err == nil {
		err = errors.New("This is an artificial error")
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("Converter workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, input string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "input", input)

	return "Received " + input, nil
}
