

# truffleGopher
Searches through a given git repository for secrets, digging deep into each commit ever pushed on its target. It is meant to find accidentaly committed secrets, but you can use it with any given regular expression you like.

Credits goes to the original project: [truffleHog](https://github.com/dxa4481/truffleHog) which I used as baseline to design this Go implementation.

Due to design differences, do not expect any compatibility with the original TruffleHog even though the end result is expected to be the same.

 - written golang
 - designed to be integrated as package, CLI is only a convenience tool (See GitPD project)
 - [x] configurable regex in YAML file
 - [x] considerable performance improvement compared to original TruffleHog
 - [x] able to store already scanned commits in a file, to make subsequent runs scan only new commits
   - [x] scans only additions and not whole diff
   - [x] implement parallel scan of multiple diffs to reduce chance to get stuck on single huge diff
   - [ ] provide easy way to benchmark performances between TruffleGopher and TruffleHog
 - [ ] send findings to given SQL database instead of sending through channel
 - [ ] provides metrics on how many commits/diff were scanned/skipped, number of findings, more?
 - [ ] able to skip blacklisted paths in diffs (e.g. vendor dir)
 - [ ] follow [go project layout](https://github.com/golang-standards/project-layout)
   - [ ] trufflegopher pkg in pkg dir
   - [ ] cmd/app/main.go and internal/app/app.go for current main.go file
 - [ ] https://godoc.org/-/about
 - [ ] better error handling
 - [ ] better logging
 - [ ] better testing
 - [ ] CI/CD jobs 


### Out of scope
- clone git repos (See GitPD)
- multile repos/orgs (See GitPD)
- entropy scans (might reconsider)


## CLI
This CLI is meant as rudumental way to run TruffleGopher on a single repo, but this project was tailored for GitPD, go have a look at it :-)

#### Install
TODO: Document lib2go installation
```
go get "github.com/michaelgrifalconi/trufflegopher"
```

#### Usage
```
trufflegopher -signatures="sample-signatures.yml" -repo="PATH_TO_TARGET_REPO" -dbfile="PATH_TO_DB_FILE"
```
```
-signatures: path to YAML file with signatures, see sample-signatures.yml as an example

-repo: path to (already cloned) target repository

-dbfile: path to existing database file or where to create the new file
```

#### Benchmark
Would you like to compare the time difference between TruffleHog and TruffleGopher?
First have a look at the benchmark.sh script and understand it, I take no responsibility if it blows up your PC.

Then run with no arguments to use a default repository as target
```
./benchmark.sh
```

Or pass your git repo as first argument, maybe you can use the Linux Kernel repo if you need to heat up your room.
```
./benchmark.sh "https://github.com/torvalds/linux.git"
```

## Contributing
Feel free to open an issue to start the discussion.
