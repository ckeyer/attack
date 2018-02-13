package http

import (
	"net/http"

	"github.com/ckeyer/api/types"
	"github.com/ckeyer/commons/httpcli"
	"github.com/ckeyer/commons/httpcli/useragent"
	"github.com/ckeyer/logrus"
)

type Job struct {
	types.HTTPOption

	index int
	// run times
	num, successful int
	stopC           <-chan struct{}
	req             *http.Request
	cli             *httpcli.Client
}

func NewJob(index int, opt types.HTTPOption, ch <-chan struct{}) *Job {
	j := &Job{
		HTTPOption: opt,
		stopC:      ch,
		index:      index,
		cli:        httpcli.NewClient(),
	}
	if j.RandUA {
		j.req, _ = useragent.NewRequest(opt.Method, opt.Url, nil)
	} else {
		j.req, _ = http.NewRequest(opt.Method, opt.Url, nil)
	}
	for k, v := range opt.Headers {
		j.req.Header.Add(k, v)
	}
	return j
}

// Request
func (j *Job) Request() *http.Request {
	return j.req
}

// Inc ...
func (j *Job) Inc() {
	j.num++
}

// Next enable next
func (j *Job) Next() bool {
	return j.num < int(j.Count)
}

// Run start a job
func (j *Job) Run() {
	if j.Count <= 0 {
		j.Count = 99999
	}
	for j.Next() {
		select {
		case <-j.stopC:
			return
		default:
			j.Inc()
			j.do()
		}
	}
	return
}

// do do a http request.
func (j *Job) do() {
	req := j.Request()
	_, err := j.cli.Do(req)
	if err != nil {
		logrus.Warnf("job %v.(%v/%v) ,http reqeust failed, %s", j.index, j.num, j.Count, err)
		return
	}
	logrus.Debugf("job %v.(%v/%v) done, %v", j.index, j.num, j.Count, req.Header)
}
