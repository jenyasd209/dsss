package models

import (
	"encoding"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type DataType uint8

const (
	Simple DataType = iota
	JSON
	Audio
	Video
)

type Prefix string

const (
	PrefixSimple Prefix = "simple"
	PrefixJSON   Prefix = "json"
	PrefixAudio  Prefix = "audio"
	PrefixVideo  Prefix = "video"
)

var DataPrefixMap = map[DataType]Prefix{
	Simple: PrefixSimple,
	JSON:   PrefixJSON,
	Audio:  PrefixAudio,
	Video:  PrefixVideo,
}

var DataTypeMap = map[Prefix]DataType{
	PrefixSimple: Simple,
	PrefixJSON:   JSON,
	PrefixAudio:  Audio,
	PrefixVideo:  Video,
}

var (
	ErrorBadDataType       = errors.New("can't convert to DataType")
	ErrorGetDataTypeFromID = errors.New("can't get data type from ID")
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

	ID() ID
	Type() DataType
	Body() Content
	Title() string
}

type ID []byte

func (id ID) String() string {
	return string(id)
}

type MetaData struct {
	Title    string   `json:"title"`
	ID       ID       `json:"id"`
	DataType DataType `json:"data_type"`
}

type Content []byte

func NewSimpleData(metadata MetaData, content Content) (sd *simpleData) {
	sd = &simpleData{
		MetaData: metadata,
		Content:  content,
	}

	h := hash(sd.Content)
	sd.MetaData.ID = composeID(h, sd.MetaData.DataType)

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

func (sd *simpleData) ID() ID {
	return sd.MetaData.ID
}

func (sd *simpleData) Type() DataType {
	return sd.MetaData.DataType
}

func (sd *simpleData) Body() Content {
	return sd.Content
}

func (sd *simpleData) Title() string {
	return sd.MetaData.Title
}

func NewJSONData(metadata MetaData, content Content) (jd *jsonData) {
	jd = &jsonData{
		MetaData: metadata,
		Content:  content,
	}

	h := hash(jd.Content)
	jd.MetaData.ID = composeID(h, jd.DataType)

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

func (jd *jsonData) ID() ID {
	return jd.MetaData.ID
}

func (jd *jsonData) Type() DataType {
	return jd.MetaData.DataType
}

func (jd *jsonData) Body() Content {
	return jd.Content
}

func (jd *jsonData) Title() string {
	return jd.MetaData.Title
}

func NewAudioData(metadata MetaData, content Content) (ad *audioData) {
	ad = &audioData{
		MetaData: metadata,
		Content:  content,
	}

	h := hash(ad.Content)
	ad.MetaData.ID = composeID(h, ad.DataType)

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

func (ad *audioData) ID() ID {
	return ad.MetaData.ID
}

func (ad *audioData) Type() DataType {
	return ad.MetaData.DataType
}

func (ad *audioData) Body() Content {
	return ad.Content
}

func (ad *audioData) Title() string {
	return ad.MetaData.Title
}

func NewVideoData(metadata MetaData, frames Content) (vd *videoData) {
	vd = &videoData{
		MetaData: metadata,
		Frames:   frames,
	}

	h := hash(vd.Frames)
	vd.MetaData.ID = composeID(h, vd.DataType)

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

func (vd *videoData) ID() ID {
	return vd.MetaData.ID
}

func (vd *videoData) Type() DataType {
	return vd.MetaData.DataType
}

func (vd *videoData) Body() Content {
	return vd.Frames
}

func (vd *videoData) Title() string {
	return vd.MetaData.Title
}
