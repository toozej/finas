#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/finas/ man | gzip -c -9 >manpages/finas.1.gz
