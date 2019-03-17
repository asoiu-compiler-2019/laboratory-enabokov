package lexis

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/enabokov/language/bnf"
)

var bnfConfig bnf.BNF

func init() {
	bnfConfig = bnf.Read()
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

func Analyze(filename string) []Token {
	lines := readLines(filename)
	stream := readInputStream(lines)
	tokens := readTokenStream(stream)

	for i := 0; i < 500; i++ {
		token := tokens.next()
		if token != nil {
			fmt.Println(*token)
		}
	}

	return nil
}
