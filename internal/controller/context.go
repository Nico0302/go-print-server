package controller

import (
	"errors"
	"fmt"
	"github.com/nico0302/go-ipp"
	"strings"

	"github.com/nico0302/go-print-server/config"
	"github.com/nico0302/go-print-server/internal/printer"
)

const (
	DefaultPresetName = "default"
)

type PrinterContext struct {
	printers map[string]printer.IPrinter
	cups     map[string]*ipp.CUPSClient
	presets  map[string]config.Preset
}

func NewPrinterContext(cfg config.Config) *PrinterContext {
	printers := make(map[string]printer.IPrinter)
	cups := make(map[string]*ipp.CUPSClient)
	for name, printerConf := range cfg.Printers {
		if printerConf.Host != "" {
			printers[name] = printer.NewPrinter(name, printer.PrinterConfig{
				Host:     printerConf.Host,
				Port:     printerConf.Port,
				UseTLS:   printerConf.UseTLS,
				Username: printerConf.Username,
				Password: printerConf.Password,
			})
		} else if printerConf.CupsName != "" {
			if cups[printerConf.CupsName] == nil {
				conf := cfg.Cups[printerConf.CupsName]
				cups[printerConf.CupsName] = ipp.NewCUPSClient(conf.Host, conf.Port, conf.Username, conf.Password, conf.UseTLS)
			}
			printers[name] = printer.NewCupsPrinter(printerConf.PrinterName, cups[printerConf.CupsName])
		}
	}
	return &PrinterContext{
		printers: printers,
		presets:  cfg.Presets,
		cups:     cups,
	}
}

func (c *PrinterContext) getPreset(name string) (config.Preset, error) {
	if name == "" {
		preset, ok := c.presets[DefaultPresetName]
		if !ok {
			return preset, errors.New("No default preset specified.")
		}
		return preset, nil
	}

	preset, ok := c.presets[strings.ToLower(name)]
	if !ok {
		return preset, errors.New(fmt.Sprintf("No preset found with name: '%s'", name))
	}
	return preset, nil
}

func (c *PrinterContext) getPrinter(name string) (printer.IPrinter, error) {
	p, ok := c.printers[strings.ToLower(name)]
	if !ok {
		return p, errors.New(fmt.Sprintf("No preset found with name: '%s'", name))
	}
	return p, nil
}
