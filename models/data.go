package models

import (
	"crypto/sha256"
	"encoding"
	"encoding/hex"
	"encoding/json"
)

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

	ID() Hash32
	Type() DataType
	Body() Content
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

func (sd *simpleData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(sd)
}

func (sd *simpleData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &sd)
}

func (sd *simpleData) ID() Hash32 {
	return sd.MetaData.CachedHash
}

func (sd *simpleData) Type() DataType {
	return sd.MetaData.DataType
}

func (sd *simpleData) Body() Content {
	return sd.Content
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

func (jd *jsonData) ID() Hash32 {
	return jd.MetaData.CachedHash
}

func (jd *jsonData) Type() DataType {
	return jd.MetaData.DataType
}

func (jd *jsonData) Body() Content {
	return jd.Content
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

func (ad *audioData) ID() Hash32 {
	return ad.MetaData.CachedHash
}

func (ad *audioData) Type() DataType {
	return ad.MetaData.DataType
}

func (ad *audioData) Body() Content {
	return ad.Content
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

func (vd *videoData) ID() Hash32 {
	return vd.MetaData.CachedHash
}

func (vd *videoData) Type() DataType {
	return vd.MetaData.DataType
}

func (vd *videoData) Body() Content {
	return vd.Frames
}

func hash(bytes []byte) Hash32 {
	return sha256.Sum256(bytes)
}

func NewEmptyData(dataType DataType) Data {
	switch dataType {
	case Simple:
		return NewSimpleData(
			MetaData{
				DataType: Simple,
			},
			nil,
		)
	case JSON:
		return NewJSONData(
			MetaData{
				DataType: JSON,
			},
			nil,
		)
	case Audio:
		return NewJSONData(
			MetaData{
				DataType: Audio,
			},
			nil,
		)
	case Video:
		return NewJSONData(
			MetaData{
				DataType: Video,
			},
			nil,
		)
	default:
		return nil
	}
}

func NewDataWithTitle(dataType DataType, title string) Data {
	switch dataType {
	case Simple:
		return NewSimpleData(
			MetaData{
				Title:    title,
				DataType: Simple,
			},
			nil,
		)
	case JSON:
		return NewJSONData(
			MetaData{
				Title:    title,
				DataType: JSON,
			},
			nil,
		)
	case Audio:
		return NewJSONData(
			MetaData{
				Title:    title,
				DataType: Audio,
			},
			nil,
		)
	case Video:
		return NewJSONData(
			MetaData{
				Title:    title,
				DataType: Video,
			},
			nil,
		)
	default:
		return nil
	}
}

func NewDataWithContent(dataType DataType, content Content) Data {
	switch dataType {
	case Simple:
		return NewSimpleData(
			MetaData{
				DataType: Simple,
			},
			content,
		)
	case JSON:
		return NewJSONData(
			MetaData{
				DataType: JSON,
			},
			content,
		)
	case Audio:
		return NewJSONData(
			MetaData{
				DataType: Audio,
			},
			content,
		)
	case Video:
		return NewJSONData(
			MetaData{
				DataType: Video,
			},
			content,
		)
	default:
		return nil
	}
}

func NewData(dataType DataType, title string, content Content) Data {
	switch dataType {
	case Simple:
		return NewSimpleData(
			MetaData{
				Title:    title,
				DataType: Simple,
			},
			content,
		)
	case JSON:
		return NewJSONData(
			MetaData{
				Title:    title,
				DataType: JSON,
			},
			content,
		)
	case Audio:
		return NewJSONData(
			MetaData{
				Title:    title,
				DataType: Audio,
			},
			content,
		)
	case Video:
		return NewJSONData(
			MetaData{
				Title:    title,
				DataType: Video,
			},
			content,
		)
	default:
		return nil
	}
}
