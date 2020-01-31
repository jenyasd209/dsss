package server

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/iorhachovyevhen/dsss/models"
)

func byteToHash32(b []byte) (h models.Hash32, err error) {
	d, err := hex.DecodeString(string(b))
	copy(h[:], d[:])

	return
}

func byteToDataType(b []byte) models.DataType {
	u, _ := binary.Uvarint(b)
	return models.DataType(u)
}

func newData(dataType []byte) models.Data {
	dt := byteToDataType(dataType)

	switch dt {
	case models.Simple:
		return models.NewSimpleData(
			models.MetaData{
				DataType: models.Simple,
			},
			nil,
		)
	case models.JSON:
		return models.NewJSONData(
			models.MetaData{
				DataType: models.JSON,
			},
			nil,
		)
	case models.Audio:
		return models.NewJSONData(
			models.MetaData{
				DataType: models.Audio,
			},
			nil,
		)
	case models.Video:
		return models.NewJSONData(
			models.MetaData{
				DataType: models.Video,
			},
			nil,
		)
	default:
		return nil
	}
}
