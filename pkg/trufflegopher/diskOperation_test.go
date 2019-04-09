package trufflegopher

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestDiskOperationLoad(t *testing.T) {

	s := Signature{Match: regexp.MustCompile("FAKE[0-9A-Z]{16}"), Description: "My fake KEY"}
	th := Trufflegopher{}
	sig := []Signature{}

	sig = append(sig, s)
	err := th.Init(1, sig)
	if err != nil {
		t.Fatal("Failed to initialize trufflegopher")
	}

	err = th.LoadState("../../test/dummy-repoDB.bin")
	if err != nil {
		t.Fatal("Failed to load state from file")
	}
	_, ok := th.SearchedDiffs.Get("5c611b9765a9be73c9ff502aec0221d3f990571044cf358fc4d5653d7c885fb8af2bea0437f9ab21")
	if !ok {
		t.Fatal("Failed to retrieve expected key from loaded map")
	}
}
func TestDiskOperationSave(t *testing.T) {
	tempDir := "../../test/tmp"

	os.RemoveAll(tempDir)
	os.Mkdir(tempDir, os.ModePerm)

	s := Signature{Match: regexp.MustCompile("FAKE[0-9A-Z]{16}"), Description: "My fake KEY"}
	th := Trufflegopher{}
	sig := []Signature{}

	sig = append(sig, s)
	err := th.Init(1, sig)
	if err != nil {
		t.Fatal("Failed to initialize trufflegopher")
	}

	th.SearchedDiffs.Set("FAKE_COMMIT_PAIR", struct{}{})
	err = th.SaveState(tempDir + "/testDB")
	if err != nil {
		t.Fatal("Failed to save state to file")
	}

	f, err := os.Open(tempDir + "/testDB")
	if err != nil {
		t.Fatal("Failed to open the created file")
	}
	defer f.Close()
	hasher := sha1.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	if strings.Compare(fmt.Sprintf("%x", hasher.Sum(nil)), "55327b5fda0500e94ef034bf2a7675f1d5ade477") != 0 {
		t.Fatal("Failed to validate created files")
	}
	os.RemoveAll(tempDir)
}
