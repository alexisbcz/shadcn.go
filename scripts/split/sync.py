"""Syncs the published modules with the monorepo.

    python3 sync.py WORKDIR [--push] [--message MSG]

Re-stages every module from the monorepo (stage.py, docs.py) into
WORKDIR/next, then updates a clone of each published repository in
WORKDIR/repos: the staged files replace the clone's (go.mod and go.sum are
kept), the alexisbcz/* requirements move to the newest published commits in
dependency order, and the module must build, vet and pass its tests on its
own. Repositories whose content changed get one commit (MSG, by the
alexisbcz noreply address); --push pushes them. Without --push it only
reports what would change.
"""

import argparse
import os
import shutil
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
OWNER = "alexisbcz"
EMAIL = "276929831+alexisbcz@users.noreply.github.com"
NAME = "Alexis Bouchez"
# Dependency order: a module comes after the modules it requires.
ORDER = ["libhtml", "cn", "cva", "twmerge", "tailwind", "icons", "components"]
DEPS = {"icons": ["libhtml"], "components": ["libhtml", "cn", "cva", "icons"]}
GOENV = dict(os.environ, GOWORK="off", GOPROXY="direct", GONOSUMDB=f"github.com/{OWNER}/*", GOFLAGS="-mod=mod")


def run(args, cwd, env=None, check=True):
    r = subprocess.run(args, cwd=cwd, env=env, text=True, capture_output=True)
    if check and r.returncode != 0:
        sys.exit(f"{' '.join(args)} in {cwd} failed:\n{r.stdout}{r.stderr}")
    return r.stdout.strip()


def root_requirements():
    """The monorepo's external requirements: module -> version."""
    out = run(["go", "list", "-m", "-f", "{{.Path}} {{.Version}}", "all"], os.path.join(HERE, "../.."))
    reqs = {}
    for line in out.splitlines():
        parts = line.split()
        if len(parts) == 2 and not parts[0].startswith("shadcn.go"):
            reqs[parts[0]] = parts[1]
    return reqs


def pin(clone, reqs):
    """Moves the clone's direct external requirements to the monorepo's
    versions (the ones its tests ran with)."""
    out = run(["go", "list", "-m", "-f", "{{if not .Indirect}}{{.Path}}{{end}}", "all"], clone, GOENV)
    for mod in out.split():
        if mod in reqs and not mod.startswith(f"github.com/{OWNER}/"):
            run(["go", "mod", "edit", f"-require={mod}@{reqs[mod]}"], clone)
    run(["go", "mod", "tidy"], clone, GOENV)


def replace_tree(src, dst):
    """Makes dst hold src's files, keeping dst's .git, go.mod and go.sum."""
    keep = {".git", "go.mod", "go.sum"}
    for name in os.listdir(dst):
        if name in keep:
            continue
        p = os.path.join(dst, name)
        shutil.rmtree(p) if os.path.isdir(p) else os.remove(p)
    for name in os.listdir(src):
        if name in keep:
            continue
        s, d = os.path.join(src, name), os.path.join(dst, name)
        shutil.copytree(s, d) if os.path.isdir(s) else shutil.copy2(s, d)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("workdir")
    ap.add_argument("--push", action="store_true")
    ap.add_argument("--message", default="sync with the monorepo")
    a = ap.parse_args()
    work = os.path.abspath(a.workdir)
    nxt, repos = os.path.join(work, "next"), os.path.join(work, "repos")
    shutil.rmtree(nxt, ignore_errors=True)
    os.makedirs(repos, exist_ok=True)
    run([sys.executable, os.path.join(HERE, "stage.py"), nxt], HERE)
    run([sys.executable, os.path.join(HERE, "docs.py"), nxt], HERE)
    reqs = root_requirements()
    changed = []
    for r in ORDER:
        clone = os.path.join(repos, r)
        if not os.path.isdir(os.path.join(clone, ".git")):
            run(["git", "clone", "-q", f"https://github.com/{OWNER}/{r}.git", clone], repos)
        else:
            run(["git", "reset", "-q", "--hard"], clone)
            run(["git", "clean", "-fdq"], clone)
            run(["git", "pull", "-q", "--ff-only"], clone)
        replace_tree(os.path.join(nxt, "split", r), clone)
        for dep in DEPS.get(r, []):
            run(["go", "get", f"github.com/{OWNER}/{dep}@main"], clone, GOENV)
        run(["go", "mod", "tidy"], clone, GOENV)
        pin(clone, reqs)
        run(["go", "build", "./..."], clone, GOENV)
        run(["go", "vet", "./..."], clone, GOENV)
        run(["go", "test", "./..."], clone, GOENV)
        run(["git", "add", "-A"], clone)
        diff = run(["git", "diff", "--cached", "--stat"], clone)
        if not diff:
            print(f"{r}: up to date")
            continue
        changed.append(r)
        print(f"{r}: changed\n{diff}")
        if a.push:
            run(["git", "-c", f"user.name={NAME}", "-c", f"user.email={EMAIL}", "commit", "-q", "-m", a.message], clone)
            run(["git", "push", "-q", "origin", "main"], clone)
            print(f"{r}: pushed {run(['git', 'log', '--format=%h', '-1'], clone)}")
        else:
            run(["git", "reset", "-q"], clone)
    print("changed:", ", ".join(changed) or "none")


main()
