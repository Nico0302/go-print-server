package printer

import (
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/nico0302/go-ipp"
)

const (
	DefaultPort = 631
)

type Printer struct {
	Name   string
	client *ipp.IPPClient
}

type PrinterConfig struct {
	Host     string
	Port     int
	UseTLS   bool `mapstructure:"tls"`
	Username string
	Password string
}

type MediaCol struct {
	MediaSource string `mapstructure:"media-source"`
	MediaType   string `mapstructure:"media-type"`
	Media       string `mapstructure:"media"`
}

type JobAttributes struct {
	MediaCol MediaCol `mapstructure:"media-col"`
}

func NewPrinter(name string, conf PrinterConfig) *Printer {
	printer := new(Printer)
	printer.Name = name
	printer.client = ipp.NewIPPClient(conf.Host, conf.Port, conf.Username, conf.Password, conf.UseTLS)
	return printer
}

func (p Printer) PrintJob(doc Document, jobAttributes JobAttributes) (int, error) {
	attr := make(map[string]interface{})
	err := mapstructure.Decode(jobAttributes, &attr)
	if err != nil {
		return 0, err
	}
	document := ipp.Document{
		Document: doc.Body,
		Size:     doc.Size,
		Name:     doc.Name,
		MimeType: string(doc.MimeType),
	}
	id, err := p.client.PrintJob(document, p.Name, attr)
	if err != nil && err != io.EOF {
		return 0, err
	}
	return id, nil
}
