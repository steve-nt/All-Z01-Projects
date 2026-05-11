# vD (ex-dotfiles)

> vD means virtualDimension

## Provides
A complete setup using podman containers which sets up:
 - go
 - nodejs
 - vim
 - git
 - bash
 - tmux
 - curl
 - jq
 - pinentry
 - screen
 - tree
 - which
 - wget
 - zsh

Also, provides a directory for repositories so one can reach those via the host
system. It is created at `$HOME/vd_repositories_dir/` and it mounts at
`/root/vd_repositories_dir` inside the container.

## Requirements
 - podman

## Launch
To launch the environment follow the steps:
```
git clone https://platform.zone01.gr/git/mothonai/vd
cd vd
./launch.sh -c <username> <email>
```

## Additional scripts

There are scripts accompanying the repository. These are stored inside the
`./bin` directory and the directory is included in the PATH environment variable
so you can execute them directly like:
```
er
```

### Script descriptions
 - `er` can be used to clone directly from platform.zone01.gr/git via HTTPS
```
$ er
Usage:
	er <username> <repo>
	er <repo>
```
 - `mkc` is a wrapper of `mkdir -p <dir> && cd <dir>`
```
$ mkc
Usage:
    mkc <directory>
```
- `eru` is a git fetcher that fetches all the updates from your remotes for
  repositories you cloned with `er`. Below an example of its output.
```
$ eru
Fetching repo updates:
Done username/repo1
Fail user2/repo1
```
