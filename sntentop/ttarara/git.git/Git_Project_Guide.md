
# 🌟 Git Project: Mastering Version Control

## 🚀 Introduction
This project introduces the fundamental concepts and workflows of **Git**, a powerful tool for version control and collaboration. By completing this project, you will:

- Gain hands-on experience with Git commands.
- Learn best practices for managing code repositories.
- Understand how to track changes and collaborate effectively.

---

## 📂  Setting Up the Project

### 🛠️  Create a Working Directory
```bash
mkdir work
cd work
```

### 🛠️ Create the `hello` Directory and `hello.sh` File
```bash
mkdir hello
cd hello
nano hello.sh
```

Add the following content to `hello.sh`:
```bash
echo "Hello, World"
```

### 🛠️ Initialize a Git Repository (--Inside hello folder--)
```bash
git init
```

### 🛠️ Check the Repository Status
```bash
git status
```

### 🛠️  Stage and Commit the `hello.sh` File
```bash
git add hello.sh
git commit -m "Initial commit with Hello World"
```

---

## ✏️ Modifying `hello.sh`

### ✏️  Update Content to Accept Arguments
Modify `hello.sh`:
```bash 
nano hello.sh 
```
```bash
#!/bin/bash
echo "Hello, $1"
```

### ✏️ Stage and Commit the Changes
```bash
git add hello.sh
git commit -m "Accept an argument"
```

### ✏️ Add a Default Value to `hello.sh`
Update `hello.sh`:
```bash 
nano hello.sh 
```
```bash
#!/bin/bash
# Default is "World"
name=${1:-"World"}
echo "Hello, $name"
```

### ✏️ Stage Interactively and Commit
```bash
git add -p hello.sh
git commit -m "Added a comment"
git add hello.sh
git commit -m "Add lines 4-5"
```

---
### 📜 View Commit History
```bash
git log
```

### 📜 Show a One-Line History
```bash
git log --oneline
```

### 📜 Display the Last 2 Commits
```bash
git log -n 2
```

### 📜 Show Commits from the Last 5 Minutes
```bash
git log --since="5 minutes ago"
```

### 📜 Customize the Log Format
```bash
git log --pretty=format:"%h %ad | %s (%H -> %d) [%an]" --date=short
```

---

### ⏪ Restore the First Snapshot
```bash
git checkout HEAD~3
```

### ⏪  View `hello.sh` Content
```bash
cat hello.sh
```

### ⏪  Restore the Second Snapshot
```bash
git checkout HEAD~2
cat hello.sh
```

---

### 🔖 Tag the Current Version as `v1`
```bash
git tag v1
```

### 🔖  Tag the Previous Version as `v1-beta`
```bash
git tag v1-beta c3baa32
```

### 🔖  Checkout Tags
```bash
git checkout v1
git checkout v1-beta
```

### 🔖  List All Tags
```bash
git tag
```

---

### ❌ Modify hello.sh with Unwanted Comments
```bash 
nano hello.sh
```
```bash 
#!/bin/bash
# This is a bad comment. We want to revert it.
name=${1:-"World"}
echo "Hello, $name"
```

## ❌  Revert Changes Before Staging
```bash
git checkout -- hello.sh
```

## ❌ Stage and Commit Unwanted Changes
Introduce unwanted changes and stage them:
```bash 
nano hello.sh 
```
```bash 
#!/bin/bash
# This is an unwanted but staged comment
name=${1:-"World"}
echo "Hello, $name"
```
## ❌ Add the comment:
```bash 
git add hello.sh
```

## ❌  Revert Staged Changes
```bash
git reset HEAD hello.sh
```
## ❌ Commit Unwanted Changes
```bash 
git commit -m "Unwanted change"
```

## ❌ Revert Committed Changes
```bash
git revert HEAD
```
---
## ❌ Tag the latest commit with "oops": 
```bash
git tag oops
```

## ❌ Remove Commits After v1:
```bash
git reset --hard v1
```

## ❌ Display Logs with Deleted Commits:
```bash 
git reflog
```
## ❌ Clean Unreferenced Commits:
git gc --prune=now

### 📂  Move `hello.sh` to `lib/`
```bash
mkdir lib
git mv hello.sh lib/
git commit -m "Moved hello.sh to lib/"
```

### 📂  Create and Commit a `Makefile`
Create a `Makefile`:
```bash 
nano Makefile 
```
```makefile
TARGET="lib/hello.sh"

run:
    bash ${TARGET}
```

Commit the file:
```bash
git add Makefile
git commit -m "Added Makefile"
```

---

