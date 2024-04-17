package compositeconverter

import (
	"bytes"
	"encoding/json"
	"fmt"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

type CustomDataConverter struct {
}

// CustomPayloadConverter creates a payload converter
func NewCustomPayloadConverter() converter.PayloadConverter {
	return &CustomDataConverter{}
}

// Payload converter implementation

func (dc *CustomDataConverter) ToPayloads(value ...interface{}) (*commonpb.Payloads, error) {
	payloads := &commonpb.Payloads{}
	for _, obj := range value {
		payload, err := dc.ToPayload(obj)
		if err != nil {
			return nil, err
		}

		payloads.Payloads = append(payloads.Payloads, payload)
	}

	return payloads, nil
}

func (dc *CustomDataConverter) ToPayload(value interface{}) (*commonpb.Payload, error) {
	var err error
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err = enc.Encode(value)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to encode argument: %T, with error: %w", value, err)
	}

	payload := &commonpb.Payload{
		Metadata: map[string][]byte{
			"encoding": []byte("raw"),
		},
		Data: buf.Bytes(),
	}

	return payload, nil
}

func (dc *CustomDataConverter) FromPayloads(payloads *commonpb.Payloads, valuePtrs ...interface{}) error {
	if payloads == nil {
		return nil
	}
	for i, payload := range payloads.Payloads {
		err := dc.FromPayload(payload, valuePtrs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (dc *CustomDataConverter) FromPayload(payload *commonpb.Payload, valuePtr interface{}) error {
	var err error
	dec := json.NewDecoder(bytes.NewBuffer(payload.GetData()))
	err = dec.Decode(valuePtr)
	if err != nil {
		return fmt.Errorf(
			"unable to decode argument: %T, with error: %v", valuePtr, err)
	}
	return nil
}

func (dc *CustomDataConverter) ToString(payload *commonpb.Payload) string {
	return string(payload.GetData())
}

func (c *CustomDataConverter) Encoding() string {
	return "json/plain"
}
