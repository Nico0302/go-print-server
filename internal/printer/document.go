package printer

import (
	"io"
)

type MimeType string

const (
	ApplicationPdf MimeType = "application/pdf"
)

type Document struct {
	Body     io.Reader
	Size     int
	Name     string
	MimeType MimeType
}
