#!/bin/bash
w | awk '{ print $1 " " $2 }' | grep -v ":" | grep -v "USER TTY" | grep -v '^$' 
