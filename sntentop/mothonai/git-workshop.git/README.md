# git-workshop

Documents about a git workshop

## Build

### Podman
If you have `podman` installed you can build the documentation with the script:
```
./compile-with-podman.sh
```

### Fedora/RHEL

Install dependencies and make the html format
```
dnf install python3-sphinx_rtd_theme make -y
make html
```

### Ubuntu / WSL
```
sudo apt-get install python3-sphinx python3-sphinx-rtd-theme
make html
```

### Fedora Bluefin (using toolbox)

Clone the repository
```bash
git clone https://platform.zone01.gr/git/mothonai/git-workshop
cd git-workshop
```

Create a toolbox
```
toolbox create
```

Enter the toolbox you just created
```
toolbox enter
```

Install dependenties
```bash
sudo dnf install -y python3-sphinx_rtd_theme.noarch
sudo dnf install -y make
```

Build the project
```bash
$ make html
```

When the project is build close the toolbox
```bash
$ exit
```

## Make use of the documentation
Both of the above, will output the documentation in HTML format under the
directory `./build/html/`.

From now on, you can read the project's documentation in the browser
directory and starting the local server with the following command

```bash
$ python3 -m http.server 8000
```
