// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

import (
	"math"
	"strings"
	"unicode"
	"github.com/mattn/go-runewidth"
)

var (
	nl = "\n"
	sp = " "
)

const defaultPenalty = 1e5

// Wrap wraps s into a paragraph of lines of length lim, with minimal
// raggedness.
func WrapString(s string, lim int) ([]string, int) {
	var strs []string
	for _, r := range s {
		n := len(strs)
		if unicode.Is(unicode.Han, r) {
			if n == 0 || strs[n-1] != "" {
				strs = append(strs, string(r))
			} else {
				strs[n-1] = string(r)
			}
			strs = append(strs, "")
		} else {
			if n == 0  {
				strs = append(strs, string(r))
			} else {
				strs[n-1] = strs[n-1] + string(r)
			}
		}
	}
	if len(strs) > 0 && strs[len(strs)-1] == "" {
		strs = strs[:len(strs)-1]
	}

	var words []string
	for _, str := range strs {
		words = append(words, (strings.Split(strings.Replace(str, nl, sp, -1), sp))...)
	}
	
	var lines []string
	max := 0
	for _, v := range words {
		max = runewidth.StringWidth(v)
		if max > lim {
			lim = max
		}
	}
	for _, line := range WrapWords(words, 1, lim, defaultPenalty) {

		for i, word := range line {
			wordRune := []rune(word)
			n := len(wordRune)
			if n > 0 && !unicode.Is(unicode.Han, wordRune[n-1]) {
				line[i] = word + sp
			}
			if n == 0 {
				line[i] = sp
			}
		}

		lines = append(lines, strings.Join(line, ""))
	}
	return lines, lim
}

// WrapWords is the low-level line-breaking algorithm, useful if you need more
// control over the details of the text wrapping process. For most uses,
// WrapString will be sufficient and more convenient.
//
// WrapWords splits a list of words into lines with minimal "raggedness",
// treating each rune as one unit, accounting for spc units between adjacent
// words on each line, and attempting to limit lines to lim units. Raggedness
// is the total error over all lines, where error is the square of the
// difference of the length of the line and lim. Too-long lines (which only
// happen when a single word is longer than lim units) have pen penalty units
// added to the error.
func WrapWords(words []string, spc, lim, pen int) [][]string {
	n := len(words)

	length := make([][]int, n)
	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = runewidth.StringWidth(words[i])
		for j := i + 1; j < n; j++ {
			length[i][j] = length[i][j-1] + spc + runewidth.StringWidth(words[j])
		}
	}
	nbrk := make([]int, n)
	cost := make([]int, n)
	for i := range cost {
		cost[i] = math.MaxInt32
	}
	for i := n - 1; i >= 0; i-- {
		if length[i][n-1] <= lim {
			cost[i] = 0
			nbrk[i] = n
		} else {
			for j := i + 1; j < n; j++ {
				d := lim - length[i][j-1]
				c := d*d + cost[j]
				if length[i][j-1] > lim {
					c += pen // too-long lines get a worse penalty
				}
				if c < cost[i] {
					cost[i] = c
					nbrk[i] = j
				}
			}
		}
	}
	var lines [][]string
	i := 0
	for i < n {
		lines = append(lines, words[i:nbrk[i]])
		i = nbrk[i]
	}
	return lines
}

// getLines decomposes a multiline string into a slice of strings.
func getLines(s string) []string {
	return strings.Split(s, nl)
}