### 🔍Explore the .git/ Directory
```bash 
cd .git
ls 
```
## 🔍Retrieve the Latest Commit Hash
```bash 
git rev-parse HEAD
```
## 🔍 Inspect Git Objects
```bash
git cat-file -t <commit_hash>
git cat-file -p <commit_hash>
```
## 🔍 Dump the Directory Tree
```bash 
git ls-tree <commit_hash>
```

## 🔍Dump the Contents of lib/ and hello.sh
```bash 
git cat-file -t <tree_hash_of_lib> 
```
```bash 
# Get the tree hash from ls-tree output
```
```bash 
git cat-file -p <blob_hash_of_hello.sh>  
```
```bash
Get the blob hash from ls-tree output
```
### 🌳 Create and Switch to a New Branch
```bash
git checkout -b greet
```

### 🌳 Add and Commit `greeter.sh`
```bash
nano greeter.sh
```
```bash
echo '#!/bin/bash
Greeter() {
    who="$1"
    echo "Hello, $who"
}' 

> lib/greeter.sh
```
```bash
git add lib/greeter.sh
git commit -m "Added greeter.sh"
```

### 🌳 Update `hello.sh` to Use `Greeter`:
```bash
nano hello.sh
```
Modify `hello.sh`:
```bash
#!/bin/bash
source lib/greeter.sh

name="$1"
if [ -z "$name" ]; then
    name="World"
fi

Greeter "$name"
```

Commit the changes:
```bash
git add lib/hello.sh
git commit -m "Updated hello.sh to use Greeter function"
```
## 🌳 Update the Makefile
```bash 
# Ensure it runs the updated lib/hello.sh file
TARGET="lib/hello.sh"

run:
	bash ${TARGET}
``` 
```bash 
git add Makefile
git commit -m "Added Makefile"
```

## 🌳 Switch back to the main branch:
```bash  
git checkout main
```

## 🌳 Compare differences between main and greet:
```bash 
git diff main greet -- Makefile

git diff main greet -- lib/hello.sh

git diff main greet -- lib/greeter.sh
```
## 🧩 Create README.md:
```bash 
echo "This is the Hello World example from the git project." > README.md
``` 

```bash 
git add README.md

git commit -m "Added README.md"
```

### 📈  Visualize with a Graph
```bash
git log --graph --oneline --all
```
```bash 
ttarara@Tarara-PC:/mnt/c/Users/xarou/Desktop/git/work/hello$ git log --graph --oneline --all
*c702842 (HEAD -> greet, origin/main, origin/greet, main) Updated README.md in the original repository
*   bebd1d1 Resolved merge conflict
|\  
| * 247cc46 hello.sh
* | efa1763 hello.sh
* | dfeaa37 Conflict create
|/  
* f18c612 Interactive hello.sh
*   f5c03de Merge branch 'main' into greet
|\  
| * b005cd5 Added README.md
* | efd54dc Updated Makefile with a comment
* | 6aa6cc3 Updated hello.sh to use Greeter function
* | 656049c Added greeter.sh
* | cf5a1ea Added Makefile
* | f16511a Moved hello.sh to lib/
|/  
* db1cf9d (tag: v1) 3rd Comment 4-5 lines
* 757ef8e Add comment line 3
| * 8ef0e46 (tag: oops) Revert "Revert "Unwanted change""
| * d158c70 Revert "Unwanted change"
| * 54b91dc Unwanted change
|/  
* c3baa32 (tag: v1-beta) 2nd Initial commit with Hello World
* 96ca8c5 1st Initial commit with Hello World

---
```

### 🧬 Merge main into greet:

```bash 
git checkout greet
git merge main
```

### 🧬 Switch to the main branch and make changes:

```bash 
git checkout main
```
```bash 
 Modify lib/hello.sh:
 nano hello.sh 
 ```
 ```bash 
#!/bin/bash

echo "What's your name"
read my_name

echo "Hello, $my_name"

```
```bash 
git add lib/hello.sh
git commit -m "Interactive hello.sh"
```

### 🧬 Merge main into greet (creating a conflict):

```bash 
git checkout greet
git merge main (This should result in a merge conflict)
```
### 🧬 Resolve the conflict:

```bash 
Open lib/hello.sh in a text editor and resolve the conflict by choosing the changes from the main branch.

git add lib/hello.sh
git commit -m "Resolved merge conflict"
``` 
### 🧬 Rebase greet onto main:git checkout greet

```bash 
git rebase main
```
- Merge greet into main:git checkout main
```bash 
git merge greet
``` 
## 🧬Explain fast-forwarding and the difference between merging and rebasing.
```bash 
┃ Fast-forwarding: Imagine merging two lanes of traffic with no cars in one lane. 
┃ You just speed up! That's fast-forwarding - a simple merge with no obstacles.

Merging ⚔️ Rebasing:

┃ Merging: Like combining two roads, keeping all the original paths visible.
┃ Rebasing: Like moving a side road to connect to the main road's end, creating a single, straight path.
┃ Merging keeps all history, rebasing streamlines it. Choose based on your project's needs!
```

