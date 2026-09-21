---
name: release
description: "Cut a new claude-launcher release by bumping the version and pushing a tag. Use when the user asks to release, cut a release, ship, or publish a new version of claude-launcher."
---

# Release

## Step 0 — Sync main

```bash
git branch --show-current
git status
```

- Must be on `main` with a clean working tree. If not, confirm with the user how to proceed.

```bash
git pull origin main
```

## Step 1 — Propose a bump level

Get the latest tag:

```bash
git describe --tags --abbrev=0
```

Take the tag printed above (e.g. `v0.2.1`) and list commits since it, substituting the
literal value in place of `<latest_tag>`:

```bash
git log <latest_tag>..HEAD --oneline
```

Based on Conventional Commit types in that log, propose a level per
[Semantic Versioning](https://semver.org/):

| Commits since last tag | Suggested level |
| --- | --- |
| Any breaking change | `major` |
| `feat` (no breaking change) | `minor` |
| Only `fix` / `chore` / `docs` / etc. | `patch` |

Present the proposal with the reasoning to the user and get confirmation (or their own choice)
before proceeding.

## Step 2 — Run the release task

```bash
mise run release -- <patch|minor|major> -y
```

`mise/tasks/release` will:

1. Compute the next version from the latest tag (the version is injected via ldflags, so there is no version file to bump)
2. Create annotated tag `v<X.Y.Z>`
3. Push the tag to `origin`

`-y` skips the interactive confirmation prompt — the confirmation with the user already
happened in Step 1, so no need for a second, terminal-only confirmation here.

## Step 3 — Wait for CI to finish

The pushed tag triggers `.github/workflows/release.yml`, which runs tests, builds binaries for
each platform, generates a checksum, and creates the GitHub Release.

Get the run's database ID:

```bash
gh run list --workflow=release.yml --limit 1 --json databaseId,status,url
```

Take the `databaseId` printed above (e.g. `30054336679`) and watch it until it finishes,
substituting the literal value in place of `<run_id>`:

```bash
gh run watch <run_id> --exit-status
```

If the command exits non-zero, the run failed — inspect the failing job
(`gh run view <run_id> --log-failed`) and resolve it before telling the user the release is
done. Once it succeeds, report the run URL and confirm the release is complete.
