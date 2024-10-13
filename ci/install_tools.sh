#!/bin/bash

check_input_present() {
    local in=$1
    local err=$2
    if [ -z "$in" ]; then
        echo "error: $err"
        exit 1
    fi
}

check_go_install() {
    go version
}

check_go_task_install() {
    task --version
}

check_go_templ_install() {
    templ --version
}

install_go() {
    local version=$1
    check_input_present $version "missing go version"

    local url="https://golang.org/dl/go${version}.linux-amd64.tar.gz"
    local tarball="go${version}.linux-amd64.tar.gz"
    wget $url -O $tarball
    rm -rf /usr/local/go
    tar -C $HOME -xzf $tarball
    export PATH=$PATH:$HOME/go/bin
    rm $tarball

    check_go_install
}

install_go_task() {
    local version=$1
    check_input_present $version "missing go-task version"

    local go_task_url="github.com/go-task/task/v3/cmd/task@$version"
    go install $go_task_url

    check_go_task_install
}

install_go_templ() {
    local version=$1
    check_input_present $version "missing go-templ version"

    local go_templ_url="github.com/a-h/templ/cmd/templ@$version"
    go install $go_templ_url

    check_go_templ_install
}

install_go $1
install_go_task $2
install_go_templ $3