## 👥 Clone the repository
```bash 
git clone hello cloned_hello
```
- Show logs for the cloned repository
```bash 
cd cloned_hello
git log
```

## 👥 Display remote repository information
```bash 
git remote show origin
``` 

## 👥 List all branches (local and remote)
```bash 
git branch -a
```

## 👥 Make changes to the original repository 
```bash 
cd ../hello
echo "This is the Hello World example from 
the git project."  -----> (Changed in the original) > README.md 
git add README.md
git commit -m "Updated README in original"
```

## 👥  Fetch changes in the cloned repository
```bash 
cd ../cloned_hello
git fetch origin
git log --decorate --graph --oneline --all 
ttarara@Tarara-PC:/mnt/c/Users/xarou/Desktop/git/work/cloned_hello$ git log --decorate --graph --oneline --all 
* c702842 (HEAD -> main, origin/main, origin/greet, origin/HEAD, greet) Updated README.md in the original repository
*   bebd1d1 Resolved merge conflict       
|\  
| * 247cc46 hello.sh
* | efa1763 hello.sh
* | dfeaa37 Conflict create
|/  
* f18c612 Interactive hello.sh
*   f5c03de Merge branch 'main' into greet
|\  
| * b005cd5 Added README.md
* | efd54dc Updated Makefile with a comment
* | 6aa6cc3 Updated hello.sh to use Greeter function
* | 656049c Added greeter.sh
* | cf5a1ea Added Makefile
* | f16511a Moved hello.sh to lib/
|/
* db1cf9d (tag: v1) 3rd Comment 4-5 lines
* 757ef8e Add comment line 3
| * 8ef0e46 (tag: oops) Revert "Revert "Unwanted change""
| * d158c70 Revert "Unwanted change"
| * 54b91dc Unwanted change
|/
* c3baa32 (tag: v1-beta) 2nd Initial commit with Hello World
* 96ca8c5 1st Initial commit with Hello World
```

## 👥 Merge changes from remote main
```bash 
git merge origin/main 
```

## 👥 Add a local branch tracking origin/greet
```bash 
git checkout -b greet origin/greet
```

## 👥 Add a remote repository > Replace with your actual remote URL
```bash 
git remote add myremote <your_remote_url>
```

## 👥 Push main and greet branches to the new remote
```bash 
git push myremote main
git push myremote greet 
``` 

## 👥 "What is the single git command equivalent to what you did before to bring changes from remote to local main branch?"
```bash 
git pull origin main
```
> This command elegantly combines two essential steps:

- - Fetching: It retrieves the latest changes from the main branch of the origin remote repository.
- - Merging: It seamlessly integrates those fetched changes into your local main branch.
```bash 
In essence, git pull acts as a streamlined shortcut, ensuring your local branch 
stays up-to-date with the remote repository effortlessly.
```

### 🧸 Create a bare repository
```bash 
git clone --bare hello hello.git
# Add hello.git as a remote to the original repository
cd hello
git remote add shared ../hello.git
```
## 🧸 Change README.md
```bash 
echo "This is the Hello World example from the git project. "
> Changed in the original and pushed to shared  > README.md

# Commit the changes
git add README.md
git commit -m "Updated README and pushed to shared"
``` 
## 🧸 Push changes to the shared repository
``` bash 
git push shared main
``` 
## 🧸 What is a bare repository? Why is it needed?
```bash 
- A bare repository is a Git repository that doesn't have a working directory. 
It only contains the version history (the .git directory contents) and is used 
primarily as a central repository for sharing code among developers.
It's like a central hub where everyone pushes their changes and pulls updates from.

* Collaboration: Bare repositories provide a central location for developers to share their code and collaborate on projects.
* Avoid Conflicts: Since there's no working directory, there's no risk of merge conflicts arising from direct changes in the bare repository itself.
* Efficient Storage: Bare repositories are optimized for efficient storage and transfer of Git data.What is a bare repository?
```
### 🌐 Create a Repository Online
> Go to a Git hosting service (e.g., GitHub, GitLab, or Bitbucket).
> Create a new repository named git.
```bash 
Push to the main Branch

git checkout main
git remote add origin htte01.gr/git/ttarara/git.git
git add .
git commit -m "Git Project"
git push -u origin main
```
```bash 
Push the greet Branch

git checkout greet
git remote add origin htte01.gr/git/ttarara/git.git
git add .
git commit -m "Git Project greet"
git push -u origin greet
```