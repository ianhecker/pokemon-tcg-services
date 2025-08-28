package cards

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Response struct {
	Data []Card `json:"data"`
}

func (response *Response) GetCardIndex(i int) (Card, error) {
	if response == nil {
		return Card{}, errors.New("response is nil")
	}
	length := len(response.Data)
	if length == 0 {
		return Card{}, errors.New("zero cards in response")
	}
	if length < i {
		return Card{}, fmt.Errorf("index: %d out of bounds of array length: %d", i, length)
	}
	return response.Data[i], nil
}

func (response *Response) UnmarshalJSON(bytes []byte) error {
	type Alias Response
	tmp := &struct {
		*Alias
	}{
		Alias: (*Alias)(response),
	}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return fmt.Errorf("error unmarshaling pricing response: %w", err)
	}
	if len(response.Data) == 0 {
		return fmt.Errorf("zero cards in response")
	}
	return nil
}
