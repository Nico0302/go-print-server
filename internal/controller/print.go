package controller

import (
	"bytes"
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
		h.POST("/urls", r.urls)
	}
}

func (r *printRoutes) printDocument(doc printer.Document, presetName string) (int, error) {
	preset, err := r.c.getPreset(presetName)
	if err != nil {
		return 0, err
	}
	p, err := r.c.getPrinter(preset.Printer)
	if err != nil {
		return 0, err
	}

	r.l.Debug(fmt.Sprintf("Print file %s on printer %s with preset %s.", doc.Name, p.GetName(), presetName), "http - print - url")

	id, err := p.PrintJob(doc, preset.JobAttributes)
	if err != nil {
		return 0, fmt.Errorf("IPP error: %w", err)
	}

	return id, nil
}

type document struct {
	Doc    printer.Document
	Preset string
}

func (r *printRoutes) printDocuments(docs []document) ([]int, error) {
	var ids []int
	for _, doc := range docs {
		id, err := r.printDocument(doc.Doc, doc.Preset)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *printRoutes) downloadFile(url string, cookies string) (*printer.Document, error) {
	file, err := r.f.DownloadFile(url, cookies)
	if err != nil {
		return nil, err
	}

	doc := printer.Document{
		Name:     "CloudPrintDocument",
		MimeType: printer.ApplicationPdf,
	}

	if file.ContentLength > 0 {
		doc.Body = file.Body
		doc.Size = int(file.ContentLength)
	} else {
		var buf bytes.Buffer
		_, err = buf.ReadFrom(file.Body)
		if err != nil {
			return nil, err
		}
		doc.Body = &buf
		doc.Size = buf.Len()
	}

	return &doc, nil
}

type urlPrintResponse struct {
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

	printDoc, err := r.downloadFile(request.Url, cookies)
	if err != nil {
		r.l.Error(err, "http - print - url")
		errorResponse(c, http.StatusInternalServerError, "could not download file")
		return
	}

	id, err := r.printDocument(*printDoc, request.Preset)
	if err != nil {
		r.l.Error(err, "http - print - url")
		errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("print error: %s", err))
		return
	}

	c.JSON(http.StatusOK, urlPrintResponse{
		JobID: id,
	})
}

type urlDocument struct {
	Url    string
	Preset string
}

type urlsPrintResponse struct {
	JobIDs []int
}

type urlsPrintRequest struct {
	Documents []urlDocument
}

func (r *printRoutes) urls(c *gin.Context) {
	cookies := c.GetHeader(fetcher.ForwardedCookieHeader)

	var request urlsPrintRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - print - url")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	var jobIDs []int
	for _, doc := range request.Documents {
		printDoc, err := r.downloadFile(doc.Url, cookies)
		if err != nil {
			r.l.Error(err, "http - print - url", doc.Url)
			errorResponse(c, http.StatusInternalServerError, "could not download file")
			return
		}

		jobID, err := r.printDocument(*printDoc, doc.Preset)
		if err != nil {
			r.l.Error(err, "http - print - url")
			errorResponse(c, http.StatusInternalServerError, fmt.Sprintf("print error: %s", err))
			return
		}

		jobIDs = append(jobIDs, jobID)
	}

	c.JSON(http.StatusOK, urlsPrintResponse{
		JobIDs: jobIDs,
	})
}
