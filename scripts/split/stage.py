import os, re, shutil, subprocess, sys
ROOT = "/Users/alex/Work/shadcn.go"
S = sys.argv[1]
OUT = os.path.join(S, "split")
OWNER = "github.com/alexisbcz"
YEAR = "2026"
HOLDER = "Alexis Bouchez"

# repo -> list of (src dir relative to ROOT, dest dir relative to repo root)
REPOS = {
    "libhtml": [("libhtml", "")],
    "cn": [("cn", "")],
    "cva": [("cva", "")],
    "twmerge": [("twmerge", "")],
    "tailwind": [("tailwind", ""), ("cmd/tailwind", "cmd/tailwind")],
    "icons": [("components/icons", "")],
    "components": [("components", "")],
}
EXCLUDE = {"components": ["icons"]}

def rewrite_imports(text):
    text = text.replace('shadcn.go/components/icons', OWNER + "/icons")
    text = text.replace('shadcn.go/components', OWNER + "/components")
    text = text.replace('shadcn.go/cmd/tailwind', OWNER + "/tailwind/cmd/tailwind")
    for m in ["libhtml", "cn", "cva", "twmerge", "tailwind"]:
        text = re.sub(r'shadcn\.go/' + m + r'\b', OWNER + "/" + m, text)
    return text

MIT = """MIT License

Copyright (c) {year} {holder}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
""".format(year=YEAR, holder=HOLDER)

for repo, parts in REPOS.items():
    dst = os.path.join(OUT, repo)
    if os.path.exists(dst):
        shutil.rmtree(dst)
    os.makedirs(dst)
    for src, sub in parts:
        ignore = shutil.ignore_patterns(*EXCLUDE.get(repo, [])) if not sub else None
        target = os.path.join(dst, sub) if sub else dst
        if sub:
            shutil.copytree(os.path.join(ROOT, src), target, dirs_exist_ok=True)
        else:
            for name in os.listdir(os.path.join(ROOT, src)):
                if name in EXCLUDE.get(repo, []):
                    continue
                p = os.path.join(ROOT, src, name)
                if os.path.isdir(p):
                    shutil.copytree(p, os.path.join(dst, name))
                else:
                    shutil.copy2(p, os.path.join(dst, name))
    for dp, dn, fn in os.walk(dst):
        for f in fn:
            if f.endswith((".go", ".js", ".mjs", ".md")):
                p = os.path.join(dp, f)
                with open(p, encoding="utf-8") as fh:
                    t = fh.read()
                n = rewrite_imports(t)
                if n != t:
                    with open(p, "w", encoding="utf-8") as fh:
                        fh.write(n)
    with open(os.path.join(dst, "LICENSE"), "w") as fh:
        fh.write(MIT)
print("staged", ", ".join(REPOS))
