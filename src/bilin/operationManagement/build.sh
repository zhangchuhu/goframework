#!/bin/bash

echo "Begin go-bindata static..."
go-bindata-assetfs static/...

make clean
echo "Begin make project!"
make tar
