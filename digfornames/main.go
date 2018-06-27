package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jdkato/prose/chunk"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
)

var findQuotedString = regexp.MustCompile(`\[\[([^\]]*)\]\]`)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const width = 40

func renderInContext(choice int, term, sentence, song string) {
	if strings.Index(song, term) > -1 {
		return
	}
	termStart := strings.Index(sentence, term)
	termEnd := termStart + len(term)
	songStart := strings.Index(sentence, song)
	songEnd := songStart + len(song)
	if termStart > -1 {
		start := max(termStart-width, 0)
		end := min(termEnd+width, len(sentence))
		parts := []string{fmt.Sprintf("%2v => ", choice)}
		for i := start; i < end; i++ {
			letter := sentence[i : i+1]
			switch {
			case i >= termStart && i <= termEnd:
				letter = color.RedString(letter)
				break
			case i >= songStart && i <= songEnd:
				letter = color.BlueString(letter)
				break
			}
			parts = append(parts, letter)
		}
		fmt.Fprintln(os.Stderr, strings.Join(parts, ""))
	}
}

func errOut(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	defer os.Exit(1)
}

func main() {
	color.NoColor = false
	buf := bufio.NewReader(os.Stdin)

	inFile, _ := os.Open(os.Args[1])
	if len(os.Args) < 1 {
		errOut("usage: digfornames <infile>")
		return
	}

	startAt := 0
	if len(os.Args) > 2 {
		start, _ := strconv.Atoi(os.Args[2])
		startAt = start - 1
	}

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	numLines := -1
	for scanner.Scan() {
		numLines++
		sentence := scanner.Text()
		if numLines < startAt {
			continue
		}
		match := findQuotedString.FindAllStringSubmatch(sentence, -1)
		if len(match) < 1 {
			errOut("no match found: " + sentence)
			return
		}
		song := match[0][1]
		fmt.Fprintln(os.Stderr, color.BlueString(song))
		words := tokenize.TextToWords(sentence)
		regex := chunk.TreebankNamedEntities

		tagger := tag.NewPerceptronTagger()
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%2v => not found", 0))
		var entities []string
		for i, entity := range chunk.Chunk(tagger.Tag(words), regex) {
			renderInContext(i+1, entity, sentence, song)
			entities = append(entities, entity)
		}
		choice, _ := buf.ReadString('\n')
		choice = strings.TrimFunc(choice, func(r rune) bool {
			return r == '\n'
		})
		choiceNum, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if choiceNum > 0 {
			fmt.Fprintln(os.Stdout, song, entities[choiceNum-1])
		}
	}
}
