package trufflegopher

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func TestEnqueueJob(t *testing.T) {
	queue := make(chan bool, 2)

	dataChan := make(chan string, 20)
	youCanRetourn := make(chan bool)

	f := func() {
		dataChan <- "foo"
		<-youCanRetourn
	}

	for i := 1; i <= 3; i++ {
		go enqueueJob(queue, f)
	}
	time.Sleep(time.Millisecond * 600)

	if len(dataChan) != 2 {
		t.Error("Number of workers are not limited as expected")
	}

	for i := 1; i <= 3; i++ {
		youCanRetourn <- true
	}
	for len(queue) != 0 {
		//Waiting for all background jobs to finish..
		time.Sleep(time.Second)
	}
}

func TestInit(t *testing.T) {
	s := Signature{}
	s.Description = "Desc"
	s.Exp = "AKIA[0-9A-Z]{16}"

	th := Trufflegopher{}

	var sig []Signature

	sig = append(sig, s)
	err := th.Init(2, sig)
	if err != nil {
		t.Error("No errors were expected")
	}
	if th.workers != 2 {
		t.Error("Wrong number of workers")
	}

	if th.signatures[0].Exp != "AKIA[0-9A-Z]{16}" {
		t.Error("Wrong signature")
	}

	if !th.signatures[0].Match.Match([]byte("AKIAIOSFODNN7EXAMPLE")) {
		t.Error("Error parsing regular expressions")
	}
	if !th.SearchedDiffs.IsEmpty() {
		t.Error("Something went wrong with map initialization")
	}

	th = Trufflegopher{}
	sig = []Signature{}
	err = th.Init(10, sig)
	if err == nil {
		t.Error("Init should fail when no signatures are passed")
	}

}
