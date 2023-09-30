package printer

import (
	"github.com/mitchellh/mapstructure"
	"github.com/nico0302/go-ipp"
)

type CupsClient struct {
	name   string
	client *ipp.CUPSClient
}

type CupsConfig struct {
	Host     string
	Port     int
	UseTLS   bool `mapstructure:"tls"`
	Username string
	Password string
}

type CupsPrinter struct {
	name   string
	client *ipp.CUPSClient
}

type CupsPrinterConfig struct {
	CupsName    string `mapstructure:"cups"`
	PrinterName string `mapstructure:"name"`
}

func NewCupsPrinter(name string, cups *ipp.CUPSClient) *CupsPrinter {
	printer := new(CupsPrinter)
	printer.name = name
	printer.client = cups
	return printer
}

func (c *CupsPrinter) PrintJob(doc Document, jobAttributes JobAttributes) (int, error) {
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
	id, err := c.client.PrintJob(document, c.name, attr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CupsPrinter) GetName() string {
	return c.name
}
