package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Cases []Case

func main() {
	var cases Cases

	f, err := os.Open("input.in")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// scanner := bufio.NewScanner(os.Stdin)	// <-- Use this if want to read input from console

	scanner.Scan()
	T, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < T; i++ {
		caseInstance := Case{}

		scanner.Scan()
		caseInstance.N, _ = strconv.Atoi(scanner.Text())

		scanner.Scan()
		caseInstance.M, _ = strconv.Atoi(scanner.Text())

		for i := 0; i < caseInstance.N; i++ {
			scanner.Scan()
			caseInstance.Letters = append(caseInstance.Letters, scanner.Text())
		}

		scanner.Scan()
		caseInstance.W = scanner.Text()

		caseInstance.ReverseW()

		caseInstance.CaseNumber = i + 1

		cases = append(cases, caseInstance)
	}

	for i := 0; i < len(cases); i++ {
		cases[i].SearchWord()
	}
}

type Case struct {
	N          int
	M          int
	Letters    []string
	W          string
	Wi         string
	CaseNumber int
	X          int
}

func (c *Case) WriteOutputToFile() {

	if _, err := os.Stat("output.in"); err == nil {
		f, err := os.OpenFile("output.in", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		s := fmt.Sprintf("Case %d: %d\n", c.CaseNumber, c.X)
		_, err = f.WriteString(s)
		if err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create("output.in")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		s := fmt.Sprintf("Case %d: %d\n", c.CaseNumber, c.X)
		_, err = f.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (c *Case) SearchWord() {
	for i := 0; i < c.N; i++ {
		input := c.Letters[i]

		for j := 0; j < c.M; j++ {
			current := input[j]

			// Start checking when current character
			// is same as first character of W
			if current == c.W[0] {
				c.Check(i, j)
			}

			// Start checking when current character
			// is same as first character of W inverse
			if current == c.Wi[0] {
				c.CheckInverse(i, j)
			}
		}
	}
	c.WriteOutputToFile()

	fmt.Printf("Case %d: %d\n", c.CaseNumber, c.X)
}

func (c *Case) Check(i, j int) {
	// Check to right
	limitRight := c.M - len(c.W) + 1
	if j < limitRight {
		c.CheckToRight(i, j)
	}

	// Check to bottom
	limitBottom := c.N - len(c.W) + 1
	if i < limitBottom {
		c.CheckToBottom(i, j)
	}

	// Check to right bottom
	if i < limitBottom && j < limitRight {
		c.CheckToRightBottom(i, j)
	}

	// Check to left bottom
	if i < limitBottom && j >= len(c.W)-1 {
		c.CheckToLeftBottom(i, j)
	}
}

func (c *Case) CheckInverse(i, j int) {
	// Check to right
	limitRight := c.M - len(c.Wi) + 1
	if j < limitRight {
		c.CheckToRightInverse(i, j)
	}

	// Check to bottom
	limitBottom := c.N - len(c.W) + 1
	if i < limitBottom {
		c.CheckToBottomInverse(i, j)
	}

	// Check to right bottom
	if i < limitBottom && j < limitRight {
		c.CheckToRightBottomInverse(i, j)
	}

	// Check to left bottom
	if i < limitBottom && j >= len(c.W)-1 {
		c.CheckToLeftBottomInverse(i, j)
	}
}

func (c *Case) AddX() {
	c.X++
}

func (c *Case) CheckThenAdd(k int) {
	if k == len(c.W)-1 {
		c.AddX()
	}
}

func (c *Case) ReverseW() {
	for i := len(c.W) - 1; i >= 0; i-- {
		c.Wi += string(c.W[i])
	}
}

func (c *Case) CheckToRight(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i][j+k] != c.W[k] {
			break
		}

		// Add X when whole is same as W
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToBottom(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i+k][j] != c.W[k] {
			break
		}

		// Add X when whole is same as W
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToRightBottom(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i+k][j+k] != c.W[k] {
			break
		}

		// Add X when whole is same as W
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToLeftBottom(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i+k][j-k] != c.W[k] {
			break
		}

		// Add X when whole is same as W
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToRightInverse(i, j int) {
	for k := 1; k < len(c.Wi); k++ {
		if c.Letters[i][j+k] != c.Wi[k] {
			break
		}

		// Add X when whole is same as W inverse
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToBottomInverse(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i+k][j] != c.Wi[k] {
			break
		}

		// Add X when whole is same as W inverse
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToRightBottomInverse(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i+k][j+k] != c.Wi[k] {
			break
		}

		// Add X when whole is same as W inverse
		c.CheckThenAdd(k)
	}
}

func (c *Case) CheckToLeftBottomInverse(i, j int) {
	for k := 1; k < len(c.W); k++ {
		if c.Letters[i+k][j-k] != c.Wi[k] {
			break
		}

		// Add X when whole is same as W inverse
		c.CheckThenAdd(k)
	}
}
