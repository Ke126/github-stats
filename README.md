# GitHub Stats

![card](card.svg)

GitHub Stats is a project which automatically generates a simple card displaying some useful stats from your GitHub account in Scalable Vector Graphics (SVG) format. This card can be viewed in any browser or linked in any of your repositories' READMEs.

## Overview

This repository contains a GitHub API client written in Go. This client uses your GitHub Personal Access Token (PAT) to retrieve information from your GitHub account, which is then used to generate an SVG using Go templates. The generated SVG is placed directly in this repository, so that it can be linked to from external sources.

This process is automated via a daily cron job run using GitHub Actions. The [update-card.yaml](.github/workflows/update-card.yaml) workflow will automatically generate a new card each day with fresh information from GitHub's API, and commit the new card to the repository if needed.

## Setup

If you are interested in using this project to get your own stats card SVG, here are the steps needed to set it up:

1. Fork this repository.
2. Create a [GitHub fine-grained Personal Access Token (PAT)](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) with the following settings:
    - Set **Token name** to anything you wish.
    - Set **Expiration** to **No expiration** (alternatively, you may create tokens with a shorter lifetime, although you will need to manually create new ones when they expire).
    - Under **Repository access**, select **All repositories**.
    - Under **Repository permissions**, set **Metadata** access to **Read-only**.
3. Save this PAT as a repository secret in your forked repository named `GH_PAT`.
4. Run the [update-card.yaml](.github/workflows/update-card.yaml) workflow, either manually on the Actions tab or automatically via the daily cron job.

Your GitHub stats card SVG should now be visible in the repository! You can link to it using a url of the form https://raw.githubusercontent.com/{user}/{repo}/refs/heads/main/card.svg, substituting in your GitHub username and repository name as necessary.
