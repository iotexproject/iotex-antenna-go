#!/bin/bash
go build -v -x -work -ldflags "-s -w" -o didlib.a -buildmode='c-archive' -gcflags='-l'
