/* egolution.go */
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var verbose *bool

// number of iterations
const defaultIter = 50000

// possible azoted bases in the genome, we work with RNA here
const bases = "AUGC"

// size in aminoacids
const genomeSize = 10

func main() {

	// parse arguments
	iter := flag.Int("i", defaultIter, "Number of iterations to run")
	verbose = flag.Bool("v", false, "Enable output")
	flag.Parse()

	// init random
	rand.Seed(time.Now().UTC().UnixNano())

	// create our base genome
	genome := getLetters(genomeSize * 3)
	score := len(translateGenome(genome))
	firstLine := fmt.Sprintf("0\t%d\t%d\n", score, len(genome))
	for i := 0; i < *iter; i++ {
		mutated := mutate(genome)
		tryScore := len(translateGenome(mutated))

		// evolution step
		if tryScore > score {
			// keep this mutation
			genome = mutated
			score = tryScore
			if *verbose {
				fmt.Println("✓", score)
			}
		}
	}
	if *verbose {
		fmt.Print("\n")
	}
	fmt.Printf("Time\tProt size\tGenome size\n")
	fmt.Println(firstLine)
	fmt.Printf("%d\t%d\t%d\n", *iter, score, len(genome))
}

func getLetters(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = bases[rand.Intn(len(bases))]
	}
	return string(b)
}

// a mutation can be infect, add or loss
func mutate(genome string) string {
	dice := randInt(0, 100)

	if dice%15 == 0 {
		if *verbose {
			fmt.Print("-")
		}
		return loss(genome)

	} else if dice%10 == 0 {
		if *verbose {
			fmt.Print("*")
		}
		return infect(genome)

	} else {
		return add(genome)
	}
}

// replace part of the genome with a virus
func infect(genome string) string {
	g := []byte(genome)
	virusSize := randInt(1, 4)
	v := []byte(getLetters(virusSize))
	insertPos := randInt(0, len(genome)-virusSize)
	for i := 0; i < virusSize; i++ {
		g[insertPos+i] = v[i]
	}
	return string(g)
}

// add some nucleotide
func add(genome string) string {
	if *verbose {
		fmt.Print("+")
	}
	return genome + getLetters(randInt(1, 500))
}

// remove a single nucleotide
func loss(genome string) string {
	g := []byte(genome)
	deletePos := randInt(0, len(genome)-1)
	return string(append(g[:deletePos], g[deletePos+1:]...))
}

func translateGenome(genome string) string {
	var protein string
	// the RNA is read every 3 letters
	for i := 0; i < len(genome)-3; i += 3 {
		aa := translate(genome[i : i+3])
		// stop codon is bad
		if aa == "!" {
			return protein
		}
		protein = protein + aa
	}
	return protein
}

func translate(codon string) string {
	aminoAcids := map[string]string{
		// ! is for stop codon
		"UAG": "!",
		"UAA": "!",
		"UGA": "!",
		"GCU": "A",
		"GCC": "A",
		"GCA": "A",
		"GCG": "A",
		"ACU": "T",
		"ACC": "T",
		"ACA": "T",
		"ACG": "T",
		"CCU": "P",
		"CCC": "P",
		"CCA": "P",
		"CCG": "P",
		"UCU": "S",
		"UCC": "S",
		"UCA": "S",
		"UCG": "S",
		"UUU": "F",
		"UUC": "F",
		"UUA": "L",
		"UUG": "L",
		"CUU": "L",
		"CUG": "L",
		"CUA": "L",
		"CUC": "L",
		"AUU": "I",
		"AUC": "I",
		"AUA": "I",
		"AUG": "M",
		"GUU": "V",
		"GUC": "V",
		"GUA": "V",
		"GUG": "V",
		"UAU": "Y",
		"UAC": "Y",
		"CAU": "H",
		"CAC": "H",
		"CAA": "Q",
		"CAG": "Q",
		"AAU": "N",
		"AAC": "N",
		"AAA": "K",
		"AAG": "K",
		"GAU": "D",
		"GAC": "D",
		"GAA": "E",
		"GAG": "E",
		"UGU": "C",
		"UGC": "C",
		"UGG": "W",
		"CGU": "R",
		"CGC": "R",
		"CGA": "R",
		"CGG": "R",
		"AGU": "S",
		"AGC": "S",
		"AGA": "R",
		"AGG": "R",
		"GGU": "G",
		"GGC": "G",
		"GGA": "G",
		"GGG": "G",
	}

	return aminoAcids[codon]
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
