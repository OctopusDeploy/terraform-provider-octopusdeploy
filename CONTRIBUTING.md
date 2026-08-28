# Contributing Guide

Thanks for contributing to this project! :+1: This project and everyone participating in it is governed by the [Octopus Deploy Code of Conduct](https://github.com/OctopusDeploy/.github/blob/main/CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior using the instructions in the code of conduct.

This guide provides an overview of the contribution workflow from opening an issue, creating a PR, reviewing, and merging the PR.

## Getting Started

This project is built, tested, and released by workflows defined in GitHub Actions (see [Actions](/actions/) for more information). Release management is controlled through [Release-Please](https://github.com/googleapis/release-please). 

Please use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary), and choose the "Squash and Merge" option to merge your PR (so we get a nice, clear Git history that makes it obvious what's changed)

### Issues

We :heart: feedback! Submitting an issue (i.e. feature, bug) is the best way to document things your experience with this project. For example, if there's a feature missing or there's behavior that doesn't match your expectations then we strongly encourage you to submit an issue. That way, contributors can track them and have interested folks (like you) by notified if/when they're resolved.

#### Create a New Issue

Use the Issues feature in GitHub to document bugs and/or features related to this project. Please ensure to apply any/all associated metadata (such as labels) in order to classify them appropriately. Also, please provide as much contextual information as you can, especially when documenting bugs. Templates are provided in this project to guide the authoring process.

#### Resolve an Issue

Issues will be triaged and modified (if necessary) by the [CODEOWNERS](CODEOWNERS) for this project. It is important to associate pull requests with issues by referencing their issue ID in the commit message. That way, issues will be able to document changes and/or fixes. This will assist visitors when reading through issue lists.

### Commit Your Change(s) through Pull Requests

This project employs [branch protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule); the `main` branch is protected. Therefore, your changes MUST be committed to a branch and submitted as a pull request. Also, this project requires the use of [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for all commit messages. Using Conventional Commits enables this project to autogenerate its [CHANGELOG.md](CHANGELOG.md) and release notes.

### Your Pull Request is Merged! Now What?

Congratulations! :tada: And thank you very much for your contribution to this project!

Once your pull request is merged, our build and test workflow will execute once again to validate changes. Afterward, your changes will be committed to the `main` branch.

## Merging Pull Requests from External Contributors

Our PR checks don't run automatically on pull requests from forks. You can approve the run from the PR, but it won't help: the checks require access to values that GitHub Actions won't expose to a workflow running in the contributor's profile rather than the Octopus organization, so they will always fail and the change can't be validated where it is. To work around this, redirect the change through a branch in this repository so the checks can run before it reaches `main`.

> [!WARNING]
> Redirecting a fork PR into a branch in this repository causes the contributor's code to run in workflows that hold repository secrets, including the GitHub App private key (`GH_APP_PRIVATE_KEY`), the App ID, and the Octopus licence. `build.yml` and `docs.yml` run on `push` to any branch, and the test workflow runs on `pull_request`, so any code merged into an in-repo branch executes with those secrets and a `contents: write` token. Treat this as running untrusted code with full access to those secrets.
>
> Before you redirect anything, read the entire diff line by line. Pay particular attention to workflow files under `.github/`, `Makefile`, `go.mod`/`go.sum`, `github-app-jwt.py`, and anything that runs during build or test. If the change modifies how CI runs, or you are not certain it is safe, do not redirect it.

> [!IMPORTANT]
> Don't merge a fork PR directly into `main`, and never redirect one you have not fully reviewed. Follow the steps below only after completing the review described above.

1. Review the contributor's diff in full and confirm it is safe to execute with repository secrets (see the warning above).
2. Create a new branch in this repository with a name similar to the contributor's branch (for example, if their branch is `fix-typo-in-docs`, create `fix-typo-in-docs`).
3. Edit the contributor's PR and change its base branch from `main` to the new branch you just created.
4. Squash and merge the PR into the new branch.
5. Continue as you would for any other change: open a PR from the intermediate branch to `main` and merge it.

A more robust alternative to reviewing by eye on every contribution is to move the secret-bearing jobs behind a [protected GitHub Environment](https://docs.github.com/en/actions/deployment/targeting-different-environments/using-environments-for-deployment#required-reviewers) that requires manual approval before the job with secrets runs. That keeps the secrets from being exposed to a run until a maintainer has explicitly approved it.
