package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

type node struct {
	number  int
	faction string
}

func (n *node) setNumber(number int) {
	n.number = number
}

func (n *node) getNumber() int {
	return n.number
}

func (n *node) setFaction(faction string) {
	n.faction = faction
}

func (n *node) getFaction() string {
	return n.faction
}

type regions struct {
	data      []region
	contested int
	region    map[string]int
	keys      []string
}

func (r *regions) createKeys() {
	for k := range r.region {
		r.keys = append(r.keys, k)
	}
	
	sort.Strings(r.keys)
	
	// sort.SliceStable(r.keys, func(i, j int) bool {
	// 	return r.region[r.keys[i]] < r.region[r.keys[j]]
	// })
}

func (r *regions) append(d region) {
	r.data = append(r.data, d)
	if len(d.getFactions()) > 1 {
		r.contested++
		return
	} else if len(d.getFactions()) == 1 {
		if r.region == nil {
			r.region = make(map[string]int)
		}
		r.region[d.getFactions()[0]]++
	}
}


type region struct {
	reg      int
	factions []string
}

func (r *region) setReg(reg int) {
	r.reg = reg
}

func (r *region) getReg() int {
	return r.reg
}

func (r *region) addFaction(faction string) {
	r.factions = append(r.factions, faction)
}

func (r *region) getFactions() []string {
	return r.factions
}

func (r *region) distinguishFactions() {
	k := make(map[string]bool)
	factions := []string{}

	for _, f := range r.factions {
		if k[f] == false {
			k[f] = true
			factions = append(factions, f)
		}
	}

	r.factions = factions
}

type Cases []Case

type Case struct {
	N          int
	M          int
	Letters    []string
	CaseNumber int
	Count      int
	RegNumbers []int
	Regions    regions
}

func (c *Case) GenerateNode(f byte, n int) node {

	iNode := node{}
	iNode.setNumber(n)
	if string(f) != "." {
		iNode.setFaction(string(f))
	}

	return iNode
}

func (c *Case) AddCount() {
	c.Count++
}

func (c *Case) CheckSize(i, j int) bool {
	if c.N == 1 && c.M == 1 {

		if string(c.Letters[i][j]) == "#" {

			c.WriteOutputToFile()
			return true
		} else if string(c.Letters[i][j]) != "#" {

			if string(c.Letters[i][j]) != "." {

				iRegion := region{}
				iRegion.setReg(1)
				iRegion.addFaction(string(c.Letters[i][j]))

				c.Regions.append(iRegion)
			}

			c.Regions.createKeys()
			c.WriteOutputToFile()
			return true
		}

		return true
	}

	return false
}

func (c *Case) CheckMapIfLeftRightSame(left, right byte, strLeft, strRight string, i int, k map[string]interface{}) {

	if k[strLeft] == nil && k[strRight] == nil {

		c.AddCount()

		iNode1 := c.GenerateNode(left, c.Count)
		k[strLeft] = &iNode1

		iNode2 := c.GenerateNode(right, c.Count)
		k[strRight] = &iNode2
	} else if k[strLeft] == nil && k[strRight] != nil {

		iNode := c.GenerateNode(left, k[strRight].(*node).number)
		k[strLeft] = &iNode
	} else if k[strLeft] != nil && k[strRight] == nil {

		iNode := c.GenerateNode(right, k[strLeft].(*node).number)
		k[strRight] = &iNode
	} else

	// REFACTOR ONE OF TWO REGION NUMBER WHEN THERE IS SAME AREA BUT HAVE DIFFERENCE REGION NUMBER
	if k[strLeft] != nil && k[strRight] != nil {

		c.Refactor(strLeft, strRight, i, k)
	}
}

func (c *Case) CheckMapIfLeftRightNotSame(left byte, strLeft string, k map[string]interface{}) {

	if k[strLeft] == nil {

		c.AddCount()

		iNode := c.GenerateNode(left, c.Count)
		k[strLeft] = &iNode
	}
}

func (c *Case) CheckMapIfLeftBelowSame(left, below byte, strLeft, strBelow string, k map[string]interface{}) {

	if k[strLeft] == nil && k[strBelow] == nil {

		c.AddCount()

		iNode1 := c.GenerateNode(left, c.Count)
		k[strLeft] = &iNode1

		iNode2 := c.GenerateNode(below, c.Count)
		k[strBelow] = &iNode2
	} else if k[strLeft] != nil && k[strBelow] == nil {

		iNode := c.GenerateNode(below, k[strLeft].(*node).number)
		k[strBelow] = &iNode
	}
}

func (c *Case) CheckMapIfLeftBelowNotSame(left byte, strLeft string, k map[string]interface{}) {

	if k[strLeft] == nil {

		c.AddCount()

		iNode := c.GenerateNode(left, c.Count)
		k[strLeft] = &iNode
	}
}

func (c *Case) Refactor(strLeft, strRight string, i int, k map[string]interface{}) {

	tmp := k[strRight].(*node).number
	for p := 0; p <= i; p++ {

		for q := 0; q < len(c.Letters[p]); q++ {

			strCurrent := fmt.Sprintf("A[%v][%v]", p, q)
			if k[strCurrent] == nil {
				continue
			}

			if k[strCurrent].(*node).number == tmp {
				k[strCurrent].(*node).setNumber(k[strLeft].(*node).number)
			}
		}
	}
}

