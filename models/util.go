package models

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"strconv"
)

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
		return Unknown, ErrorBadDataType
	}

	return
}

func strToDataType(s string) (DataType, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return Unknown, ErrorBadDataType
	}

	return intToDataType(i)
}

func intToDataType(i int) (DataType, error) {
	if i < 0 {
		return Unknown, ErrorBadDataType
	}

	return DataType(i), nil
}

func DataTypeFromID(hexID []byte) (DataType, error) {
	id := make([]byte, hex.DecodedLen(len(hexID)))
	_, err := hex.Decode(id, hexID)
	if err != nil {
		return Unknown, ErrorBadHexID
	}

	prefix := id[:len(id)-32]

	dt, ok := DataTypeMap[Prefix(prefix)]
	if !ok {
		return Unknown, ErrorGetDataTypeFromID
	}

	return dt, nil
}

func newID(dataType DataType) ID {
	uuid := make([]byte, 32)
	_, err := rand.Read(uuid)
	if err != nil {
		log.Fatal(err)
	}

	return composeID(uuid, dataType)
}

func composeID(uuid []byte, dataType DataType) ID {
	prefix := []byte(DataPrefixMap[dataType])

	id := append(prefix, uuid...)

	hexID := make([]byte, hex.EncodedLen(len(id)))
	hex.Encode(hexID, id)

	return hexID
}
