# Exercise 1: Implement a Custom Data Converter

During this exercise, you will: 

* Output typical payloads from a Temporal Workflow using the default Data Converters
* Implement a Custom Data Converter that encrypts Workflow output
* Implement a Failure Converter and demonstrate parsing its output
* Implement a Composite Data Converter that combines multiple conversion steps

Make your changes to the code in the `practice` subdirectory (look for 
`TODO` comments that will guide you to where you should make changes to 
the code). If you need a hint or want to verify your changes, look at 
the complete version in the `solution` subdirectory.


## Part A: Implement a Custom Data Converter

1. Similar to other Temporal features like opting in to Worker Versioning,
   defining a Custom Data Converter is only a one-line change to your existing
   Worker and Starter code. The Example in the Practice subdirectory of this
   exercise is missing the necessary change to use a Custom Data Converter,
   meaning that you can run it out of the box, and produce JSON output using the
   Default Data Converter. You'll do this first, so you have an idea of the
   expected output. First, start the Worker:

   ```shell
   go run worker/main.go
   ```

   Next, run the Workflow starter:

   ```shell
   go run starter/main.go
   ```

   After that, you can use the `temporal` CLI to show the Workflow result:

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
   the string "Received Plain text input". In the next step, you'll add a Custom
   Data Converter.
2. To add a Custom Data Converter, you don't need to change anything in your
   Workflow code. You only need to add a `DataConverter` parameter to
   `client.Dial()` where it is used in both `starter.go` and `worker.go`.
3. Next, take a look in `data_converter.go`. This contains the Custom Converter
   code you'll be using. The `Encode()` function should marshal a payload to
   JSON then compress it using Go's [snappy](https://github.com/google/snappy)
   codec, and set the file metadata. The `Decode()` function does the same thing
   in reverse. Add the missing calls to the `Encode()` function (you can use the
   `Decode()` function as a hint).
4. Now you can re-run the Workflow with your Custom Conveter. Stop your Worker
   (with `Ctrl+C` in a blocking terminal) and restart it with `go run
   worker/main.go`, then re-run the workflow with `go run starter/main.go`.
   Finally, get the result again with `temporal workflow show -w
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
   attributes, in which case it moves your message and stack_trace fields to a
   Payload that's run through your codec. To do this, you can override the
   default Failure Converter with a single additional parameter,
   `EncodeCommonAttributes: true`. Make this change to `client.Dial()` where it
   is used in both `starter.go` and `worker.go`, as you did before.
2. To test your Failure Converter, change your Workflow to return an artificial
   error. Change the `ExecuteActivity` call to throw an error where there isn't
   one, like so:

   ```go
	err = workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
	if err == nil {
		err = errors.New("This is an artificial error")
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
   ```

   Don't forget to add the `errors` package to `workflow.go` as well. Next, try
   re-running your Workflow, and it should fail.
3. Run `temporal workflow show -w converters_workflowID` to get the status of your
   failed Workflow. Notice that the `Failure:` field should now display an encoded
   result, rather than a plain text error:

   ```
   ...
   Status: FAILED
   Failure: &Failure{Message:Encoded failure,Source:GoSDK,StackTrace:,Cause:nil,FailureType:Failure_ApplicationFailureInfo,}
   ```


## Part C: Implement a Composite Data Converter

1. Finally, you can implement a Composite Data Converter. A Composite Data
   Converter is used to apply custom, type-specific Payload Converters in a
   specified order. Rather than overriding the default Data Converter and
   Failure Converter as you did in the last two parts of this exercise, you can
   construct a Composite Data Converter that provies a set of rules in a custom
   order. For example, the default Go Data Converter looks like this:

   ```go
   defaultDataConverter = NewCompositeDataConverter(
      NewNilPayloadConverter(),
      NewByteSlicePayloadConverter(),
      NewProtoJSONPayloadConverter(),
      NewProtoPayloadConverter(),
      NewJSONPayloadConverter(),
   )
   ```

   Order is important. Both the `ProtoJsonPayload` and `ProtoPayload` converters
   check for the same `proto.Message` interface. The first match will always be
   used for serialization. Deserialization is controlled by metadata, therefore
   both converters can deserialize the corresponding data format (JSON or binary
   proto). For this exercise, you can try just omitting some of those converter
   interfaces, if for example you don't want your workflow to convert Nil or
   ByteSlice Payloads. Within the `client.Options{}` block, declare a
   `DataConverter` without the Nil or ByteSlice converters:

   ```go
   ClientOptions: client.Options{
		DataConverter: converter.NewCompositeDataConverter(
			converter.NewProtoJSONPayloadConverter(),
			converter.NewProtoPayloadConverter(),
			converter.NewJSONPayloadConverter(),
		),
   },
   ```

   Make this change to the `client.Options{}` block in both `starter.go` and
   `worker.go`, then restart your worker and re-run your Workflow.
2. Run `temporal workflow show -w converters_workflowID` once more to get what will
   now be stock, unencrypted output following your Composite Data Converter logic:

   ```
   Result:
     Status: COMPLETED
     Output: ["Received Plain text input"]
   ```

   You may not need to use a Composite Data Converter as often as you will need
   to override the stock Data Converter to add encryption, but if you have
   complex Workflow logic, you may need to do both.


### This is the end of the exercise.
