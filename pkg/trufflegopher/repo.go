package trufflegopher

import (
	"fmt"
	"log"
	"time"

	git "github.com/libgit2/git2go/v33"
)

//TODO: review Cgo usage and check if we have to free some memory manually?

//More info at https://gobyexample.com/closures
func getOdbForEachCallback(r *git.Repository, t *Trufflegopher) git.OdbForEachCallback {
	return func(id *git.Oid) error {
		commit, err := r.LookupCommit(id) //this is an expensive call, but can we avoid it somehow?
		if err != nil {
			//object is not a commit, move on
			return nil
		}

		commitCounter++
		var commitID, parentID string
		var commitTree, parentTree *git.Tree

		commitID = commit.Id().String()
		commitTree, err = commit.Tree()
		if err != nil {
			log.Fatal(err)
		}

		parent := commit.Parent(0)

		//commit has a parent
		if parent != nil { //commit has a parent
			parentID = parent.Id().String()
			parentTree, err = parent.Tree()
			if err != nil {
				log.Fatal(err)
			}
		} else { //was first commit, parent will be an empty tree
			parentID = "0000000000000000000000000000000000000000"
			parentTree = &git.Tree{}
		}
		// try to add commit pair to already scanned DB, if was already there, do nothing and return
		if !t.SearchedDiffs.SetIfAbsent(parentID+commitID, struct{}{}) {
			skippedCounter++
			return nil
		}

		diffOpt, err := git.DefaultDiffOptions()
		if err != nil {
			log.Fatal(err)
		}

		d, err := r.DiffTreeToTree(parentTree, commitTree, &diffOpt)
		if err != nil {
			log.Println("failed to diff")
			log.Fatal(err)
		}
		numDeltas, err := d.NumDeltas()
		if err != nil {
			log.Println("failed to get numdeltas")
			log.Fatal(err)
		}
		for i := 0; i < numDeltas; i++ {
			patch, err := d.Patch(i)
			if err != nil {
				log.Println("failed to patch")
				log.Fatal(err)
			}
			diffCounter++
			patchString, err := patch.String()
			if err != nil {
				log.Println("failed to patch string")
				log.Fatal(err)
			}
			func(pcHash string, pcMessage string, pcDate time.Time, t *Trufflegopher) { //TODO: on failure, we must remove the commit pair from DB to avoid losing findings
				enqueueJob(t.executionQueue, func() {
					scan(patchString, pcHash, pcMessage, pcDate, t)
				})
			}(commit.Id().String(), commit.Message(), commit.Author().When, t)
		}
		return nil
	}
}

func repoScan(repoPath string, t *Trufflegopher) error {

	r, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Println("Failed to open repo " + repoPath)
	}
	odb, err := r.Odb()
	if err != nil {
		log.Fatal(err)
	}

	err = odb.ForEach(getOdbForEachCallback(r, t))
	if err != nil {
		log.Fatal(err)
	}

	for len(t.executionQueue) != 0 {
		fmt.Println("Waiting for all background jobs to finish..")
		time.Sleep(time.Second)
	}

	// fmt.Println("DIFF: ", diffCounter)
	// fmt.Println("COMMIT: ", commitCounter)
	// fmt.Println("SKIP: ", skippedCounter)
	return nil
}
