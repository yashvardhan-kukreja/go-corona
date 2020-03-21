#!/bin/bash
rm -rf go-corona
go build go-corona.go
sudo cp go-corona /usr/local/bin