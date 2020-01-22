package models

import (
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"errors"
)

var ErrContentApplied = errors.New("data has been submitted")

type DataType uint8

const (
	Simple DataType = iota
	JSON
	Audio
	Video
)

type Hash32 [32]byte

func (h Hash32) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash32) IsEmpty() bool {
	return h == Hash32{}
}

type Data interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	CachedHash() Hash32
	Submit() error
}

type MetaData struct {
	Title      string   `json:"title"`
	cachedHash Hash32   `json:"cached_hash"`
	DataType   DataType `json:"data_type"`
}

type Content []byte

type SimpleData struct {
	MetaData
	Content Content `json:"content"`
}

func (s *SimpleData) Submit() error {
	if !s.MetaData.cachedHash.IsEmpty() {
		return ErrContentApplied
	}

	s.MetaData.cachedHash = hash(s.Content)

	return nil
}

func (s *SimpleData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

func (s *SimpleData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *SimpleData) CachedHash() Hash32 {
	return s.MetaData.cachedHash
}

type JsonData struct {
	MetaData
	Content Content `json:"body"`
}

func (jsonData *JsonData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*jsonData)
}

func (jsonData *JsonData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, jsonData)
}

func (jsonData *JsonData) CachedHash() Hash32 {
	return jsonData.MetaData.cachedHash
}

func (jsonData *JsonData) Submit() error {
	if !jsonData.MetaData.cachedHash.IsEmpty() {
		return ErrContentApplied
	}

	jsonData.MetaData.cachedHash = hash(jsonData.Content)

	return nil
}

type AudioData struct {
	MetaData
	Content Content `json:"content"`
}

func (ad *AudioData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*ad)
}

func (ad *AudioData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ad)
}

func (ad *AudioData) CachedHash() Hash32 {
	return ad.MetaData.cachedHash
}

func (ad *AudioData) Submit() error {
	if !ad.MetaData.cachedHash.IsEmpty() {
		return ErrContentApplied
	}

	ad.MetaData.cachedHash = hash(ad.Content)

	return nil
}

type VideoData struct {
	MetaData
	Frames []Content `json:"frames"`
}

func (vd *VideoData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*vd)
}

func (vd *VideoData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, vd)
}

func (vd *VideoData) CachedHash() Hash32 {
	return vd.MetaData.cachedHash
}

func (vd *VideoData) Submit() error {
	if !vd.MetaData.cachedHash.IsEmpty() {
		return ErrContentApplied
	}

	vd.MetaData.cachedHash = hash(vd.Frames[0])

	return nil
}

func hash(bytes []byte) Hash32 {
	return sha256.Sum256(bytes)
}
