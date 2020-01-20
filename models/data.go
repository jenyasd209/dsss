package models

import (
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"log"
)

type DataType uint8

const (
	JSON DataType = iota
	Simple
	Audio
	Video
)

type Hash32 [32]byte

func (h Hash32) String() string {
	return hex.EncodeToString(h[:])
}

type Data interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Hash() Hash32
}

type Content []byte

func (c *Content) IsEmpty() bool {
	return *c == nil
}

type AudioData struct {
	metaData `json:"meta_data"`
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
	metaData metaData
	frames   []Content
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

type DataBuilder interface {
	SetMetadata(data metaData) Data
	SetContent(content Content) Data
}