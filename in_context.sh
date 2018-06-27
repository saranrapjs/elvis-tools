#!/bin/sh
while read line
do
  res=$(grep --color -n -m 1 "$line" elvis.txt)
  echo $res"; [["$line"]]"
done < "${1:-/dev/stdin}"
# cat filtered.txt | ./in_context.sh | sort -t: -n -k1
