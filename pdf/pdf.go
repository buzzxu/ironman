package pdf

import (
	"errors"
	"github.com/buzzxu/boys/common/files"
	"github.com/buzzxu/boys/common/radom"
	"github.com/buzzxu/ironman/conf"
	"github.com/thecodingmachine/gotenberg-go-client/v7"
	"log"
	"net/http"
	"time"
)

type Oper int

const (
	OFFICE Oper = iota
	MARKDOWN
	HTML
	URL
	MERGE
)

type PDFFile struct {
	Url      string
	Body     string
	Data     *[]byte
	Tag      string
	FileName string
}
type PDF struct {
	Files     []*PDFFile
	Oper      Oper
	Ranges    string
	Landscape bool
	Scale     float64
	Margins   [4]float64
	Delay     float64
	Timeout   float64
	FileName  string
}

var client gotenberg.Client

func Client() {
	if conf.ServerConf.Pdf == nil {
		log.Fatal("PDF 客户端配置失败. pdf: \n\t url: http://127.0.0.1:3000")
	}
	timeOut := conf.ServerConf.Pdf.Timeout
	if timeOut == 0 {
		timeOut = 30
	}
	httpClient := &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
	}
	client = gotenberg.Client{Hostname: conf.ServerConf.Pdf.Url, HTTPClient: httpClient}
}

func To(file *PDF) (*http.Response, error) {
	req, err := Request(file)
	if err != nil {
		return nil, err
	}
	return client.Post(req)
}

func Save(file *PDF, dest string) error {
	req, err := Request(file)
	if err != nil {
		return err
	}
	return client.Store(req, dest)
}

func Request(file *PDF) (gotenberg.Request, error) {
	var files = make(map[string]gotenberg.Document, len(file.Files))
	for _, f := range file.Files {
		document, err := Document(f)
		if err != nil {
			return nil, err
		}
		var tag = f.Tag
		if tag == "" {
			tag = "index"
		}
		files[tag] = document
	}

	var fileName = file.FileName
	if fileName == "" {
		fileName = radom.Alphanumeric(12) + ".pdf"
	}
	switch file.Oper {
	case OFFICE:
		req := gotenberg.NewOfficeRequest(files["index"])
		req.Landscape(file.Landscape)
		if file.Timeout != 0 {
			req.WaitTimeout(file.Timeout)
		}
		if file.Ranges != "" {
			req.PageRanges(file.Ranges)
		}
		if file.Timeout > 0 {
			req.WaitTimeout(file.Timeout)
		}
		req.ResultFilename(fileName)
		return req, nil
	case MARKDOWN:
		req := gotenberg.NewMarkdownRequest(files["index"], files["markdown"])
		req.Landscape(file.Landscape)
		if file.Timeout != 0 {
			req.WaitTimeout(file.Timeout)
		}
		if file.Ranges != "" {
			req.PageRanges(file.Ranges)
		}
		if files["header"] != nil {
			req.Header(files["header"])
		}
		if files["footer"] != nil {
			req.Footer(files["footer"])
		}
		if files["assets"] != nil {
			req.Assets(files["assets"])
		}
		if len(file.Margins) > 0 {
			req.Margins(file.Margins)
		}
		if file.Scale > 0 {
			req.Scale(file.Scale)
		}
		if file.Delay != 0 {
			req.WaitDelay(file.Delay)
		}
		if file.Ranges != "" {
			req.PageRanges(file.Ranges)
		}
		if file.Timeout > 0 {
			req.WaitTimeout(file.Timeout)
		}
		req.ResultFilename(fileName)
		return req, nil
	case HTML:
		req := gotenberg.NewHTMLRequest(files["index"])
		if files["header"] != nil {
			req.Header(files["header"])
		}
		if files["footer"] != nil {
			req.Footer(files["footer"])
		}
		if files["assets"] != nil {
			req.Assets(files["assets"])
		}
		if len(file.Margins) > 0 {
			req.Margins(file.Margins)
		}
		if file.Scale > 0 {
			req.Scale(file.Scale)
		}
		if file.Delay != 0 {
			req.WaitDelay(file.Delay)
		}
		if file.Ranges != "" {
			req.PageRanges(file.Ranges)
		}
		if file.Timeout > 0 {
			req.WaitTimeout(file.Timeout)
		}

		req.ResultFilename(fileName)
		return req, nil
	case URL:
		req := gotenberg.NewURLRequest(file.Files[0].Url)
		if files["header"] != nil {
			req.Header(files["header"])
		}
		if files["footer"] != nil {
			req.Footer(files["footer"])
		}
		if len(file.Margins) > 0 {
			req.Margins(file.Margins)
		}
		if file.Scale > 0 {
			req.Scale(file.Scale)
		}
		if file.Delay != 0 {
			req.WaitDelay(file.Delay)
		}
		if file.Ranges != "" {
			req.PageRanges(file.Ranges)
		}
		if file.Timeout > 0 {
			req.WaitTimeout(file.Timeout)
		}
		req.ResultFilename(fileName)
		return req, nil
	case MERGE:
		documents := make([]gotenberg.Document, 0)
		for _, v := range files {
			documents = append(documents, v)
		}
		req := gotenberg.NewMergeRequest(documents...)
		req.WaitTimeout(file.Timeout)
		req.ResultFilename(fileName)
		return req, nil
	}
	return nil, errors.New("request fail")
}

func Document(file *PDFFile) (documet gotenberg.Document, err error) {
	var bytes *[]byte
	if file.Url != "" {
		bytes, err = files.URL(file.Url)
		if err != nil {
			return nil, err
		}
	} else if file.Data != nil {
		bytes = file.Data
	} else if file.Body != "" {
		return gotenberg.NewDocumentFromString(file.FileName, file.Body)
	} else {
		return nil, errors.New("无法获取数据")
	}
	documet, err = gotenberg.NewDocumentFromBytes(file.FileName, *bytes)
	return
}
