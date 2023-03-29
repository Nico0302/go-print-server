package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nico0302/go-print-server/internal/fetcher"
	"github.com/nico0302/go-print-server/internal/printer"
	"github.com/nico0302/go-print-server/pkg/logger"
)

type printRoutes struct {
	c *PrinterContext
	f *fetcher.Fetcher
	l logger.Interface
}

func newPrintRoutes(handler *gin.RouterGroup, c *PrinterContext, f *fetcher.Fetcher, l logger.Interface) {
	r := &printRoutes{c, f, l}

	h := handler.Group("/print")
	{
		h.POST("/url", r.url)
	}
}

func (r *printRoutes) printDocument(doc printer.Document, presetName string) (int, error) {
	preset, err := r.c.getPreset(presetName)
	if err != nil {
		return 0, err
	}
	printer, err := r.c.getPrinter(preset.Printer)
	if err != nil {
		return 0, err
	}

	r.l.Debug(fmt.Sprintf("Print file %s on printer %s with preset %s.", doc.Name, printer.Name, presetName), "http - print - url")

	id, err := printer.PrintJob(doc, preset.JobAttributes)
	if err != nil {
		return 0, fmt.Errorf("IPP error: %w", err)
	}

	return id, nil
}

type printResponse struct {
	JobID int
}

type urlPrintRequest struct {
	Url    string
	Preset string
}

func (r *printRoutes) url(c *gin.Context) {
	cookies := c.GetHeader(fetcher.ForwardedCookieHeader)

	var request urlPrintRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - print - url")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	file, err := r.f.DownloadFile(request.Url, cookies)
	if err != nil {
		r.l.Error(err, "http - print - url", request.Url)
		errorResponse(c, http.StatusInternalServerError, "could not download file")
		return
	}

	doc := printer.Document{
		Body:     file.Body,
		Size:     fetcher.GetFileSize(file),
		MimeType: printer.ApplicationPdf,
	}

	id, err := r.printDocument(doc, request.Preset)
	if err != nil {
		r.l.Error(err, "http - print - url")
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("print error: %s", err))
		return
	}

	c.JSON(http.StatusOK, printResponse{
		JobID: id,
	})
}
