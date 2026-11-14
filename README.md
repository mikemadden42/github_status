# github_status

![example workflow](https://github.com/mikemadden42/github_status/actions/workflows/go.yml/badge.svg)

Print the current status of each GitHub component, from
[githubstatus.com](https://www.githubstatus.com).

## Install

```sh
go install github.com/mikemadden42/github_status@latest
```

## Usage

```sh
github_status
```

```
All Systems Operational

Git Operations                                              operational
Webhooks                                                    operational
API Requests                                                operational
Issues                                                      operational
Pull Requests                                               operational
Actions                                                     operational
Packages                                                    operational
Pages                                                       operational
Copilot                                                     operational
Codespaces                                                  operational
Copilot AI Model Providers                                  operational
```

Statuses are colored: green for operational, yellow for degraded performance or
partial outage, red for major outage and blue for maintenance.

### Scripting

With `-check`, `github_status` exits with status 1 if any component is not
operational:

```sh
github_status -check || echo "GitHub is having problems"
```

It also exits with status 1 if the status page can't be reached.

## Development

```sh
go test ./...
make check      # golangci-lint
make build      # goreleaser snapshot build into dist/
make snapshot   # full goreleaser snapshot release (local only, nothing is published)
```
