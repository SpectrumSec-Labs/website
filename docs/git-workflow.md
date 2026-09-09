# Git workflow

Practical command reference for day-to-day work on this repo.

- Remote: `origin` → `https://github.com/SpectrumSec-Labs/website.git`
- Default branch: `main` — **protected**: PR required, CI must pass, linear
  history, no force-push, no direct commits.
- PRs are **squash-merged** on GitHub. One PR = one commit on `main`.
- After merge, delete the branch on GitHub (there's a button on the PR).

---

## 1. Ship a feature

```sh
# start from an up-to-date main (see §2 for the full sync)
git switch main
git pull

# branch — use the repo's prefixes: content/  fix/  feat/  chore/  test/
git switch -c feat/short-slug

# ... edit ...

# build must be clean before you push (see §5)
./bin/hugo.exe --gc --minify --panicOnWarning

git add -A
git commit -m "feat: what changed and why"
git push -u origin feat/short-slug     # -u only needed the first push
```

Then open a PR on GitHub, get it green, squash-merge it, and click
**Delete branch**.

---

## 2. Sync local `main` after your PR is merged

Because PRs are **squash-merged**, the commit on `main` is a *new* commit your
local branch has never seen. Your local `feat/...` branch will look "unmerged"
to git even though the change is in `main` — that's expected.

```sh
git switch main
git pull --ff-only          # fast-forwards main to include the squash commit
git fetch --prune           # deletes stale origin/<branch> refs (branches gone from GitHub)
git branch -D feat/short-slug   # -D (not -d): squash-merge means git can't see it as merged
```

`git pull --ff-only` refuses if local `main` has diverged (it never should —
you don't commit to `main` directly). If it refuses, stop and look:
`git log --oneline origin/main..main` shows what's local-only.

### One-shot cleanup of every branch deleted on the remote

After `git fetch --prune`, branches whose upstream is gone show `: gone]` in
`git branch -vv`. Delete them all (PowerShell):

```powershell
git branch -vv | Where-Object { $_ -match ': gone\]' } |
  ForEach-Object { git branch -D $_.Substring(2).Trim().Split(' ')[0] }
```

bash/zsh:

```sh
git branch -vv | awk '/: gone]/ {print $1}' | xargs -r git branch -D
```

> Check the list first — a `: gone]` branch that still has unpushed commits
> (`git log --oneline main..<branch>`) will lose that work. `git reflog` keeps
> it recoverable for ~90 days regardless.

---

## 3. Keep a long-running branch current with `main`

```sh
git switch feat/big-thing
git fetch origin
git rebase origin/main       # linear history — rebase, don't merge
# resolve conflicts, then: git rebase --continue
git push --force-with-lease  # your own branch only, never main
```

---

## 4. First-time / fresh-clone setup

```sh
git clone https://github.com/SpectrumSec-Labs/website.git
cd website
git config pull.ff only          # never create surprise merge commits on pull
git config fetch.prune true      # every fetch/pull prunes gone remote refs
git config branch.autoSetupMerge always
```

The `bin/` toolchain (hugo.exe, go/, etc.) is gitignored — install it
separately per `CLAUDE.md` → "Local tooling".

---

## 5. Before every push

```sh
./bin/hugo.exe --gc --minify --panicOnWarning   # CI fails on any deprecation/warning
```

If you touched `tools/newsfetch`:

```sh
cd tools/newsfetch
GOROOT=../../bin/go ../../bin/go/bin/go.exe test -mod=vendor ./...
GOROOT=../../bin/go ../../bin/go/bin/go.exe vet  -mod=vendor ./...
gofmt -l $(git ls-files '*.go' | grep -v /vendor/)   # must print nothing
```

---

## 6. Undo / recover

| Situation | Command |
|---|---|
| Unstage a file | `git restore --staged <file>` |
| Discard uncommitted changes to a file | `git restore <file>` |
| Discard **all** uncommitted changes | `git reset --hard && git clean -fd` |
| Amend the last commit (before push) | `git commit --amend` |
| Move last commit to a new branch | `git switch -c feat/x` then `git switch - && git reset --hard HEAD~1` |
| Undo a `pull`/`reset` that moved a branch wrongly | `git reflog`, find the good SHA, `git reset --hard <sha>` |
| Recover a deleted branch | `git reflog`, then `git branch <name> <sha>` |
| Stash work to switch tasks | `git stash push -m "wip"` … `git stash pop` |

---

## 7. Inspecting state

```sh
git status                             # working tree
git branch -vv                         # local branches + upstream + ahead/behind
git log --oneline --graph --all -20    # who's where
git log --oneline origin/main..main    # commits local main has that remote doesn't (should be empty)
git log --oneline main..origin/main    # commits remote main has that you don't
git diff origin/main...HEAD            # what your branch changes vs main
```

---

## 8. The automated news branch

`origin/automated/news` is opened by the `update-news` workflow to refresh
`data/news.json`. It auto-merges once CI is green. Don't branch from it, don't
commit to it — if it has conflicts, re-run the workflow rather than fixing by
hand.

---

## 9. Deploy tags

`release.yml` creates tags like `deploy-20260909T130427Z-caf50e2` on every
release build. They're informational — the VPS pulls signed GitHub Releases,
not tags. Fetch them with `git fetch --tags`; don't create or move them by hand.
