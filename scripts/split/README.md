# split

These scripts split the monorepo into the public modules under
github.com/alexisbcz: libhtml, cn, cva, twmerge, tailwind, icons and
components.

1. `python3 stage.py DIR` copies each package to `DIR/split/<repo>` and
   rewrites the `shadcn.go/...` imports.
2. `python3 docs.py DIR` writes each README and NOTICE.
3. Write each `go.mod`, then lint the prose: `vale --config vale.ini DIR/split/*/README.md`.
4. `SPLIT_DIR=DIR ./publish.sh <repo>` creates the repository and pushes its
   initial commit. Publish libhtml, cn, cva, twmerge and tailwind first, then
   icons, then components, and set the `require` lines of each from the
   commits already published.

## Sync

`python3 sync.py DIR` stages every module again and reports which published
repositories would change. `python3 sync.py DIR --push --message MSG` commits
and pushes them, in dependency order. Each repository keeps its `go.mod` and
`go.sum`, takes the newest commits of the modules it requires, and must
build, vet and pass its tests on its own.
