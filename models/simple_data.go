package models

import (
	"crypto/sha256"
	"encoding/json"
)

func NewSimpleData() *simpleDataBuilder {
	return &simpleDataBuilder{}
}

type simpleData struct {
	metaData
	content  Content
}

type simpleDataBuilder struct {
	simpleData
}

func (s *simpleDataBuilder) Build() simpleData {
	return s.simpleData
}

func (s *simpleDataBuilder) SetMetadata(data metaData) *simpleDataBuilder {
	s.simpleData.metaData = data
	s.simpleData.dataType = Simple
	return s
}

func (s *simpleDataBuilder) SetContent(content Content) *simpleDataBuilder {
	if !s.content.IsEmpty() {
		return s
	}

	s.simpleData.content = content
	s.metaData.setCashedHash(s.Hash())

	return s
}

func (s *simpleData) GetMetadata() *metaData {
	return &s.metaData
}

func (s *simpleData) GetContent() *Content {
	return &s.content
}

func (s *simpleData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(&struct {
		Title string
		CashedHash Hash32
		DataType DataType
		Content Content
	}{
		Title: s.metaData.title,
		CashedHash: s.metaData.cashedHash,
		DataType: s.metaData.dataType,
		Content: s.content,
	})
}

func (s *simpleData) UnmarshalBinary(data []byte) error {
	anon := struct {
		Title string
		CashedHash Hash32
		DataType DataType
		Content Content
	}{}

	err := json.Unmarshal(data, &anon)
	if err != nil {
		return err
	}

	s.title = anon.Title
	s.cashedHash = anon.CashedHash
	s.dataType = anon.DataType
	s.content = anon.Content

	return nil
}

func (s *simpleData) Hash() Hash32 {
	return sha256.Sum256(s.content)
}
