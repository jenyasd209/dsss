package models

import (
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"log"
)

type Hash32 [32]byte

func (h Hash32) String() string {
	bytes := h[:]

	return hex.EncodeToString(bytes)
}

type Data interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Hash() Hash32
}

type Content []byte

type MetaData struct {
	Title  string `json:"title"`
	Format string `json:"format"`
}

type SimpleData struct {
	MetaData `json:"meta_data"`
	Content  `json:"content"`
}

func (s *SimpleData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*s)
}

func (s *SimpleData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *SimpleData) Hash() Hash32 {
	data, err := s.MarshalBinary()
	if err != nil {
		log.Println(err)
	}

	return sha256.Sum256(data)
}

type JsonData struct {
	Content
}

func (jsonData *JsonData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*jsonData)
}

func (jsonData *JsonData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, jsonData)
}

func (jsonData *JsonData) Hash() Hash32 {
	data, err := jsonData.MarshalBinary()
	if err != nil {
		log.Println(err)
	}

	return sha256.Sum256(data)
}

type AudioData struct {
	MetaData `json:"meta_data"`
	Content  `json:"content"`
}

func (ad *AudioData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*ad)
}

func (ad *AudioData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ad)
}

func (ad *AudioData) Hash() Hash32 {
	data, err := ad.MarshalBinary()
	if err != nil {
		log.Println(err)
	}

	return sha256.Sum256(data)
}

type VideoData struct {
	MetaData `json:"meta_data"`
	Frames   []Content `json:"frames"`
}

func (vd *VideoData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*vd)
}

func (vd *VideoData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, vd)
}

func (vd *VideoData) Hash() Hash32 {
	data, err := vd.MarshalBinary()
	if err != nil {
		log.Println(err)
	}

	return sha256.Sum256(data)
}
