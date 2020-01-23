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
	IsCorrect() bool
	Type() DataType
}

type MetaData struct {
	Title      string   `json:"title"`
	CachedHash Hash32   `json:"cached_hash"`
	DataType   DataType `json:"data_type"`
}

type Content []byte

func NewSimpleData(metadata MetaData, content Content) (sd *simpleData) {
	sd = &simpleData{
		MetaData: metadata,
		Content:  content,
	}

	sd.MetaData.CachedHash = hash(sd.Content)

	return
}

type simpleData struct {
	MetaData
	Content Content `json:"content"`
}

func (sd *simpleData) IsCorrect() bool {
	return sd.MetaData.CachedHash == hash(sd.Content)
}

func (sd *simpleData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(sd)
}

func (sd *simpleData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &sd)
}

func (sd *simpleData) CachedHash() Hash32 {
	return sd.MetaData.CachedHash
}

func (sd *simpleData) Type() DataType {
	return sd.MetaData.DataType
}

func NewJSONData(metadata MetaData, content Content) (jd *jsonData) {
	jd = &jsonData{
		MetaData: metadata,
		Content:  content,
	}

	jd.MetaData.CachedHash = hash(jd.Content)

	return
}

type jsonData struct {
	MetaData
	Content Content `json:"body"`
}

func (jd *jsonData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*jd)
}

func (jd *jsonData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, jd)
}

func (jd *jsonData) CachedHash() Hash32 {
	return jd.MetaData.CachedHash
}

func (jd *jsonData) IsCorrect() bool {
	return jd.MetaData.CachedHash == hash(jd.Content)
}

func (jd *jsonData) Type() DataType {
	return jd.MetaData.DataType
}

func NewAudioData(metadata MetaData, content Content) (ad *audioData) {
	ad = &audioData{
		MetaData: metadata,
		Content:  content,
	}

	ad.MetaData.CachedHash = hash(ad.Content)

	return ad
}

type audioData struct {
	MetaData
	Content Content `json:"content"`
}

func (ad *audioData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*ad)
}

func (ad *audioData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ad)
}

func (ad *audioData) CachedHash() Hash32 {
	return ad.MetaData.CachedHash
}

func (ad *audioData) IsCorrect() bool {
	return ad.MetaData.CachedHash == hash(ad.Content)
}

func (ad *audioData) Type() DataType {
	return ad.MetaData.DataType
}

func NewVideoData(metadata MetaData, frames Content) (vd *videoData) {
	vd = &videoData{
		MetaData: metadata,
		Frames:   frames,
	}

	vd.MetaData.CachedHash = hash(vd.Frames)

	return
}

type videoData struct {
	MetaData
	Frames Content `json:"frames"`
}

func (vd *videoData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*vd)
}

func (vd *videoData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, vd)
}

func (vd *videoData) CachedHash() Hash32 {
	return vd.MetaData.CachedHash
}

func (vd *videoData) IsCorrect() bool {
	return vd.MetaData.CachedHash == hash(vd.Frames)
}

func (vd *videoData) Type() DataType {
	return vd.MetaData.DataType
}

func hash(bytes []byte) Hash32 {
	return sha256.Sum256(bytes)
}
