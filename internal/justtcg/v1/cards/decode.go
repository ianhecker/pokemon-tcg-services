package cards

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Decode(data []byte) (ResponseDTO, error) {

	var dto ResponseDTO
	decoder := json.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&dto)
	if err != nil {
		return ResponseDTO{}, fmt.Errorf("decode: responseDTO error: %w", err)
	}

	if decoder.More() {
		return ResponseDTO{}, fmt.Errorf("decode: extra data after JSON")
	}
	return dto, nil
}
