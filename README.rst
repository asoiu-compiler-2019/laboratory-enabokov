=================
Language compiler
=================

*Written in purpose of investigating how any language works under the hood*

**Grammar**
===========
*BNF is written in format of YAML config file. It's easy to use within the program.
Every language has own implementation of YAML marshalling/unmarshalling, that enables this feature.*

- **basics** - defines the primitives of BNF (e.g alter -> alternation, repete -> repetition)

- **definitions** - explains unicode char

- **primitives** - points basic components of primitive types

- **variableIdentifier** - grammar for building lexemes

- **types** - possible types in the new language

- **keywords** - possible keywords in the new language

- **operators** - possible operators in the new language

- **punctuation** - possible puncts in the new language

- **comments** - symbol that enforces analyzers to skip certain string

- **integerLiterals** - structure of Integer literals

  - **alternatives** - used to define alternative way for Integer numbers

- **floatLiterals** - structure of Float literals

  - **alternatives** - same as for previous one

- **stringLiterals** - structure of String literals

  - **alternatives** - same as for previous one

- **arrayLiterals** - structure of Array literal

  - **alternatives** - same as for previous one

- **structTypes** - structure of Struct literal

  - **alternatives** - same as for previous one

- **functionTypes** - structure of Function literal

  - **alternatives** - same as for previous one

- **bnf.go** - parse grammar rules and applies for parsing

**Lexis analyzer**
==================
*Responsible for converting each primitive or group of them into lexical token from the given source code*

.. code-block::

  func main() {} -> []string{"f", "u", "n", "c", " ", "m", "a", "i", "n", " ", "{", "}"} -> []tokens{"func", "main", "(", ")"}

- **check.go** - responsible for clarifying belonging of the token, whether identifier, variable, number, string, etc.

- **reader.go** - defines naming classes and reads chars from source code grouping them into tokens

- **stream.go** - defines stateful struct

  - **Peek()** - returns next token, but do NOT step forward

  - **Next()** - returns next token, shifting forward

  - **EOF()** - checks, whether there is an end of file

  - **Croak()** - notify about errors including line and column information

- **reader.go** - entrypoint

**Syntax analyzer**
===================
*Responsible for grouping lexical tokens into Abstract Syntax Tree (AST)*

Source code
^^^^^^^^^^^

.. code-block::

   package main

   func main(a, b) {}

Abstract syntax tree
^^^^^^^^^^^^^^^^^^^^

.. code-block::

   syntax.TokenProgram{
      Class:      "program",
      Expression: {
         syntax.TokenPackage{Class:"package", Value:"main"},
         syntax.TokenFunction{
            Class:  "function",
            Name:   "main",
            Params: {
               {Class:"variable", Name:"a", Type:""},
               {Class:"variable", Name:"b", Type:""},
            },
            Body: nil,
         },
      },
   }

   2019/04/02 04:11:52 OK

- **check.go** - check what structure group of token describes (e.g contional operator, assignment, function) 

- **defaults.go** - defines tokens structures, that are used within *syntax* package and *semantics* as well.

- **parser.go** - after check the groups of tokens -> *parser* starts grabbing all tokens until handles logical component.

- **syntax.go** - entrypoint

**Semantics analyzer**
======================
*Responsible for validating the groups of tokens and their meaningful in the code*

Source code
^^^^^^^^^^^

.. code-block::

   package test

   import "fmt"

   def do(a, b) {
      var a int
      var c int

      fmt.Println(a, b, c)
   }

Validation
^^^^^^^^^^

.. code-block::

   2019/04/02 04:30:23 Variable `a` is already defined in `do`

- **walk.go** - traverse Abstract syntax tree and looks for ambigious situations. If found -> trigger error notifier

- **semantics.go** - entrypoint

Main
====
*Start entrypoint for compiling source code*

.. code-block::

   go build compile.go
   ./compile examples/complex.ena

- **compile.go** - accepts path to source code `*.ena` for further parsing

Dependencies
============
*This project depends on some packages*

go.mod
^^^^^^

.. code-block::

   // list deps packages
   require (
      github.com/kr/pretty v0.1.0
      gopkg.in/yaml.v2 v2.2.2
   )

Examples
========
*Defined 3 examples of source code*

- **simple.ena**

- **duplicated_vars.ena**

- **complex.ena**
