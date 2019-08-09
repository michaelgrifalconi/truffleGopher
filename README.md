

# truffleGopher
Searches through a given git repository for secrets, digging deep into each commit ever pushed on its target. It is meant to find accidentaly committed secrets, but you can use it with any given regular expression you like.

Credits goes to the original project: [truffleHog](https://github.com/dxa4481/truffleHog) which I used as baseline to design this Go implementation.

Due to design differences, do not expect any compatibility with the original TruffleHog even though the main goal is the same.

 - [x] written golang
 - [x] can be used as package in your project or as standalone tool in docker image
 - [x] configurable regex in YAML file
 - [x] considerable performance improvement compared to original TruffleHog
 - [x] able to store already scanned commits in a file, to make subsequent runs scan only new commits
   - [x] scans only additions and not whole diff
   - [x] implement parallel scan of multiple diffs to reduce chance to get stuck on single huge diff
   - [ ] provide easy way to benchmark performances between TruffleGopher and TruffleHog
 - [ ] allow SQL DB as destination of findings instead of just print(standalone) or channel(package)
 - [ ] provides metrics on how many commits/diff were scanned/skipped, number of findings, etc?
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
- clone git repos (See [GitPD](https://github.com/michaelgrifalconi/gitPD))
- scan multiple repos/orgs (See [GitPD](https://github.com/michaelgrifalconi/gitPD))
- entropy scans (might reconsider)

#### Usage

```
Not yet ready
```

#### Benchmark

```
Not yet ready
```

## Contributing
Feel free to open an issue to start the discussion.
