#!/usr/bin/env bash
# Publishes the docs site to GitHub Pages (https://alexisbcz.github.io/shadcn.go):
# exports it as static files with cmd/docs -export and force-pushes them as
# the single commit of the gh-pages branch, which Pages serves.
#
# It runs locally, not in CI: the site reads the upstream checkout in ui/
# (gitignored), which CI does not have.
#
#	scripts/deploy-pages.sh
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
repo=https://github.com/alexisbcz/shadcn.go.git
base_path=/shadcn.go
site_url=https://alexisbcz.github.io/shadcn.go

out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT

cd "$root"
go run ./cmd/docs -export "$out/site" -base-path "$base_path" -site-url "$site_url"

rev=$(git rev-parse --short HEAD)
cd "$out/site"
git init -q -b gh-pages
git add -A
git -c user.name="Alexis Bouchez" -c user.email=276929831+alexisbcz@users.noreply.github.com \
	commit -q -m "Docs site at $rev"
git push -f "$repo" gh-pages
