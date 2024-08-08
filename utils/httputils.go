package utils

import (
	"github.com/antchfx/xmlquery"
	"io"
	"strings"
)

func GetXmlPath(res, path string) string {
	doc, err := xmlquery.Parse(strings.NewReader(res))
	if err != nil {
		panic(err)
	}
	result := ""
	if n := doc.SelectElement(path); n != nil {
		result = n.InnerText()
	}
	return result
}

func CloseBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		panic(err)
	}
}
