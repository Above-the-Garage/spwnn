/*
Package spwnn - neural spelling corrector worker
*/
package spwnn

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

// 26 letters plus underscore (to mark start AND end of word, or characters not
// in the alphabet like apostraphe)

const sAlphabetSize = 28

// SpwnnDictionary contains all necessary state to analyze a word
type SpwnnDictionary struct {
	wordCount       int
	words           []string
	wordScore       []float64
	lenDiff         []int
	neuronProtect   [sAlphabetSize * sAlphabetSize]bool
	neuralIndex     [sAlphabetSize * sAlphabetSize][]string
	neuralIndexSize [sAlphabetSize * sAlphabetSize]int
}

// GetWordCount returns the total number of words in the dictionary
func GetWordCount(dict *SpwnnDictionary) int {
	return dict.wordCount
}

func addStartEnd(word string) string {
	return "_" + word + "_"
}

// RemoveSpaces removes spaces and new lines and carriage returns that might be in a string
func RemoveSpaces(line string) string {
	line = strings.Replace(line, " ", "", -1) // -1 means 'all'
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, "\r", "", -1)
	return line
}

// GetWords returns the array of words in the dictionary
func GetWords(dict *SpwnnDictionary) []string {
	return dict.words
}

// MakeOneExactly returns 1.0 if the argument is pretty close to 1.0
// otherwise returns the argument
func MakeOneExactly(f float64) float64 {
	if math.Abs(f-1.0) < 0.0000000001 {
		return 1.0
	}
	return f
}

func charToIndex(ch byte) int {
	res := byte(0)
	if ch >= 'a' && ch <= 'z' {
		res = ch - 'a'
	} else if ch >= 'A' && ch <= 'Z' {
		res = ch - 'A'
	} else if ch == '_' {
		res = 26
	} else {
		res = 27
	}
	// assert res < sAlphabetSize
	return int(res)
}

func indexToChar(n int) byte {
	return byte(n) + 'a' // which causes all non-alpha (mapped to 27) are converted to opening brace '{'
}

func computeIndex(n1, n2 int) int {
	return n1*sAlphabetSize + n2
}

func getNeuronIndex(ch1, ch2 byte) int {
	n1 := charToIndex(ch1)
	n2 := charToIndex(ch2)
	return computeIndex(n1, n2)
}

func addNeuron(dict *SpwnnDictionary, ch1, ch2 byte, word string) {
	nIndex := getNeuronIndex(ch1, ch2)
	if !dict.neuronProtect[nIndex] {
		dict.neuralIndex[nIndex] = append(dict.neuralIndex[nIndex], word)
		dict.neuralIndexSize[nIndex]++
		dict.neuronProtect[nIndex] = true
	}
}

func rememberWord(dict *SpwnnDictionary, word string) {
	// Make permanent copy of the word
	dict.words = append(dict.words, word)
	dict.wordScore = append(dict.wordScore, 0.0)
	dict.lenDiff = append(dict.lenDiff, 0)
	dict.wordCount++
}

func clearList(dict *SpwnnDictionary, ch1, ch2 byte) {
	nIndex := getNeuronIndex(ch1, ch2)
	dict.neuralIndex[nIndex] = make([]string, 0)
}

func clearNetwork(dict *SpwnnDictionary) {
	dict.words = make([]string, 0)
	dict.wordCount = 0
	for ch1 := 'a'; ch1 <= 'z'; ch1++ {
		for ch2 := 'a'; ch2 <= 'z'; ch2++ {
			clearList(dict, byte(ch1), byte(ch2))
		}
	}
}

func addWordToNetwork(dict *SpwnnDictionary, word string) {
	rememberWord(dict, word)
	for index := range dict.neuronProtect {
		dict.neuronProtect[index] = false
	}
	length := len(word)
	for i := 0; i < length-1; i++ {
		addNeuron(dict, word[i], word[i+1], word)
	}
}

// ReadDictionary reads "knownWords.txt" into a SpwnnDictionary
func ReadDictionary(noisy bool) *SpwnnDictionary {
	file, err := os.Open("knownWords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var dict SpwnnDictionary

	counter := 0
	clearNetwork(&dict)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var word string
		word = scanner.Text()
		word = RemoveSpaces(word)
		word = addStartEnd(word)
		addWordToNetwork(&dict, word)
		counter++
		if noisy && counter > 1000 {
			fmt.Print(".")
			counter = 0
		}
	}
	sort.Strings(dict.words[:])
	if noisy {
		fmt.Printf("\n%d words in dictionary\n", dict.wordCount)
	}
	return &dict
}

// PrintNeuron prints all of words associated with a letter pair
func PrintNeuron(dict *SpwnnDictionary, input string) {
	input = RemoveSpaces(input)
	if len(input) < 2 {
		fmt.Printf("Bad input for printNeuron ('%s')\n", input)
		return
	}
	ch1 := input[0]
	ch2 := input[1]
	n1 := charToIndex(ch1)
	n2 := charToIndex(ch2)
	nIndex := computeIndex(n1, n2)
	fmt.Println(ch1, ch2, dict.neuralIndex[nIndex])
}

