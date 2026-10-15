#!/bin/zsh
# publish.sh <repo>: one "initial commit", new public repo alexisbcz/<repo>, push.
set -e
r=$1
S=${SPLIT_DIR:?set SPLIT_DIR to the directory holding split/}
cd $S/split/$r
rm -rf .git
git init -q -b main
git add -A
git -c user.name="Alexis Bouchez" -c user.email="276929831+alexisbcz@users.noreply.github.com" commit -q -m "initial commit"
gh repo create alexisbcz/$r --public --source . --push >/dev/null
git log --format='%h %an <%ae> %s' -1
