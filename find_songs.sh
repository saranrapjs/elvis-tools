#!/bin/bash
# usage: cat file.txt | ./find_songs.sh
regex="\"[A-Z](?:[A-Za-z0-9',]+[^A-Za-z0-9',\"]*){1,10}\""
cat - | grep -o -P $regex | sort | uniq
