package trufflegopher

import (
	"bufio"
	"strings"
	"time"
)

func scan(diff string, commitHash string, commitMessage string, commitDate time.Time, t *Trufflegopher) {
	var f Finding
	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "+") { //Scanning only additions in the diff text
			for _, s := range t.signatures {
				if s.Match.MatchString(scanner.Text()) {
					f.Diff = scanner.Text()
					f.CommitHash = commitHash
					f.CommitMessage = commitMessage
					f.CommitDate = commitDate
					f.Reason = s.Description
					t.FindingsChan <- f
					break
				}
			}
		}
	}
}
