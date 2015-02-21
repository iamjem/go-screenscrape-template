package main

import (
	"errors"
	"flag"
	"fmt"
	screenscrape "github.com/iamjem/go-screenscrape-template/lib"
	"net/url"
)

type UrlFlag []*url.URL

func (flag *UrlFlag) String() string {
	return fmt.Sprintf("%v", *flag)
}

func (flag *UrlFlag) Set(val string) error {
	urlVal, err := url.ParseRequestURI(val)
	if err != nil {
		return err
	}

	if !urlVal.IsAbs() {
		return errors.New(fmt.Sprintf("Invalid URL '%s' - all values must be absolute.", urlVal))
	}

	*flag = append(*flag, urlVal)
	return nil
}

func (flag *UrlFlag) Get() interface{} {
	return []*url.URL(*flag)
}

var urlFlag UrlFlag

func init() {
	// Initialize urlFlag variable, and define url flag
	urlFlag = make(UrlFlag, 0)
	flag.Var(&urlFlag, "url", "Absolute URL to known comic source.")
}

func main() {
	//  Parse command line flags, and call screenscrape.Run() if one or more URLs
	flag.Parse()
	if len(urlFlag) > 0 {
		screenscrape.Run(urlFlag...)
	}
}
