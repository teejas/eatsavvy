package utils

import (
	"bytes"
	"encoding/gob"
	"log/slog"
)

func ToBytes(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(v)
	if err != nil {
		slog.Error("[utils.ToBytes] Failed to encode value", "error", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func FromBytes(b []byte, v interface{}) error {
	buffer := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(v)
	if err != nil {
		slog.Error("[utils.FromBytes] Failed to decode value", "error", err)
		return err
	}
	return nil
}
