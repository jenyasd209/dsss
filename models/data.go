package models

import (
	"crypto/sha256"
	"encoding"
	"encoding/json"
)

type Hash32 [32]byte

type Data interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Hash() Hash32
}

type SimpleData struct {
	Content []byte `json:"content"`
}

func (s *SimpleData) Hash() Hash32{
	return sha256.Sum256(s.Content)
}

func (s *SimpleData) MarshalBinary() (data []byte, err error){
	return json.Marshal(s)
}

func (s *SimpleData) UnmarshalBinary(data []byte) error  {
	return json.Unmarshal(data, s)
}

type JsonData struct {
	SimpleData
}

type AudioData struct {
	SimpleData
}

type VideoData struct {
	Frames []SimpleData
}