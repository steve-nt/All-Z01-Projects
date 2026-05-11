# How to merge repos

```cd /mnt/c/Users/(user name)/Desktop/(folder) ```
    - Tip: when I use linux environment from windows and use wsl in vs code, this is the path I should use.
## Remove Submodule
🪜 Step-by-step:

1. Remove the submodule link (but keep files locally):

```git rm --cached practisetest```


2. Delete the .git folder inside practisetest (this removes its separate Git history):

```rm -rf practisetest/.git```

3. Add it again as normal files:

```git add practisetest```


4. Commit the change:

```git commit -m "fix practisetest folder (remove submodule)"```


5. Push again:

```git push```


✅ After this, practisetest will be a normal directory inside your main repo — no more submodule issues.

### Error 

```fatal: not a git repository (or any parent up to mount point /mnt)```

### Meaning & Solution

This means:

“Hey, I don’t see a .git folder here — you’re not inside a Git repo right now.”

You need to go inside the repo folder first (the one that actually has .git/ inside it).

Let’s find it and then run the command again.

🪜 Step 1. Move into your repo

using the above path
You can check you’re in the right place by listing hidden files:

```ls -a```

You should see something like:

```.git  (files name)  README.md  ...```


If you see the .git folder — perfect ✅

🪜 Step 2. Remove the submodule (inside the repo)

Now run:

```git rm --cached practisetest```

🪜 Step 3. Delete the inner .git folder
```rm -rf practisetest/.git```

🪜 Step 4. Add the folder back as normal files
```git add (file name)```

🪜 Step 5. Commit and push
```git commit -m "fix practisetest(file name) folder (remove submodule)"```
```git push``` 

## Merge without preserving the history

🪜 Step 1. Prepare your new “main” repo

Go to your **working folder** (where you want the new repo to live):

```cd /mnt/c/Users/(user name)/Desktop/(folder)```
```mkdir (repo's name)```
```cd (repo's name)```
```git init```

Now you have an empty Git repository ready to receive your projects.

🪜 Step 2. Copy your 4 repos inside (without their .git folders)

Let’s assume your old repos are here:

```/mnt/c/Users/gener/Desktop/ZoneGit/repo1```

Then copy them like this:

```cp -r /mnt/c/Users/gener/Desktop/ZoneGit/repo1 ./repo1```
```cp -r /mnt/c/Users/gener/Desktop/ZoneGit/repo2 ./repo2```
```cp -r /mnt/c/Users/gener/Desktop/ZoneGit/repo3 ./repo3```
```cp -r /mnt/c/Users/gener/Desktop/ZoneGit/repo4 ./repo4```

🧠 Important: Each of these repos has its own .git folder, which causes the submodule problem.
We don’t want those.

So run:

```rm -rf repo1/.git repo2/.git repo3/.git repo4/.git```


Now they’re just regular folders with files, not submodules.

🪜 Step 3. Add and commit everything
```git add .```
```git commit -m "Add all four projects into one organized repo"```

🪜 Step 4. Create a remote (on Zone01 or GitHub)

If you haven’t yet:

```git remote add origin https://platform.zone01.gr/git/mkapola/mega-repo.git```


Then push:

```git push -u origin main```


✅ Done!
Now you’ll have one beautiful, clean repo with all 4 projects neatly inside.

## Merge and preserve the history

That’s possible too, using git subtree.
It’s a bit more advanced, but very powerful:

**clone** each of your existing **repos** to your **local folder** first,
so that you have the actual files and their Git history on your computer.

Then, you’ll use those local copies as the “sources” for your new combined mega-repo.

Example:

```git remote add repo1 /path/to/repo1```
```git fetch repo1```
```git subtree add --prefix=repo1 repo1 main```


Do this for each repo (repo2, repo3, repo4).
That way, you keep all commit history inside the new combined repo.

Now, here’s what each command does 👇

🪄 1. ```git remote add repo1 /path/to/repo1```

What it does:
Adds your old repo (repo1) as a remote source to your new combined repo.

🧠 Think of it like saying:

“Hey Git, besides my current repo, also keep an eye on that other one called repo1 — it lives at this path.”

🔍 After this, you can fetch and see the commits from that repo.

🗂 Example:

If your folder structure looks like this:

mega-repo/
repo1/


You’d run:

```git remote add repo1 ../repo1```


Now when you type:

```git remote -v```


You’ll see something like:

```repo1   ../repo1 (fetch)```
```repo1   ../repo1 (push)```
```origin  https://platform.zone01.gr/git/mkapola/mega-repo.git (fetch)```
```origin  https://platform.zone01.gr/git/mkapola/mega-repo.git (push)```

📥 2. git fetch repo1

What it does:
Downloads all the commit history and branches from that remote (the old repo)
→ but doesn’t yet merge or change your files.

🧠 Think of it like:

“Get me everything from repo1, so I can use it locally if I want.”

You’ll now have all of repo1’s commits, branches, and tags available in your current repo’s database.

🗂 After fetching, you can check what you got:
```git branch -r```


You’ll see something like:

repo1/main
origin/main

