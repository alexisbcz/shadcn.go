#!/bin/sh
# Builds and serves the upstream documentation site (ui/apps/v4) on
# http://localhost:4000: the pixel reference for the docs port (DOCS_PLAN.md).
# Requires node and bun; runs entirely offline once dependencies are installed.
set -e
cd "$(dirname "$0")/../../ui"
npx --yes pnpm@10.33.4 install --frozen-lockfile --ignore-scripts
cd apps/v4
npx --yes pnpm@10.33.4 run registry:build   # generates styles/<base>-<style> and the docs example index
NEXT_TELEMETRY_DISABLED=1 npx next build
exec env NEXT_TELEMETRY_DISABLED=1 npx next start --port 4000
