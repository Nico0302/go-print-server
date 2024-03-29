package printer

import (
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/nico0302/go-ipp"
)

type IPrinter interface {
	GetName() string
	PrintJob(doc Document, jobAttributes JobAttributes) (int, error)
}

type Printer struct {
	name   string
	client *ipp.IPPClient
}

type PrinterConfig struct {
	Host     string
	Port     int
	UseTLS   bool `mapstructure:"tls"`
	Username string
	Password string
}

func NewPrinter(name string, conf PrinterConfig) *Printer {
	printer := new(Printer)
	printer.name = name
	printer.client = ipp.NewIPPClient(conf.Host, conf.Port, conf.Username, conf.Password, conf.UseTLS)
	ipp.NewCUPSClient(conf.Host, conf.Port, conf.Username, conf.Password, conf.UseTLS)
	return printer
}

func (p *Printer) PrintJob(doc Document, jobAttributes JobAttributes) (int, error) {
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
	id, err := p.client.PrintJob(document, p.name, attr)
	if err != nil && err != io.EOF {
		return 0, err
	}
	return id, nil
}

func (p *Printer) GetName() string {
	return p.name
}
