#!/usr/bin/env bash
# Git setup, podman image builder and launcher for vD system
#
# Author:  もとない (motonai)
# Email: 217199434+motonai@users.noreply.github.com
# License: GPLv3

function _vd_launch_usage(){
    printf "./%s - vD launcher\n\n" "$(basename $0)" >&2
    printf "./%s -h\n" "$(basename $0)" >&2
    printf "./%s --help\n" "$(basename $0)" >&2
    printf "./%s -c <username> <email>\n" "$(basename $0)" >&2
}
username=""
email=""
if [ "$1" == "-h" ] || [ "$1" == "--help" ]
    then
        _vd_launch_usage
        exit 1
    fi

if [ $# -eq 3 ]
then
    if [ "$1" != "-c" ]
    then
        _vd_launch_usage
        exit 1
    fi
    username="$2"
    email="$3"
else
    username="$(grep 'Author:' $(realpath $0) | head -n 1 | cut -d ' ' -f 3-)"
    email="$(grep 'Email:' $(realpath $0) | head -n 1 | cut -d ' ' -f 3-)"
fi

if [ ! -n "$(git config --global core.editor)" ]
then
    git config --global core.editor vim
fi
if [ ! -n "$(git config --global credential.helper)" ]
then
    git config --global credential.helper store
fi
if [ ! -n "$(git config --global user.email)" ]
then
    git config --global user.email "${email}"
fi
if [ ! -n "$(git config --global user.name)" ]
then
    git config --global user.name "${username}"
fi
if [ ! -n "$(git config --global init.defaultBranch)" ]
then
    git config --global init.defaultBranch main
fi

mv ${HOME}/.gitconfig $(pwd)/image/.gitconfig
mv ${HOME}/.git-credentials $(pwd)/image/.git-credentials

podman build -f ./image/ContainerFile -t vd

rm $(pwd)/image/.gitconfig
rm $(pwd)/image/.git-credentials

if [ ! -d $HOME/vd_repositories_dir ]
then
    mkdir $HOME/vd_repositories_dir
fi

podman run -it \
    --rm \
    --name vdi \
    -h vdi \
    -p 8181:8181 \
    --network=host \
    --security-opt label=disable \
    --volume=$HOME/vd:/root/vd:rw \
    --volume=$HOME/vd_repositories_dir:/root/vd_repositories_dir:rw \
    --volume=/etc/localtime:/etc/localtime:ro \
    vd
