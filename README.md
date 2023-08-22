# ELRS Rainbow Table Generator

## Description

This is a set of tools to reverse an ExpressLRS binding UID back into the binding phrase that generated it. It works only with simple binding phrases.

## Usage

```bash
# Get the words.txt file
$ ./get_wordstxt.sh
# Build the rainbow table (takes 1-2 minutes)
$ go run build_table.go
Restored 0 entries from CSV
Discovered 617981 new entries
Wrote 617981 entries to CSV
# Look up a binding phrase (takes 30 seconds)
$ go run lookup.go 180,155,166,195,245,227
[...] Looking up binding phrase for 180,155,166,195,245,227
[...] Loaded 14890190 entries from CSV
[+++] Found binding phrase for 180,155,166,195,245,227: spektrum
```

## Warnings

Once downloaded, words.txt and found.csv will contain profanities and other nasty things by the nature of the way language works.
