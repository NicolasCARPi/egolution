/* egolution.go */
package main

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"
)

var genome string

// radiation level
const rad = 20

// number of iterations
const iter = 500000

// possible letters in the genome
const letters = " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.!?:,;"

const badLetters = "kqwxyz"

const genomeSize = 800

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	genome := getLetters(genomeSize)
	score := getScore(genome)
	fmt.Printf("[Time 0] %s (%2.2f)\n", genome, score)
	for i := 0; i < iter; i++ {
		virus := getLetters(rad)
		mutated := mutate(genome, virus)
		tryScore := getScore(mutated)

		if tryScore > score {
			// keep this mutation
			genome = mutated
			score = tryScore
			fmt.Print("+")
		}

		if i%25 == 0 {
			fmt.Print(".")
		}

	}
	fmt.Printf("\n[Time %d] %s (%2.1f)\n", iter, genome, score)
}

func getLetters(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func mutate(genome string, virus string) string {
	b := []rune(genome)
	v := []rune(virus)
	insertPos := randInt(0, len(genome)-len(virus))
	for i := 0; i < rad; i++ {
		b[insertPos+i] = v[i]
	}
	return string(b)
}

func getScore(genome string) float64 {
	bad := 0.0
	good := 0.0
	r := []rune(genome)

	// walk the genome
	for i := 0; i < len(genome); i++ {

		// count the space
		if unicode.IsSpace(r[i]) {
			good += 1
		}
		// numbers are not wanted
		if unicode.IsNumber(r[i]) {
			bad += 1
		}
		// we don't want punctuation or uppercase
		if unicode.In(r[i], unicode.Punct) || unicode.In(r[i], unicode.Upper) {
			bad += 1
		}
		// ideal number of space is every 5 letters
		ideal := float64(len(genome)) / 5.0
		if int(good) == int(ideal) {
			good += 1
		}

	}
	return (good / (bad + 1)) / genomeSize * 100000
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
