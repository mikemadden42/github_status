# TODO

## Bugs / correctness

- [x] Fix `.golangci.yml`. It fails `golangci-lint config verify` (v2), so the
      linter selection is ignored (staticcheck runs despite "disable-all"):
  - [x] `version: 2` should be `version: "2"`
  - [x] `linters.disable-all` should be `linters.default: none`
  - [x] `linters-settings` should be `linters.settings`
  - [x] goimports settings go under `formatters.settings.goimports.local-prefixes` (list)
  - [x] `issues.exclude-dirs` should be `linters.exclusions.paths`
- [ ] Fix `.goreleaser.yaml`. It fails `goreleaser check`:
  - [ ] Add `version: 2`
  - [ ] `snapshot.name_template` should be `snapshot.version_template`
  - [ ] Remove the no-op `go generate ./...` hook
- [x] HTTP handling in `github_status.go`:
  - [x] Use an `http.Client` with a timeout instead of `http.Get`
  - [x] Close `response.Body`
  - [x] Check `response.StatusCode` before decoding
- [x] Replace `panic` with an error message on stderr and `os.Exit(1)` (or `log.Fatal`)

## Code quality

- [x] Trim the `summary` struct to the fields actually used; `interface{}` should be `any`
- [x] Print the overall status (`Status.Description`)
- [x] Replace the if/else chain with a `switch` (staticcheck QF1003) covering
      `operational`, `degraded_performance`, `partial_outage`, `major_outage`,
      `under_maintenance`
- [x] Filter the "Visit www.githubstatus.com for more information" filler row
- [x] Use `json.NewDecoder(resp.Body).Decode` instead of `io.ReadAll` + `Unmarshal`
- [x] Use the readable URL `https://www.githubstatus.com/api/v2/summary.json`
- [x] Add a package comment (revive)
- [ ] Exit non-zero when any component isn't operational

## Tests / CI

- [x] Split into `fetch(url)` and `render(w io.Writer, summary)`; test `render`
      with a JSON fixture and `fetch` with `httptest.NewServer`
- [ ] Add a lint job (`golangci/golangci-lint-action`) to CI
- [x] CI targets Go `oldstable` and `stable` (via `actions/setup-go` aliases)
- [x] Raise `go.mod` to `go 1.26` to match the oldest Go version CI tests
- [ ] Bump `actions/setup-go` / `actions/checkout`
- [ ] Update dependencies (`go get -u ./... && go mod tidy`); `x/sys` v0.25.0 → v0.50.0,
      `fatih/color` v1.18.0 → v1.19.0

## Housekeeping

- [ ] Delete the Jenkinsfile unless Jenkins is in use (build-only, no tests, echo-only post blocks)
- [ ] `make release` uses `--snapshot`, so it never publishes. Rename the target or drop the flag
- [ ] Add `dist/` and `github_status` to `.gitignore`
- [ ] Expand the README with usage, install and example output
