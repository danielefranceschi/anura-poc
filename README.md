# Anura-POC

Anura: a build artifact repository written in Go.

## What?

After working a lot with Nexus and JFrog Artifactory, I wanted to build an artifact repository using Go.

Because, as a Devops, _I'm tired of deploying Java programs._
Expecially, tired of the fact that the JVM requires finetuning _after_ you finetune the container.

## What happened?

I started by cloning the Gitea project and removing everything that was not necessary for the artifact repository.
So I deleted all the files regarding git, repositories, issues, pull requests, SSH keys, and so on.

Right now it builds, starts, and shows most of the admin interfaces.

## (initial) roadmap

- [X] always build with sqlite (using https://gitlab.com/cznic/sqlite)
- [X] goreleaser builds
- [ ] add artifact repository model
- [ ] add artifact repository API
- [ ] dynamic routes for repos
- [ ] add artifact repository UI
- [ ] fix tests
- [ ] reduce devcontainer clutter

## Local build and run

```shell
goreleaser build --single-target --snapshot --clean
dist/anura_linux_arm64/anura -C $(pwd)/custom -w $(pwd)
```
