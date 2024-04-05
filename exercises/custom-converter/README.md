# Exercise 1: Implementing a Custom Data Converter

During this exercise, you will: 

* Output typical payloads from a Temporal Workflow using the default Data Converter
* Implement a Custom Data Converter that encrypts Workflow output
* Implement a Failure Converter and demonstrate parsing its output

Make your changes to the code in the `practice` subdirectory (look for 
`TODO` comments that will guide you to where you should make changes to 
the code). If you need a hint or want to verify your changes, look at 
the complete version in the `solution` subdirectory.


## Part A: Implement a Custom Data Converter

1. Defining a Custom Data Converter is a straightforward change to your existing
   Worker and Starter code. The example in the `practice` subdirectory of this
   exercise is missing the necessary change to use a Custom Data Converter.
   meaning that you can run it out of the box, and produce JSON output using the
   Default Data Converter. You'll do this first, so you have an idea of the
   expected output. First, start the Worker:

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
     ID          Time                     Type
      1  2024-01-16T17:10:53Z  WorkflowExecutionStarted
      2  2024-01-16T17:10:53Z  WorkflowTaskScheduled
      3  2024-01-16T17:10:53Z  WorkflowTaskStarted
      4  2024-01-16T17:10:53Z  WorkflowTaskCompleted
      5  2024-01-16T17:10:53Z  MarkerRecorded
      6  2024-01-16T17:10:53Z  MarkerRecorded
      7  2024-01-16T17:10:53Z  ActivityTaskScheduled
      8  2024-01-16T17:10:53Z  ActivityTaskStarted
      9  2024-01-16T17:10:53Z  ActivityTaskCompleted
     10  2024-01-16T17:10:53Z  WorkflowTaskScheduled
     11  2024-01-16T17:10:53Z  WorkflowTaskStarted
     12  2024-01-16T17:10:53Z  WorkflowTaskCompleted
     13  2024-01-16T17:10:53Z  WorkflowExecutionCompleted

   Result:
     Status: COMPLETED
     Output: ["Received Plain text input"]
   ```

   You should now have an idea of how this Workflow runs ordinarily — it outputs
   the string `Received Plain text input`. In the next step, you'll add a Custom
   Data Converter.
4. To add a Custom Data Converter, you only need to add a `DataConverter`
   parameter to `client.Dial()` where it is used in both `starter/main.go` and
   `worker/main.go`. Make this change and the save both files. You don't need to
   change anything in your Workflow code.
5. Next, open `data_converter.go`. This contains the Custom Converter code
   you'll be using. The `Encode()` function should marshal a payload to JSON
   then compress it using Go's [snappy](https://github.com/google/snappy) codec,
   and set the file metadata. The `Decode()` function already does the same
   thing in reverse. Add the missing calls to the `Encode()` function (you can
   use the `Decode()` function as a hint). Then save the file.
6. Now you can re-run the Workflow with your Custom Converter. Stop your Worker
   (with `Ctrl+C` in a blocking terminal) and restart it with `go run
   worker/main.go`, then re-run the Workflow with `go run starter/main.go`.
   You should once again receive `Received Plain text input`.
7. Finally, get the result again with `temporal workflow show -w
   converters_workflowID`. This time, your output will be encoded:

   ```
   ...
   Result:
     Status: COMPLETED
     Output: [encoding binary/snappy: payload encoding is not supported]
   ```

  The `payload encoding is not supported` message is normal — the Temporal
  Cluster itself can't use the `Decode` function directly without a Codec
  Server, which you'll create in the next exercise. In the meantime, you have
  successfully implemented a Custom Data Converter, and in the next step, you'll
  add more features to it. 


## Part B: Implement a Failure Converter

1. The next feature you may add is a Failure Converter. Failure messages and
   stack traces are not encoded as codec-capable Payloads by default; you must
   explicitly enable encoding these common attributes on failures. If your
   errors might contain sensitive information, you can encrypt the message and
   stack trace by configuring the default Failure Converter to use your encoded
   attributes, in which case it moves your `message` and `stack_trace` fields to
   a Payload that's run through your codec. To do this, you can override the
   default Failure Converter with a single additional parameter,
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
   ...
   Status: FAILED
   Failure: &Failure{Message:Encoded failure,Source:GoSDK,StackTrace:,Cause:nil,FailureType:Failure_ApplicationFailureInfo,}
   ```


### This is the end of the exercise.

