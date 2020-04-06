package models

import (
	"encoding"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Content []byte

type DataType uint8

const (
	Unknown DataType = iota
	Simple
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
	ErrorMetaDataIsNil     = errors.New("metadata is nil")
	ErrorContentIsNil      = errors.New("metadata is nil")
)

type Data interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Meta() *MetaData
	Body() *Content
}

type ID []byte

func (id ID) String() string {
	return hex.EncodeToString(id)
}

func NewMetaData(title string, dataType DataType) *MetaData {
	return &MetaData{
		Title:    title,
		ID:       newID(dataType),
		DataType: dataType,
	}
}

type MetaData struct {
	Title    string   `json:"title"`
	ID       ID       `json:"id"`
	DataType DataType `json:"data_type"`
}

func (m *MetaData) GetID() ID {
	return m.ID
}

func (m *MetaData) GetType() DataType {
	return m.DataType
}

func (m *MetaData) GetTitle() string {
	return m.Title
}

func NewSimpleData(metadata *MetaData, content Content) *simpleData {
	return &simpleData{
		MetaData: *metadata,
		Content:  content,
	}
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

func (sd *simpleData) Meta() *MetaData {
	return &sd.MetaData
}

func (sd *simpleData) Body() *Content {
	return &sd.Content
}

func NewJSONData(metadata *MetaData, content Content) *jsonData {
	return &jsonData{
		MetaData: *metadata,
		Content:  content,
	}
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

func (jd *jsonData) Meta() *MetaData {
	return &jd.MetaData
}

func (jd *jsonData) Body() *Content {
	return &jd.Content
}

func NewAudioData(metadata *MetaData, content Content) *audioData {
	return &audioData{
		MetaData: *metadata,
		Content:  content,
	}
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

func (ad *audioData) Meta() *MetaData {
	return &ad.MetaData
}

func (ad *audioData) Body() *Content {
	return &ad.Content
}

func NewVideoData(metadata *MetaData, content Content) *videoData {
	return &videoData{
		MetaData: *metadata,
		Content:  content,
	}
}

type videoData struct {
	MetaData
	Content Content `json:"content"`
}

func (vd *videoData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*vd)
}

func (vd *videoData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, vd)
}

func (vd *videoData) Meta() *MetaData {
	return &vd.MetaData
}

func (vd *videoData) Body() *Content {
	return &vd.Content
}

//NewData creates new Data. If some parameter is nil than Data will be nil too
func NewData(metaData *MetaData, content Content) (Data, error) {
	if metaData == nil {
		return nil, ErrorMetaDataIsNil
	}

	if content == nil {
		return nil, ErrorContentIsNil
	}

	switch metaData.DataType {
	case Simple:
		return NewSimpleData(
			metaData,
			content,
		), nil
	case JSON:
		return NewJSONData(
			metaData,
			content,
		), nil
	case Audio:
		return NewAudioData(
			metaData,
			content,
		), nil
	case Video:
		return NewVideoData(
			metaData,
			content,
		), nil
	default:
		return nil, errors.New("dataType is wrong")
	}
}
