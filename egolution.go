/* egolution.go */
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
	"unicode"
)

var genome string

// radiation level
const rad = 2

// number of iterations
const defaultIter = 500000

// possible letters in the genome
const letters = " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.!?:,;"

const badLetters = "kqwxyz"

const genomeSize = 45

func main() {
	iter := flag.Int("i", defaultIter, "Number of iterations to run")
	quiet := flag.Bool("q", false, "Suppress output")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	genome := getLetters(genomeSize)
	score := getScore(genome)
	firstLine := fmt.Sprintf("0\t%s\t%2.2f\t%d\n", genome[:20], score, len(genome))
	for i := 0; i < *iter; i++ {
		mutated := mutate(genome)
		tryScore := getScore(mutated)

		if tryScore > score {
			// keep this mutation
			genome = mutated
			score = tryScore
			if !*quiet {
				fmt.Print("+")
			}
		}

		if !*quiet {
			if i%25 == 0 {
				fmt.Print(".")
			}
		}

	}
	if !*quiet {
		fmt.Print("\n")
	}
	fmt.Printf("Time\tGenome\tScore\tLength\n")
	fmt.Println(firstLine)
	fmt.Printf("%d\t%s\t%2.1f\t%d\n", *iter, genome[:20], score, len(genome))
}

func getLetters(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// a mutation can be infect, add or loss
func mutate(genome string) string {
	dice := randInt(0, 100)

	if dice%15 == 0 {
		return loss(genome)

	} else if dice%10 == 0 {
		return add(genome)

	} else {
		return infect(genome)
	}
}

// add a virus of length rad to the genome
func infect(genome string) string {
	g := []rune(genome)
	v := []rune(getLetters(rad))
	insertPos := randInt(0, len(genome)-len(v))
	for i := 0; i < rad; i++ {
		g[insertPos+i] = v[i]
	}
	return string(g)
}

// add a single nucleotide
func add(genome string) string {
	return genome + getLetters(1)
}

// remove a single nucleotide
func loss(genome string) string {
	g := []rune(genome)
	deletePos := randInt(0, len(genome)-1)
	full := append(g[:deletePos], g[deletePos+1:]...)
	return string(full)
}

func getScore(genome string) float64 {
	bad := 1.0
	good := 1.0
	g := []rune(genome)

	// walk the genome
	for i := 0; i < len(genome); i++ {

		// count the space
		if unicode.IsSpace(g[i]) {
			good += 0.1
		}
		// we don't want punctuation, numbers or uppercase
		if unicode.In(g[i], unicode.Punct) ||
			unicode.In(g[i], unicode.Upper) ||
			unicode.IsNumber(g[i]) {
			bad += 1
		}
		// ideal number of space is every 5 letters
		ideal := float64(len(genome)) / 5.0
		if int(good) == int(ideal) {
			good += 1
		}

	}
	return ((good / bad) / genomeSize) * 100
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
