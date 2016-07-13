/* egolution.go */
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// radiation level
const rad = 3

// number of iterations
const defaultIter = 50000

// possible azoted bases in the genome, we work with RNA here
const bases = "AUGC"

// list of amino acids
const aa = "ARNDCEQGHILKMFPSTWYV"

// size in aminoacids
const genomeSize = 10

func main() {
	iter := flag.Int("i", defaultIter, "Number of iterations to run")
	quiet := flag.Bool("q", false, "Suppress output")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	// create our base genome
	genome := getLetters(genomeSize * 3)
	score := getScore(genome)
	firstLine := fmt.Sprintf("0\t%s\t%2.2f\t%d\n", genome, score, len(genome))
	for i := 0; i < *iter; i++ {
		mutated := mutate(genome)
		tryScore := getScore(translateGenome(mutated))

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
	fmt.Printf("%d\t%s\t%2.1f\t%d\n", *iter, genome, score, len(genome))
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
		return loss(genome)

	} else if dice%10 == 0 {
		return add(genome)

	} else {
		return infect(genome)
	}
}

// add a virus of length rad to the genome
func infect(genome string) string {
	g := []byte(genome)
	v := []byte(getLetters(3 * rad))
	insertPos := randInt(0, len(genome)-len(v))
	for i := 0; i < rad*3; i++ {
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
	g := []byte(genome)
	deletePos := randInt(0, len(genome)-1)
	return string(append(g[:deletePos], g[deletePos+1:]...))
}

func translateGenome(genome string) string {
	var protein string
	// the RNA is read every 3 letters
	for i := 0; i < len(genome)-3; i += 3 {
		codon := genome[i : i+3]
		protein = protein + translate(codon)
	}
	return protein
}

// assess the quality of a protein
func getScore(protein string) float64 {
	score := 1.0

	for _, a := range protein {
		if a == 70 || a == 73 || a == 75 {
			score += 2
		}
		// stop codon is bad
		if a == 33 {
			return 0
		}
	}

	return score
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
