# jarvis
SRE toolbox

# Requirements

- Git
- gcc (build-essential on debian based); used for [Makefile](Makefile)
- Go, for installation [instructions](https://golang.org/doc/install)

# Architecture

Intelligence are located under the `internal/pkg` folder. Frontends are under the `cmd` folder.
Usage of go modules.

# How to release?

1. You need a personnal token [instructions](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)
2. To release, run the following command `GITHUB_TOKEN=XXXXX make release`
you could also store your personnal token globally, thanks to git `git config --global github.token "....."`
