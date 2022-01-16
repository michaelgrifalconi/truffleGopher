

# truffleGopher
Searches through a given git repository for secrets, digging deep into each commit ever pushed on its target. It is meant to find accidentaly committed secrets, but you can use it with any given regular expression you like.

Credits goes to the original project: [truffleHog](https://github.com/dxa4481/truffleHog) which I used as baseline to design this Go implementation.

Due to design differences, do not expect any compatibility with the original TruffleHog even though the main goal is the same.

# Why?
Because it's way faster.

It is already faster on a first run, but the biggest advantage comes from the cached result set.
A common use case of such tools is to check very often the repository for credential leaks. This tool will remember what commit were checked and not scan them again.
Be careful! Right now it only remembers the visited commits. If you change the regexp used for the scan, you must delete the commit list since it must be scanned again.

Might be nice to remember also what regexp was evaluated for each commit to make it even more efficient.


### Out of scope
- clone git repos (See [GitPD](https://github.com/michaelgrifalconi/gitPD))
- scan multiple repos/orgs (See [GitPD](https://github.com/michaelgrifalconi/gitPD))
- entropy scans (might reconsider)

#### Usage

```
git clone YOUR_TARGET_REPO_HERE tmp/target-repo
docker run --rm -v "$(pwd)":/tmp michaelgrifalconi/tg:v1 -signatures="/tmp/tmp/trufflegopher-rules.yml" -repo="/tmp/tmp/target-repo" -dbfile="/tmp/tmp/scanDB"
```

#### Benchmark

```
./scripts/build-image
./scripts/benchmark

# Sample result on my laptop with github.com/golang/tools.git

# First run: 
truffleHog:    47s
truffleGopher: 33s

# Second run: (with cached result set)
truffleHog:    46s
truffleGopher: 7s
```

## Development status

 - [x] written golang
 - [x] can be used as package in your project or as standalone tool in docker image
 - [x] configurable regex in YAML file
 - [x] considerable performance improvement compared to original TruffleHog
 - [x] able to store already scanned commits in a file, to make subsequent runs scan only new commits
   - [x] scans only additions and not whole diff
   - [x] implement parallel scan of multiple diffs to reduce chance to get stuck on single huge diff
   - [x] provide easy way to benchmark performances between TruffleGopher and TruffleHog
 - [ ] allow SQL DB / CSV as destination of findings instead of just print(standalone) or channel(package)
 - [ ] provides metrics on how many commits/diff were scanned/skipped, number of findings, etc?
 - [ ] able to skip blacklisted paths in diffs (e.g. vendor dir)
 - [ ] follow [go project layout](https://github.com/golang-standards/project-layout)
   - [x] trufflegopher pkg in pkg dir
   - [x] cmd/app/main.go and internal/app/app.go for current main.go file
 - [ ] https://godoc.org/-/about
 - [ ] better error handling
 - [ ] better logging
 - [ ] better testing
 - [ ] CI/CD jobs 
 - [ ] cache with commit ID also the regex rule applied, to allow to add new rules and still use the cached commit list

## Contributing
Feel free to open an issue to start the discussion.
