#!/bin/bash

cd
mkdir -p implant
cd implant

screen -d -m -S hosting_implant python -m SimpleHTTPServer 8080
exit 0
