#!/usr/bin/env bash
build_time="`date +%Y-%m-%d,%H:%M:%S`"
echo $build_time
build_git_hash=`git rev-parse HEAD`
echo $build_git_hash

rm -rf rocketMQ

GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.BuildTime=${build_time} -X main.BuildGitHash=${build_git_hash}" -o rocketMQ
