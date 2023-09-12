#!/bin/bash
echo $(date +%y.%m).$(git rev-list --count HEAD --since="$(date +'%Y-%m-01')") >version.txt
