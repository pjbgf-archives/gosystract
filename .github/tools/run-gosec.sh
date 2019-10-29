#!/usr/bin/env bash
set -e
set -o pipefail

wget https://github.com/securego/gosec/releases/download/2.0.0/$GOSECNAME.tar.gz -O /tmp/$GOSECNAME.tar.gz
echo "490c2a0434b2b9cbb2f4c5031eafe228023f1ac41b36dddd757bff9e1de76a2b /tmp/$GOSECNAME.tar.gz" | sha256sum -c -

mkdir -p /tmp/$GOSECNAME
tar -C /tmp/$GOSECNAME -xzf /tmp/$GOSECNAME.tar.gz

/tmp/$GOSECNAME/gosec -conf gosec.json ./...