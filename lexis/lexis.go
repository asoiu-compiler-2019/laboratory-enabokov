package lexis

import (
	"bufio"
	"log"
	"os"

	"github.com/enabokov/language/bnf"
)

var BnfConfig bnf.BNF

func init() {
	BnfConfig = bnf.Read()
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}

	return lines
}

func Analyze(filename string) TokenStream {
	lines := readLines(filename)
	stream := readInputStream(lines)
	return readTokenStream(stream)
}
