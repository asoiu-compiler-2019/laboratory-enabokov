package bnf

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Alternatives struct {
	Alternative []map[string]string `yaml:"alternatives"`
}

type BNF struct {
	Basics             map[string]string `yaml:"basics"`
	Definitions        map[string]string `yaml:"definitions"`
	Primitives         map[string]string `yaml:"primitives"`
	VariableIdentifier string            `yaml:"variableIdentifier"`
	PossibleType       []string          `yaml:"possibleType"`
	Keywords           []string          `yaml:"keywords"`
	Punctuation        []string          `yaml:"punctuation"`
	IntegerLiterals    Alternatives      `yaml:"integerLiterals"`
	FloatLiterals      Alternatives      `yaml:"floatLiterals"`
	StringLiterals     Alternatives      `yaml:"stringLiterals"`
	ArrayLiterals      Alternatives      `yaml:"arrayLiterals"`
	StructType         Alternatives      `yaml:"structTypes"`
	FunctionType       Alternatives      `yaml:"functionTypes"`
}

func (config *BNF) fill(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Failed to read bnf.yml file")
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func Read() (config BNF) {
	config.fill(os.Getenv("BNF_FILE_PATH"))
	return config
}
