package trufflegopher

import (
	"errors"
	"regexp"

	cmap "github.com/orcaman/concurrent-map"
)

//TODO: remove
var diffCounter int
var commitCounter int
var skippedCounter int

//Trufflegopher is the core of the project, instanciate only once and call multiple time its ScanRepo function
//signatures and searchedDiff database will be shared between different parallel scans
type Trufflegopher struct {
	SearchedDiffs  cmap.ConcurrentMap
	signatures     []Signature
	FindingsChan   chan Finding
	workers        int
	executionQueue chan bool
}

//TODO: review queue system, atm is more like a 'parallel workers allowed' and not a queue, it blocks
// if workers are busy and queue size is 0
func enqueueJob(executionQueue chan bool, item func()) {
	executionQueue <- true
	go func() {
		item()
		<-executionQueue
	}()
}

//Init will setup Trufflegopher with necessary configuration.
func (t *Trufflegopher) Init(workers int, signatures []Signature) error {
	if len(signatures) < 1 {
		return errors.New("Not enough signatures")
	}
	if signatures[0].Match == nil {
		//This means we got a file of strings and not already compiled regexp, so we compile them now
		for i := range signatures {
			signatures[i].Match = regexp.MustCompile(signatures[i].Exp)
		}
	}
	t.workers = workers
	//TODO: why 1 or 20 workers take same CPU% and same time overall??
	// seems like the scan is quite fast, does not make it to pile up
	t.executionQueue = make(chan bool, t.workers)
	t.FindingsChan = make(chan Finding, 200)
	t.SearchedDiffs = cmap.New()
	t.signatures = signatures
	return nil
}

//SaveState saves the current state of already scanned commits to a file
func (t *Trufflegopher) SaveState(path string) error {
	return diskOperationSave(path, t)
}

//LoadState loads a previously saved state of already scanned commits from a file
func (t *Trufflegopher) LoadState(path string) error {
	return diskOperationLoad(path, t)
}

//ScanRepo given a git repo location iterate on all commits
//TODO: provide some metrics (commits, line scanned, findings, repos, branches..)
func (t *Trufflegopher) ScanRepo(path string) error {
	return repoScan(path, t)
}
