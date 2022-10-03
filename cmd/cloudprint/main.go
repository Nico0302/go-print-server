package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/phin1x/go-ipp"
)

var (
	printcln    *ipp.IPPClient
	printername string
	httpcln     *http.Client
)

type UrlPrintReq struct {
	URL           string `json:"url"`
	Preset        string `json:"preset"`
	ForwardCookie string `json:"forwardCookie"`
}

func urlPrint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("NEW REQUEST")
	d := json.NewDecoder(req.Body)
	r := &UrlPrintReq{}
	err := d.Decode(r)
	if err != nil {
		http.Error(w, "Invalid Params", 400)
	}

	dreq, err := http.NewRequest(http.MethodGet, r.URL, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if r.ForwardCookie != "" {
		cookies := strings.Split(strings.Replace(r.ForwardCookie, " ", "", -1), ";")
		for _, c := range cookies {
			name, value, _ := strings.Cut(c, "=")
			dreq.AddCookie(&http.Cookie{Name: name, Value: value})
		}
	}

	dresp, err := httpcln.Do(dreq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("File Download Error: %s", err.Error()), 400)
		return
	}
	defer dresp.Body.Close()

	size := dresp.ContentLength

	if size < 0 {
		size = int64(0)
	}

	doc := ipp.Document{
		Document: dresp.Body,
		Size:     int(size),
		MimeType: "application/pdf",
	}

	attrs := make(map[string]interface{})
	mediaCol := make(map[string]interface{})
	if r.Preset == "shipping-labels" {
		mediaCol["media-source"] = "tray-1"
		mediaCol["media-type"] = "labels"
		attrs[ipp.AttributeMedia] = "iso_a5_148x210mm"
	} else {
		mediaCol["media-source"] = "tray-2"
		mediaCol["media-type"] = "unspecified"
		attrs[ipp.AttributeMedia] = "iso_a4_210x297mm"
	}
	attrs[ipp.AttributeMediaCol] = mediaCol

	_, err = printcln.PrintJob(doc, printername, attrs)
	if err != nil {
		fmt.Println("PrintError:", err)
		http.Error(w, fmt.Sprintf("Printing error: %s", err.Error()), 400)
		return
	}

	fmt.Fprintln(w, "Print Successful!")
}

func main() {
	httpcln = &http.Client{}

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	printcln = ipp.NewIPPClient(os.Getenv("IPP_HOST"), 631, os.Getenv("IPP_USER"), os.Getenv("IPP_PASSWORD"), false)
	printername = os.Getenv("IPP_PRINTER_NAME")

	err = printcln.TestConnection()
	if err != nil {
		panic(err)
	}

	// attrs, _ := printcln.GetPrinterAttributes(printer, []string{"media-type-supported"})
	// fmt.Println(attrs)

	http.HandleFunc("/print", urlPrint)
	//http.HandleFunc("/upload", multipartPrint)
	http.ListenAndServe(":8080", nil)
}
