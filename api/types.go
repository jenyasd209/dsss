package api

import (
	"github.com/iorhachovyevhen/dsss/models"
	"path/filepath"
)

var fileType = map[string]models.DataType{
	".json": models.JSON,
	"":      models.Simple,
	".txt":  models.Simple,
	".doc":  models.Simple,
	".docx": models.Simple,
	".odt":  models.Simple,
	".pdf":  models.Simple,
	".rtf":  models.Simple,
	".wps":  models.Simple,
	".mp3":  models.Audio,
	".mpa":  models.Audio,
	".ogg":  models.Audio,
	".wav":  models.Audio,
	".wma":  models.Audio,
	".3g2":  models.Video,
	".3gp":  models.Video,
	".avi":  models.Video,
	".flv":  models.Video,
	".h264": models.Video,
	".mp4":  models.Video,
	".mpg":  models.Video,
	".swf":  models.Video,
	".wmv":  models.Video,
}

var contentType = map[string]string{
	".json": "application/json",
	"":      "text/plain",
	".txt":  "text/plain",
	".doc":  "application/msword",
	".docx": "application/msword",
	".odt":  "application/vnd.oasis.opendocument.text",
	".pdf":  "application/pdf",
	".rtf":  "application/rtf",
	".wps":  "application/xml",
	".mp3":  "audio/mpeg",
	".mpa":  "audio/mpeg",
	".ogg":  "application/ogg",
	".wav":  "audio/x-wav",
	".wma":  "audio/x-ms-wma",
	".3g2":  "video/3gpp2",
	".3gp":  "video/3gpp",
	".avi":  "video/x-msvideo",
	".flv":  "video/x-flv",
	".h264": "video/H264",
	".mp4":  "video/mp4",
	".swf":  "application/x-shockwave-flash",
	".wmv":  "video/x-ms-wmv",
}

func DataTypeFromFilename(filename string) models.DataType {
	ext := filepath.Ext(filename)

	dt, ok := fileType[ext]
	if !ok {
		dt = models.Simple
	}

	return dt
}

func ContentTypeFromName(filename string) string {
	ext := filepath.Ext(filename)

	ct, ok := contentType[ext]
	if !ok {
		ct = "text/plain"
	}

	return ct
}
