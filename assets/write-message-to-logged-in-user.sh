#!/bin/bash

if (( $# != 2 )); then
	echo "Need 2 arguments, not $#" >&2
	exit 1
fi

user=$(echo "$1" | cut -d " " -f1)
tty=$(echo "$1" | cut -d " " -f2)
message="$2"

echo "$message" | write "$user" "$tty"