func (c *Case) Calculate() {
	k := make(map[string]interface{})

	// Loop whole row
	for i := 0; i < c.N; i++ {

		// Loop whole column of row
		for j := 0; j < c.M; j++ {

			// Don't calculate if the size is 1x1
			if c.CheckSize(i, j) {
				return
			}

			left := c.Letters[i][j]
			strLeft := fmt.Sprintf("A[%v][%v]", i, j)

			// COMPARING CAN BE DONE UNTIL THE COLUMN IS ONE BEFORE THE LAST
			if j < c.M-1 {

				right := c.Letters[i][j+1]
				strRight := fmt.Sprintf("A[%v][%v]", i, j+1)

				// DON'T COMPARE AND CONTINUE WHEN LEFT IS #
				if string(left) == "#" {

					if string(right) != "#" {

						// Fill right to map when right is not #
						if i == c.N-1 && j == c.M-2 && k[strRight] == nil {

							c.AddCount()

							iNode := c.GenerateNode(right, c.Count)
							k[strRight] = &iNode
						}
					}
					continue
				}

				// COMPARE LEFT & RIGHT
				// If left and right are the same
				if string(right) != "#" {

					c.CheckMapIfLeftRightSame(left, right, strLeft, strRight, i, k)
				}

				// COMPARE LEFT & RIGHT
				// If left and right are not the same
				if string(right) == "#" {

					c.CheckMapIfLeftRightNotSame(left, strLeft, k)
				}
			}

			// COMPARING CAN BE DONE UNTIL THE ROW IS ONE BEFORE THE LAST
			if i < c.N-1 {

				below := c.Letters[i+1][j]
				strBelow := fmt.Sprintf("A[%v][%v]", i+1, j)

				// DON'T COMPARE AND CONTINUE WHEN LEFT IS #
				if string(left) == "#" {

					// SET THE BELOW WHEN BELOW IS NOT #
					if string(below) != "#" {

						// Fill below to map when below is not #
						if i == c.N-2 && j == c.M-1 && k[strBelow] == nil {

							c.AddCount()

							iNode := c.GenerateNode(below, c.Count)
							k[strBelow] = &iNode
						}
					}

					continue
				}

				// COMPARE LEFT & BELOW
				// If left and below are the same
				if string(below) != "#" {

					c.CheckMapIfLeftBelowSame(left, below, strLeft, strBelow, k)
				}

				// COMPARE LEFT & BELOW
				// If left and below are not the same
				if string(below) == "#" {

					c.CheckMapIfLeftBelowNotSame(left, strLeft, k)
				}
			}
		}
	}

	// GET ALL NUMBER OF EACH REGION GOTTEN
	// Put each of region number into int slice but there are duplicate region number in int slice so we just need only one of each number
	key := make(map[int]bool)
	for _, j := range k {
		if key[j.(*node).number] == false {
			key[j.(*node).number] = true
			c.RegNumbers = append(c.RegNumbers, j.(*node).number)
		}
	}

	// CALCULATE THE NUMBER OF TERRITORIES OF EACH FACTION AND THE NUMBER OF CONTESTED
	for _, n := range c.RegNumbers {
		iRegion := region{}
		iRegion.setReg(n)
		for _, v := range k {

			if n == v.(*node).number {

				if len(v.(*node).getFaction()) > 0 {
					iRegion.addFaction(v.(*node).getFaction())
				}
			}
		}

		if len(iRegion.getFactions()) > 0 {
			iRegion.distinguishFactions()
			c.Regions.append(iRegion)
		}
	}

	// CREATE KEYS TO HELP LATER TO SHOW REGIONS ORDERED BY KEY AND VALUE
	c.Regions.createKeys()

	// DISPLAY OUTPUT
	c.WriteOutputToFile()
}

func (c *Case) WriteOutputToFile() {

	if _, err := os.Stat("output.in"); err == nil {
		f, err := os.OpenFile("output.in", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		r := fmt.Sprintf("Case %v:\n", c.CaseNumber)
		_, err = f.WriteString(r)
		if err != nil {
			panic(err)
		}

		for _, v := range c.Regions.keys {
			s := fmt.Sprintf("%v %v\n", v, c.Regions.region[v])
			_, err = f.WriteString(s)
			if err != nil {
				panic(err)
			}
		}

		if len(c.Regions.keys) != 0 || c.Regions.contested != 0 {
			t := fmt.Sprintf("contested %d\n", c.Regions.contested)
			_, err = f.WriteString(t)
			if err != nil {
				panic(err)
			}
		}
	} else if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create("output.in")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		r := fmt.Sprintf("Case %v:\n", c.CaseNumber)
		_, err = f.WriteString(r)
		if err != nil {
			panic(err)
		}

		for _, v := range c.Regions.keys {
			s := fmt.Sprintf("%v %v\n", v, c.Regions.region[v])
			_, err = f.WriteString(s)
			if err != nil {
				panic(err)
			}
		}

		if len(c.Regions.keys) != 0 || c.Regions.contested != 0 {
			t := fmt.Sprintf("contested %d\n", c.Regions.contested)
			_, err = f.WriteString(t)
			if err != nil {
				panic(err)
			}
		}
	}

}

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

		caseInstance.CaseNumber = i + 1

		cases = append(cases, caseInstance)
	}

	for i := 0; i < len(cases); i++ {
		cases[i].Calculate()
	}
}
