#!/bin/bash

# Don't really care that we'll have duplicates, the golang process will take care of that.

set -xe

# Random dictionary lists
curl --location https://github.com/dwyl/english-words/raw/master/words.txt > words.txt
curl --location https://archive.org/download/mobywordlists03201gut/SINGLE.TXT >> words.txt
curl --location https://archive.org/download/mobywordlists03201gut/ACRONYMS.TXT >> words.txt
curl --location https://archive.org/download/mobywordlists03201gut/COMPOUND.TXT >> words.txt
curl --location https://archive.org/download/mobywordlists03201gut/NAMES.TXT >> words.txt
# RockYou password list
curl --location https://github.com/brannondorsey/naive-hashcat/releases/download/data/rockyou.txt >> words.txt

