package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/michaelgrifalconi/trufflegopher/pkg/trufflegopher"
	"gopkg.in/yaml.v2"
)

//SignatureFileStruct holds the array of signatures to scan given repo
type SignatureFileStruct struct {
	Signatures []trufflegopher.Signature `yaml:"signatures"`
}

//TODO: Consider moving most of this code to /internal dir according to https://github.com/golang-standards/project-layout
func main() {

	signaturesFile := flag.String("signatures", "", "regex signatures yaml file")
	targetRepo := flag.String("repo", "", "repo target path")
	scanDBFile := flag.String("dbfile", "scanDB", "filepath to use for storing already scanned diff database")
	flag.Parse()

	if *signaturesFile == "" {
		log.Println("Please provide a valid path to signature files `-signatures=file.yml`")
		os.Exit(1)
	}
	if *targetRepo == "" {
		log.Println("Please provide a valid path to a target repository")
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(*signaturesFile)
	if err != nil {
		log.Fatal(err)
	}

	s := SignatureFileStruct{}
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	t := new(trufflegopher.Trufflegopher)
	err = t.Init(10, s.Signatures)

	err = t.LoadState(*scanDBFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	go func() {
		for f := range t.FindingsChan {
			fmt.Println("==============================================")
			fmt.Println(f.Reason)
			fmt.Println(f.CommitDate)
			fmt.Println(f.CommitHash)
			fmt.Println(f.CommitMessage)
			fmt.Println(f.Diff)
		}
	}()
	err = t.ScanRepo(*targetRepo)
	if err != nil {
		log.Fatal(err)
	}

	err = t.SaveState(*scanDBFile)
	if err != nil {
		log.Fatal(err)
	}
}
