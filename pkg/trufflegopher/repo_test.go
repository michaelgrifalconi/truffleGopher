package trufflegopher

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestRepoScan(t *testing.T) {
	tempDir := "../../test/tmp"

	os.RemoveAll(tempDir)
	os.Mkdir(tempDir, os.ModePerm)
	_, err := unzip("../../test/dummy-repo.zip", tempDir)
	if err != nil {
		t.Fatal("Failed to unzip file")
	}
	s := Signature{Match: regexp.MustCompile("FAKE[0-9A-Z]{16}"), Description: "My fake KEY"}
	th := Trufflegopher{}
	sig := []Signature{}

	sig = append(sig, s)
	err = th.Init(1, sig)
	if err != nil {
		t.Fatal("Failed to initialize trufflegopher")
	}

	err = repoScan(tempDir, &th)
	if err != nil {
		t.Error("Failed to scan test repo")
	}
	f := <-th.FindingsChan
	if strings.Compare(f.Reason, "My fake KEY") != 0 {
		t.Fatal("Failed to find FAKE issue on commit")
	}

	os.RemoveAll(tempDir)
}
