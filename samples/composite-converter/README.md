# Customizing the Composite Data Converter

This example shows how you can customize the Composite Data Converter by adding an additional Payload Converter. This sample uses the same logic as Exercise #1, but instead of providing a Custom Codec, it provides a new Payload Converter.

## Part A: Add an Additional Payload Converter to the Composite Converter

In `starter/main.go`, in the `main()` function before creating your `client()`, a new Composite Data Converter is initialized as follows:

```go
dataConverter := converter.NewCompositeDataConverter(
    converter.NewNilPayloadConverter(),
    converter.NewByteSlicePayloadConverter(),
    converter.NewProtoJSONPayloadConverter(),
    converter.NewProtoPayloadConverter(),
    compositeconverter.NewCustomPayloadConverter(),
    converter.NewJSONPayloadConverter(),
)
```

This matches the sequence of the default Composite Data Converter, while adding a `NewCustomPayloadConverter()` as the second-to-last step before the `NewJSONPayloadConverter()`. The `NewJSONPayloadConverter()` is a fallback that serializes most values (eg strings) that are not caught by the preceding Nil/ByteSlice/Protobuf converters, so this is generally a good place to add a new Payload Converter.

## Part B: Using your customized Composite Converter

As with the Custom Codec example, you can customize your data conversion behavior by adding a `DataConverter` parameter to your `client.Options()` field.

```go
c, err := client.Dial(client.Options{
    // Set DataConverter here
    DataConverter: dataConverter,
})
```

## Part C: Defining the Behavior of your new Payload Converter

In this sample, `data_converter.go` contains a complete implementation of a new Payload Converter. As mentioned, it incorporates `ToPayloads` and `FromPayloads` calls which perform the actual variable handling and serialization/deserialization to cast your input into a format that Temporal can store. The logic of this custom Payload Converter is currently an near-exact replica of the stock JSON converter, so by adding this `NewCustomPayloadConverter`, you should not change the stock behavior of Temporal's Go SDK at all. However, it should provide a basis for doing so if you need to accommodate any other data types.

### This is the end of the sample.