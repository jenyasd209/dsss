package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func hash(bytes []byte) Hash32 {
	return sha256.Sum256(bytes)
}

func ConvertToDataType(value interface{}) (dataType DataType, err error) {
	switch v := value.(type) {
	case float64:
		dataType = DataType(v)
	case string:
		return strToDataType(v)
	case int:
		return intToDataType(v)
	case uint:
		dataType = DataType(v)
	default:
		return 0, ErrorBadDataType
	}

	return
}

func strToDataType(s string) (DataType, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrorBadDataType
	}

	return intToDataType(i)
}

func intToDataType(i int) (DataType, error) {
	if i < 0 {
		return 0, ErrorBadDataType
	}

	return DataType(i), nil
}

func DataTypeFromID(hexID []byte) (DataType, error) {
	id, err := hex.DecodeString(string(hexID))
	if err != nil {
		return 0, ErrorGetDataTypeFromID
	}

	prefix := id[:len(id)-32]

	dt, ok := DataTypeMap[Prefix(prefix)]
	if !ok {
		return dt, ErrorGetDataTypeFromID
	}

	return dt, nil
}

func composeID(hash32 Hash32, dataType DataType) ID {
	prefix := []byte(DataPrefixMap[dataType])

	var id []byte
	id = append(id, prefix...)
	id = append(id, hash32[:]...)

	return ID(hex.EncodeToString(id))
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
		return NewAudioData(
			MetaData{
				DataType: Audio,
			},
			nil,
		)
	case Video:
		return NewVideoData(
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
		return NewAudioData(
			MetaData{
				Title:    title,
				DataType: Audio,
			},
			nil,
		)
	case Video:
		return NewVideoData(
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
		return NewAudioData(
			MetaData{
				DataType: Audio,
			},
			content,
		)
	case Video:
		return NewVideoData(
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
		return NewAudioData(
			MetaData{
				Title:    title,
				DataType: Audio,
			},
			content,
		)
	case Video:
		return NewVideoData(
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
