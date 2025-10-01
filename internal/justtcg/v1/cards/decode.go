package cards

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type DecodeOption func(*json.Decoder)

func UseNumber() DecodeOption {
	return func(d *json.Decoder) {
		d.UseNumber()
	}
}

func Decode(data []byte, opts ...DecodeOption) (ResponseDTO, error) {
	var dto ResponseDTO

	dec := json.NewDecoder(bytes.NewReader(data))
	for _, o := range opts {
		o(dec)
	}

	err := dec.Decode(&dto)
	if err != nil {
		return ResponseDTO{}, fmt.Errorf("decode: responseDTO error: %w", err)
	}

	if dec.More() {
		return ResponseDTO{}, fmt.Errorf("decode: extra data after JSON")
	}
	return dto, nil
}
