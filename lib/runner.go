package screenscrape

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/url"
	"strings"
	"sync"
	"text/template"
)

var newChaptersEmailTemplate = template.Must(template.New("newChaptersTemplate").Parse(`
New comics have arrived!
{{ range .Chapters }}
{{ . }}
{{ end }}
`))

func Run(sources ...*url.URL) {
	var wg sync.WaitGroup
	var tasks = make([]Task, len(sources))

	// Get task by host for each URL
	for i, taskUrl := range sources {
		switch taskUrl.Host {
		case "www.acmewebcomics.com":
			tasks[i] = new(AcmeComicTask)
		default:
			log.Panicf("Unknown host %v", taskUrl.Host)
		}
	}

	wg.Add(len(tasks))

	// run each task
	for i, task := range tasks {
		go func(t Task, u *url.URL) {
			t.Run(u)
			wg.Done()
		}(task, sources[i])
	}

	// wait for all tasks to complete
	wg.Wait()

	// collect successful tasks
	results := map[string]*TaskResult{}
	resultSources := []string{}

	for _, task := range tasks {
		r := task.Result()
		if r.Error == nil {
			results[r.Source.String()] = r
			resultSources = append(resultSources, r.Source.String())
		} else {
			log.Warnf("Task for '%v' encountered err '%v'", r.Source, r.Error)
		}
	}

	// query for existing rows
	if len(resultSources) == 0 {
		log.Warn("No sources found")
		return
	}

	// hold all of the new chapters
	newLatest := []string{}

	var records []*Record
	if _, err := dbmap.Select(&records,
		fmt.Sprintf(`SELECT r.* FROM records r WHERE r.source IN ('%s') ORDER BY r.source ASC`, strings.Join(resultSources, "', '"))); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error querying database")
		return
	}

	// first check existing records
	for _, record := range records {
		if result, ok := results[record.Source]; ok {
			if result.Result.String() != record.Latest {
				// new latest URL, save URL and update record
				log.Infof("Found new chapter '%v'", result.Result)
				newLatest = append(newLatest, result.Result.String())
				record.Latest = result.Result.String()
				if _, err := dbmap.Update(record); err != nil {
					log.WithFields(log.Fields{
						"error":     err,
						"record_id": record.Id,
					}).Error("Error updating record")
				}
			}
		}
		delete(results, record.Source)
	}

	// anything left in the result map is totally new
	for source, result := range results {
		log.Infof("Found new chapter '%v'", result.Result)
		newLatest = append(newLatest, result.Result.String())
		if err := dbmap.Insert(&Record{
			Source: source,
			Latest: result.Result.String(),
		}); err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"source": source,
				"latest": result.Result,
			}).Error("Error inserting record")
		}
	}

	// if theres new chapters, send the email
	if len(newLatest) > 0 {
		var mailContent bytes.Buffer
		ctx := struct {
			Chapters []string
		}{
			newLatest,
		}

		log.Info("Sending new chapter email")
		if err := newChaptersEmailTemplate.Execute(&mailContent, ctx); err == nil {
			if err := SendMail(mailContent.String()); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error sending email")
			}
		}
	}
}
