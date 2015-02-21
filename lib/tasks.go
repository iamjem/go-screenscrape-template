package screenscrape

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/url"
)

type TaskResult struct {
	Source *url.URL
	Result *url.URL
	Error  error
}

type Task interface {
	Run(source *url.URL)
	Result() *TaskResult
}

// Acme Comics
type AcmeComicTask struct {
	source   *url.URL
	result   *url.URL
	runError error
}

func (t *AcmeComicTask) Run(source *url.URL) {
	t.source = source

	doc, err := goquery.NewDocument(source.String())
	if err != nil {
		t.runError = err
		return
	}

	latest, exists := doc.Find(".chapter-list ul li a").First().Attr("href")
	if exists {
		latestUrl, err := url.ParseRequestURI(latest)
		if err != nil {
			t.runError = err
			return
		}
		t.result = latestUrl
	} else {
		t.runError = errors.New("Unable to find element that matches selector.")
	}
}

func (t *AcmeComicTask) Result() *TaskResult {
	return &TaskResult{
		Source: t.source,
		Result: t.result,
		Error:  t.runError,
	}
}
