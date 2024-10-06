# Anura-POC

Anura: a build artifact repository written in Go.

## What?

After working a lot with Nexus and JFrog Artifactory, I wanted to build an artifact repository using Go.
Because, as a Devops, _I'm tired of deploying Java programs._

I started by cloning the Gitea project and removing everything that was not necessary for the artifact repository.
So I deleted all the files regarding git, repositories, issues, pull requests, SSH keys, and so on.

Right now it builds, starts, and shows most of the admin interfaces.

I will try to add all the repository APIs and UI.