// PrintIndexSizes prints the size of each list associated with a letter pair
func PrintIndexSizes(dict *SpwnnDictionary) {
	for i := 0; i < sAlphabetSize; i++ {
		for j := 0; j < sAlphabetSize; j++ {
			nIndex := computeIndex(i, j)
			size := dict.neuralIndexSize[nIndex]
			if size != 0 {
				fmt.Printf("%c, %c => %d\n", indexToChar(i), indexToChar(j), size)
			}
		}
	}
}

// MaxIndexSize the size of the largest list assocaited with a letter pair
func MaxIndexSize(dict *SpwnnDictionary) int {
	maxSize := 0
	for i := 0; i < sAlphabetSize; i++ {
		for j := 0; j < sAlphabetSize; j++ {
			nIndex := computeIndex(i, j)
			size := dict.neuralIndexSize[nIndex]
			if size > maxSize {
				maxSize = size
			}
		}
	}
	return maxSize
}

func findWord(dict *SpwnnDictionary, word string) int {
	// words must be sorted as per Go rules for this binary search to work
	word = RemoveSpaces(word)
	low := -1
	high := dict.wordCount
	for high-low > 1 {
		var i int
		i = (high + low) / 2
		if word == dict.words[i] {
			return i
		}
		if word < dict.words[i] {
			high = i
		} else {
			low = i
		}
	}
	if word == dict.words[high] {
		return high
	}
	return -1
}

func increaseWordScore(dict *SpwnnDictionary, word string, factor float64, lenDiff int) {
	index := findWord(dict, word)
	dict.wordScore[index] = dict.wordScore[index] + factor
	dict.lenDiff[index] = lenDiff
}

func intMax(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

// SpwnnResult is a word and a corresponding score
type SpwnnResult struct {
	Score   float64
	LenDiff int
	Word    string
}

// ByScore sorts results list by score distance from 1.0
// This must be exported for "Sort" ... ?
type ByScore []SpwnnResult

func (r ByScore) Len() int           { return len(r) }
func (r ByScore) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByScore) Less(i, j int) bool { return math.Abs(r[i].Score-1.0) < math.Abs(r[j].Score-1.0) }

// ByLength sorts results by distance from length of original word
type ByLength []SpwnnResult

//
func (r ByLength) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

//
func (r ByLength) Less(i, j int) bool { return r[i].LenDiff < r[j].LenDiff }

//
func (r ByLength) Len() int { return len(r) }

func intAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// CorrectSpelling finds words similar to given word; also returns the number of words "touched" by the algorithm
func CorrectSpelling(dict *SpwnnDictionary, word string, strictLen bool) ([]SpwnnResult, int) {

	// clear scores if leftover from previous run
	for i := 0; i < dict.wordCount; i++ {
		dict.wordScore[i] = 0.0
		dict.lenDiff[i] = 0
	}

	// Allow words with "_" already attached to beginning and end;
	// if missing, add it as beginning and end markers.
	if word[0] != '_' {
		word = addStartEnd(word)
	}
	// Iterate through all pairs of letters
	for i := 0; i < len(word)-1; i++ {
		// Find list of words for this letter pair
		n1 := charToIndex(word[i])
		n2 := charToIndex(word[i+1])
		index := computeIndex(n1, n2)
		wordList := dict.neuralIndex[index]
		// For each word that has that letter pair,
		// compute contribution of that letter pair.
		var contribution float64
		for _, aWord := range wordList {
			// A human only has access to the word they see,
			// so word pair contributions are based on that.
			// Also in testing it results in far fewer ties.
			contribution = 1.0 / float64(len(word)-1)
			lenDiff := intAbs(len(aWord) - len(word))
			increaseWordScore(dict, aWord, contribution, lenDiff)
		}
	}

	// Find best score received by any word; also count number of words touched
	var bestScore float64
	var wordsTouched int
	for i := 0; i < dict.wordCount; i++ {
		if dict.wordScore[i] != 0.0 {
			dict.wordScore[i] = MakeOneExactly(dict.wordScore[i])
			wordsTouched++
		}
		// Best score is the score closest to 1.0
		if math.Abs(bestScore-1.0) > math.Abs(dict.wordScore[i]-1.0) {
			bestScore = dict.wordScore[i]
		}
	}

	// find near-winners
	var results []SpwnnResult
	results = make([]SpwnnResult, 0)
	for i := 0; i < dict.wordCount; i++ {
		if math.Abs(dict.wordScore[i]-bestScore) == 0.0 {
			if strictLen && dict.lenDiff[i] != 0 {
				continue
			}
			var res SpwnnResult
			res.Score = dict.wordScore[i]
			res.LenDiff = dict.lenDiff[i]
			res.Word = dict.words[i]
			results = append(results, res)
		}
	}

	// sort with best scores first
	sort.Sort(ByScore(results))
	// sort by length, shorter words are better
	sort.Sort(ByLength(results))

	return results, wordsTouched
}
