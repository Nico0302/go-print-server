package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nico0302/go-print-server/config"
	"github.com/nico0302/go-print-server/internal/printer"
)

const (
	DefaultPresetName = "default"
)

type PrinterContext struct {
	printers map[string]printer.Printer
	presets  map[string]config.Preset
}

func NewPrinterContext(cfg config.Config) *PrinterContext {
	printers := make(map[string]printer.Printer)
	for name, printerConf := range cfg.Printers {
		printers[name] = *printer.NewPrinter(name, printerConf)
	}
	return &PrinterContext{
		printers: printers,
		presets:  cfg.Presets,
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

func (c *PrinterContext) getPrinter(name string) (printer.Printer, error) {
	printer, ok := c.printers[strings.ToLower(name)]
	if !ok {
		return printer, errors.New(fmt.Sprintf("No preset found with name: '%s'", name))
	}
	return printer, nil
}