🧱 3. git subtree add --prefix=repo1 repo1 main

What it does:
Takes the contents of the main branch from the remote called repo1,
and merges it into your current repo — placing all its files inside a subfolder named repo1/.

✅ It also preserves the entire commit history of that repo inside your current repo.
That’s the magic of git subtree ✨

🧠 Conceptually:

You’re saying:

“Take the project from repo1’s main branch, and graft (attach) it into this repo, under a folder named repo1/.”

The result:

mega-repo/
└── repo1/
    ├── file1.go
    ├── file2.go
    └── ...


And all of repo1’s commits are now part of your mega-repo’s history.

🧩 So in summary
Command	Meaning	Effect
```git remote add repo1 /path/to/repo1```	Connect your old repo to the new one as a remote	Lets you access its commits
```git fetch repo1```	Download that repo’s history	Makes its branches available locally
```git subtree add --prefix=repo1 repo1 main```	Merge that repo’s files (and history) under repo1/ folder	Combines everything cleanly
⚡️ Bonus Tip

You can repeat this for each repo:

```git remote add repo2 /path/to/repo2```
```git fetch repo2```
```git subtree add --prefix=repo2 repo2 main```


Now you’ll have:

mega-repo/
├── repo1/
├── repo2/
├── repo3/
└── repo4/


All with their own commit histories preserved 🎯

4. Push your new combined repo to Zone01

Once everything looks right:

```git remote add origin https://platform.zone01.gr/git/mkapola/mega-repo.git```
```git push -u origin main```

### Error 

```remote (file name) already exists.```

### Meaning & Solution

🧩 What’s happening

You already have a remote called (file name) (even though the path was wrong).

The real folder name is (file name) (wrong spelling)

That’s why Git couldn’t find the repo — the path you gave doesn’t exist.

✅ Step-by-step fix

1. Remove the wrong remote:

```git remote remove practicetest```


2. Add the correct one (notice the correct spelling):

```git remote add practisetest ~/Desktop/ZoneGit/(file name)```


3. Fetch the repository data:

```git fetch practisetest```


4. Add it as a subtree inside your main repo:

```git subtree add --prefix=practisetest practisetest main```

🧠 What this does

```git remote add practisetest``` … → Tells your main repo where to find the practisetest repo.

```git fetch practisetest``` → Downloads its history into your main repo.

```git subtree add``` … → Copies all its content into a subfolder inside your main repo while preserving commit history.

After this, run:

```ls```


and you should see:

```README.md  practisetest/```


Then you can safely commit and push:

```git add .```
```git commit -m "Add practisetest repo as subtree"```
```git push```

### Error 

🧩 Error explained
```fatal: ambiguous argument 'HEAD': unknown revision or path not in the working tree.```
```fatal: working tree has modifications.  Cannot add.```


This happens for two main reasons:

-Your current repo (piscinepractice-repo) is empty but has a commit or files that Git thinks are “modifications.”

-git subtree add requires a clean working tree — no uncommitted changes.

**HEAD is ambiguous**

If your repo was just initialized and has no commits yet, Git doesn’t know what HEAD is.

**Subtree needs at least one commit in the main repo to work.**

✅ Step-by-step fix
1. Make sure your main repo has at least one commit

From piscinepractice-repo:

```git status```


If it shows files not staged for commit, stage them:

```git add .```
```git commit -m "Initial commit for main repo"```


✅ Now your repo has a HEAD and a clean working tree.

2. Make sure there are no uncommitted changes
```git status```


Should show: ```nothing to commit, working tree clean```

If not, commit or stash changes first:

```git add .```
```git commit -m "Save changes before subtree"```

3. Add the subtree
```git subtree add --prefix=(file name) (file name) main```


```--prefix=(file name)``` → where the files from the subtree will go inside your main repo

```(file name)``` → the remote we added

```main``` → the branch of the subtree repo we want to merge

🧠 Optional check

After adding, verify:

```ls (file name)```
```git log --oneline --graph --all```


You should now see all the files from practisetest and its commit history.

💡 Tip:
Every time you use git subtree add, your main repo must be clean and must have at least one commit. Otherwise Git doesn’t know where to attach the subtree.


Make sure you are pointing to the actual test repo folder, not piscinepractice-repo.
For example, if your folder structure is:

ZoneGit/
├── piscinepractice-repo/
├── practisetest/
├── squad/
├── **test/**         ← THIS is the separate repo


Then your commands should be:

```git remote add test /mnt/c/Users/gener/Desktop/ZoneGit/test```
```git fetch test```
```git subtree add --prefix=test test main```


--prefix=test → folder name in your main repo

test → remote name

main → branch of that repo (check with git branch -a in test if needed)

🧠 Double-check

Before subtree add, you can check the branches of the test repo:

```cd /mnt/c/Users/gener/Desktop/ZoneGit/test```
```git branch```


Make sure there is a branch called main (or master — use the correct name).

If the branch is master, then the subtree command should be:

```cd /mnt/c/Users/gener/Desktop/ZoneGit/piscinepractice-repo```
```git subtree add --prefix=test test master```