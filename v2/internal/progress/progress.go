package progress

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/projectdiscovery/clistats"
	"github.com/projectdiscovery/gologger"
)

// Progress is a progress instance for showing program stats
type Progress struct {
	active       bool
	stats        clistats.StatisticsClient
	tickDuration time.Duration
}

// NewProgress creates and returns a new progress tracking object.
func NewProgress(active bool) *Progress {
	var tickDuration time.Duration
	if active {
		tickDuration = 5 * time.Second
	} else {
		tickDuration = -1
	}

	var progress Progress
	if active {
		stats, err := clistats.New()
		if err != nil {
			gologger.Warningf("Couldn't create progress engine: %s\n", err)
		}
		progress.active = active
		progress.stats = stats
		progress.tickDuration = tickDuration
	}

	return &progress
}

// Init initializes the progress display mechanism by setting counters, etc.
func (p *Progress) Init(hostCount int64, rulesCount int, requestCount int64) {
	if p.active {
		p.stats.AddStatic("templates", rulesCount)
		p.stats.AddStatic("hosts", hostCount)
		p.stats.AddStatic("startedAt", time.Now())
		p.stats.AddCounter("requests", uint64(0))
		p.stats.AddCounter("errors", uint64(0))
		p.stats.AddCounter("total", uint64(requestCount))
		if err := p.stats.Start(makePrintCallback(), p.tickDuration); err != nil {
			gologger.Warningf("Couldn't start statistics: %s\n", err)
		}
	}
}

// AddToTotal adds a value to the total request count
func (p *Progress) AddToTotal(delta int64) {
	if p.active {
		p.stats.IncrementCounter("total", int(delta))
	}
}

// Update progress tracking information and increments the request counter by one unit.
func (p *Progress) Update() {
	if p.active {
		p.stats.IncrementCounter("requests", 1)
	}
}

// Drop drops the specified number of requests from the progress bar total.
// This may be the case when uncompleted requests are encountered and shouldn't be part of the total count.
func (p *Progress) Drop(count int64) {
	if p.active {
		// mimic dropping by incrementing the completed requests
		p.stats.IncrementCounter("errors", int(count))
	}
}

const bufferSize = 128

func makePrintCallback() func(stats clistats.StatisticsClient) {
	builder := &strings.Builder{}
	builder.Grow(bufferSize)

	return func(stats clistats.StatisticsClient) {
		builder.WriteRune('[')
		startedAt, _ := stats.GetStatic("startedAt")
		duration := time.Since(startedAt.(time.Time))
		builder.WriteString(fmtDuration(duration))
		builder.WriteRune(']')

		templates, _ := stats.GetStatic("templates")
		builder.WriteString(" | Templates: ")
		builder.WriteString(clistats.String(templates))
		hosts, _ := stats.GetStatic("hosts")
		builder.WriteString(" | Hosts: ")
		builder.WriteString(clistats.String(hosts))

		requests, _ := stats.GetCounter("requests")
		total, _ := stats.GetCounter("total")

		builder.WriteString(" | RPS: ")
		builder.WriteString(clistats.String(uint64(float64(requests) / duration.Seconds())))

		errors, _ := stats.GetCounter("errors")
		builder.WriteString(" | Errors: ")
		builder.WriteString(clistats.String(errors))

		builder.WriteString(" | Requests: ")
		builder.WriteString(clistats.String(requests))
		builder.WriteRune('/')
		builder.WriteString(clistats.String(total))
		builder.WriteRune(' ')
		builder.WriteRune('(')
		//nolint:gomnd // this is not a magic number
		builder.WriteString(clistats.String(uint64(float64(requests) / float64(total) * 100.0)))
		builder.WriteRune('%')
		builder.WriteRune(')')
		builder.WriteRune('\n')

		fmt.Fprintf(os.Stderr, "%s", builder.String())
		builder.Reset()
	}
}

// fmtDuration formats the duration for the time elapsed
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

// Stop stops the progress bar execution
func (p *Progress) Stop() {
	if p.active {
		if err := p.stats.Stop(); err != nil {
			gologger.Warningf("Couldn't stop statistics: %s\n", err)
		}
	}
}
