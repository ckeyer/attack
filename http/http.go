package http

import (
	"runtime"
	"sync"
	"time"

	"github.com/ckeyer/api/types"
	"github.com/ckeyer/logrus"
)

type Manager struct {
	types.HTTPOption
	sync.Mutex

	jobs  map[int]*Job
	stopC chan struct{}
}

// NewManager new a http attack manager
func NewManager(opt types.HTTPOption) *Manager {
	mgr := &Manager{
		HTTPOption: opt,
		jobs:       map[int]*Job{},
		stopC:      make(chan struct{}),
	}

	return mgr
}

// Execute do http attack
func Execute(opt types.HTTPOption) error {
	if err := CheckOption(opt); err != nil {
		logrus.Fatalln(err)
	}

	mgr := NewManager(opt)
	mgr.setupEnv()

	for i := 0; i < int(opt.Goroutine); i++ {
		mgr.runAJob(i+1, opt)
	}
	logrus.Infof("new %v jobs, %s %s", len(mgr.jobs), opt.Method, opt.Url)

	var wg sync.WaitGroup
	wg.Add(len(mgr.jobs))
	for _, job := range mgr.jobs {
		go func(j *Job) {
			defer wg.Done()
			j.Run()
		}(job)
	}
	wg.Wait()
	logrus.Info("over")
	time.Sleep(time.Second)
	return nil
}

// setupEnv set goroutine nums
func (mgr *Manager) setupEnv() {
	num := int(mgr.Goroutine)
	if num < runtime.NumCPU() {
		runtime.GOMAXPROCS(num)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	}
}

// runAJob start a new job
func (mgr *Manager) runAJob(index int, opt types.HTTPOption) {
	j := NewJob(index, opt, mgr.stopC)

	mgr.Lock()
	defer mgr.Unlock()
	if _, ok := mgr.jobs[index]; !ok {
		mgr.jobs[index] = j
	} else {
		logrus.Panicf("exists index %s", index)
	}
}
