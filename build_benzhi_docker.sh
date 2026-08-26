#!/bin/sh
set -eu

docker build -f benzhi.Dockerfile -t prod23-go-project:benzhi .
