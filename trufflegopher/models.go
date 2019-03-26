package trufflegopher

import (
	"regexp"
	"time"
)

//Signature struct contains regexp to run agains every commit, more signatures = more computation
type Signature struct {
	Match       *regexp.Regexp
	Description string `yaml:"description"`
	Exp         string `yaml:"exp"`
}

//Finding struct will contain details about the signature match on the scanned diff
type Finding struct {
	Diff          string
	CommitHash    string
	CommitMessage string
	CommitDate    time.Time
	Branch        string
	Reason        string
}
