# Exercise 1: Implementing a Custom Codec

During this exercise, you will: 

* Output typical payloads from a Temporal Workflow using the default Data Converter
* Implement a Custom Codec that compresses Workflow output
* Implement a Failure Converter and demonstrate parsing its output

Make your changes to the code in the `practice` subdirectory (look for 
`TODO` comments that will guide you to where you should make changes to 
the code). If you need a hint or want to verify your changes, look at 
the complete version in the `solution` subdirectory.


## Part A: Implement a Custom Codec

1. Defining a Custom Codec is a straightforward change to your existing Worker
   and Starter code. The example in the `practice` subdirectory of this exercise
   is missing the necessary change to use a Custom Codec, meaning that you can
   run it out of the box, and produce JSON output using the Default Data
   Converter. You'll do this first, so you have an idea of the expected output.
   First, start the Worker:

   ```shell
   go run worker/main.go
   ```

2. Next, run the Workflow starter:

   ```shell
   go run starter/main.go
   ```

   You should receive `Received Plain text input` in this terminal.

3. After that, you can use the `temporal` CLI to show the Workflow result:

   ```shell
   temporal workflow show -w converters_workflowID
   ```

   ```
   Progress:
     ID           Time                     Type
       1  2024-08-26T20:28:51Z  WorkflowExecutionStarted
       2  2024-08-26T20:28:51Z  WorkflowTaskScheduled
       3  2024-08-26T20:28:51Z  WorkflowTaskStarted
       4  2024-08-26T20:28:51Z  WorkflowTaskCompleted
       5  2024-08-26T20:28:51Z  ActivityTaskScheduled
       6  2024-08-26T20:28:51Z  ActivityTaskStarted
       7  2024-08-26T20:28:51Z  ActivityTaskCompleted
       8  2024-08-26T20:28:51Z  WorkflowTaskScheduled
       9  2024-08-26T20:28:51Z  WorkflowTaskStarted
      10  2024-08-26T20:28:51Z  WorkflowTaskCompleted
      11  2024-08-26T20:28:51Z  WorkflowExecutionCompleted

   Results:
     Status          COMPLETED
     Result          "Received Plain text input"
     ResultEncoding  json/plain
   ```

   You should now have an idea of how this Workflow runs ordinarily — it outputs
   the string `Received Plain text input`. In the next step, you'll add a Custom
   Codec.
4. To add a Custom Codec, you need to add a `DataConverter` parameter
   to `client.Dial()`. Both your Client and your Worker perform this call, and
   you technically only need to update one of them before your Workflow runs,
   but to be safe, edit both `starter/main.go` and `worker/main.go`. Make this
   change and the save both files. You don't need to change anything in your
   Workflow code.
5. Next, open `data_converter.go`. This contains the Custom Codec code
   you'll be using. The `Encode()` function should marshal a payload to JSON
   then compress it using Go's [snappy](https://github.com/google/snappy) codec,
   and set the file metadata. The `Decode()` function already does the same
   thing in reverse. Add the missing calls to the `Encode()` function (you can
   use the `Decode()` function as a hint). Then save the file.
6. Now you can re-run the Workflow with your Custom Codec. Stop your Worker
   (with `Ctrl+C` in a blocking terminal) and restart it with `go run
   worker/main.go`, then re-run the Workflow with `go run starter/main.go`.
   You should once again receive `Received Plain text input`.
7. Finally, get the result again with `temporal workflow show -w
   converters_workflowID`. This time, your output will be encoded:

   ```
   Results:
     Status          COMPLETED
     Result          {"metadata":{"encoding":"YmluYXJ5L3NuYXBweQ=="},"data":"NdAKFgoIZW5jb2RpbmcSCmpzb24vcGxhaW4SGyJSZWNlaXZlZCBQbGFpbiB0ZXh0IGlucHV0Ig=="}
     ResultEncoding  binary/snappy
   ```

  The Temporal Cluster itself can't use the `Decode` function directly without a
  Codec Server, which you'll create in the next exercise. In the meantime, you
  have successfully customized your Data Converter with a codec, and in the next
  step, you'll add more features to it. 


## Part B: Implement a Failure Converter

1. The next feature you may add is a Failure Converter. To do this, you can
   override the default Converter by providing the parameter `DataConverter` with your custom Codec. To make sure the Error Message gets encoded, you also need provide an additional parameter,
   `EncodeCommonAttributes: true`. Make this change to `client.Dial()` where it
   is used in both `starter/main.go` and `worker/main.go`, as you did before.
   You will also need to import `go.temporal.io/sdk/temporal` into these files.
2. To test your Failure Converter, change your Workflow to return an artificial
   error. Change the `ExecuteActivity` call to throw an error where there isn't
   one, like so:

   ```go
	err = workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
	if err == nil {
		err = errors.New("This is an artificial error") // add this statement
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
   ```

   Don't forget to add the `errors` package to `workflow.go` as well.
3. After making this change, stop your Worker (with `Ctrl+C` in a blocking
   terminal) and restart it with `go run worker/main.go`, then re-run the Workflow
   with `go run starter/main.go`. The Workflow should fail.
4. Run `temporal workflow show -w converters_workflowID` to get the status of your
   failed Workflow. Notice that the `Failure:` field should now display an encoded
   result, rather than a plain text error:

   ```
   Results:
     Status   FAILED
     Failure
       Message: Encoded failure
   ```


### This is the end of the exercise.

