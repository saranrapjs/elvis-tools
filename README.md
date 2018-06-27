# elvis-tools

This is a grab bag of scripts I used for combing thru a digital text & finding songs. [Here's a long-winded description of the whole thing.](http://regexking.info/2018/06/25/elvis-tools.html)

## Prerequisites

These are alternatingly bash, node and Go scripts. The Go scripts in particular have some dependencies, which should be fetch-able using `go get ./...`.

What follows are the steps in the order I followed them.


### `find_songs.sh`

```bash
cat big-input-file.txt | ./find_songs.sh > 1-possible-songs.txt
```

### `mostly_capitals.js`

Takes the grep result and excludes anything that's not mostly capitalized (crude "songfulness" heuristic).

```bash
cat 1-possible-songs.txt | ./mostly_capitals.js > 2-mostly-capitals.txt
```

### `in_context.sh`

Takes the filtered grep result, greps again with line numbers in the source text, and (with the following snippet) sorts/unique-ifies things so you end up with a unique list of mmatches in the order they appeared in the text. It also includes the originally matched song name in brackets (it's used in the next step).

```bash
cat 2-mostly-capitals.txt | ./in_context.sh | sort -t: -n -k1 > 3-in-context.txt
```

### `digfornames`

This takes the ordered matches and, for each one, interactively prompts you to select a nearby artist-like name. The result is piped out to stdout.

```bash
cat 3-in-context.txt | go run digfornames/main.go > 4-songs-with-artists.txt
```

### `findyoutube`

This takes a list of songs with artists, searches YouTube for them, and runs a little server which interactively presents the results and prompts you to select the one you want. The result is piped out to stdout.

```bash
cat 4-songs-with-artists.txt | go run findyoutube/main.gog > 5-youtube-urls.txt
```

### Bonus command

This isn't part of `elvis-tools`, but say you had a list of YouTube URL's, you could use this command to convert them to mp3's:

```bash
youtube-dl --batch-file 5-youtube-urls.txt --extract-audio --audio-format mp3 --output "%(autonumber)s%(title)s.%(ext)s
```

And furthermore, you can combine them into one giant mp3 with this:

```bash
cat *.mp3 > elvis.mp3
```
