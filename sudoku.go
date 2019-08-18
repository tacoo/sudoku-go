package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	squares  []string
	unitlist [][]string
	units    map[string][][]string
	peers    map[string][]string
)

const (
	digits = "123456789"
	rows   = "ABCDEFGHI"
	cols   = digits
)

func cross(a, b string) (crossed []string) {
	for _, x := range a {
		for _, y := range b {
			crossed = append(crossed, string(x)+string(y))
		}
	}
	return crossed
}

func init() {
	squares = cross(rows, cols)
	for _, c := range cols {
		unitlist = append(unitlist, cross(rows, string(c)))
	}
	for _, r := range rows {
		unitlist = append(unitlist, cross(string(r), cols))
	}
	for _, a := range []string{"ABC", "DEF", "GHI"} {
		for _, b := range []string{"123", "456", "789"} {
			unitlist = append(unitlist, cross(a, b))
		}
	}
	units = make(map[string][][]string)
	for _, v := range squares {
		var unit [][]string
		for _, u := range unitlist {
			if contains(u, v) {
				unit = append(unit, u)
			}
		}
		units[v] = unit
	}
	peers = make(map[string][]string)
	for _, v := range squares {
		peerSet := make(map[string]bool)
		for _, u := range units[v] {
			for _, uu := range u {
				if uu != v {
					peerSet[uu] = true
				}
			}
		}
		var peer []string
		for k := range peerSet {
			peer = append(peer, k)
		}
		peers[v] = peer
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func assign(values map[string]string, s, d string) {
	remain := strings.Replace(values[s], d, "", -1)
	for _, c := range remain {
		v := string(c)
		eliminate(values, s, v)
	}
}

func eliminate(values map[string]string, s, d string) {
	if strings.Index(values[s], d) == -1 {
		return
	}
	values[s] = strings.Replace(values[s], d, "", -1)
	if len(values[s]) == 0 {
		return
	}
	if len(values[s]) == 1 {
		for _, s2 := range peers[s] {
			eliminate(values, s2, values[s])
		}
	}
	for _, u := range units[s] {
		var dplaces []string
		for _, s1 := range u {
			if strings.Index(values[s1], d) != -1 {
				dplaces = append(dplaces, s1)
			}
		}
		if len(dplaces) == 0 {
			return
		}
		if len(dplaces) == 1 {
			assign(values, dplaces[0], d)
		}
	}
}

func search(values map[string]string) map[string]string {
	solved := true
	for _, s := range squares {
		if len(values[s]) != 1 {
			solved = false
		}
	}
	if solved {
		return values
	}
	var (
		min int
		n   string
	)
	for _, s := range squares {
		size := len(values[s])
		if size > 1 {
			if min == 0 || min > size {
				min = size
				n = s
			}
		}
	}
	var results []map[string]string
	for _, v := range values[n] {
		copied := make(map[string]string)
		for k, v := range values {
			copied[k] = v
		}
		vv := string(v)
		assign(copied, n, vv)
		results = append(results, search(copied))
	}
	for _, x := range results {
		if len(x) > 0 {
			return x
		}
	}
	return nil
}

func solve(grid string) map[string]string {
	values := make(map[string]string)
	input := make(map[string]string)
	for i, v := range squares {
		values[v] = digits
		input[v] = string(grid[i])
	}
	for k, v := range input {
		if strings.Index(digits, v) != -1 {
			assign(values, k, v)
		}
	}
	return search(values)
}

func isSolved(values map[string]string) bool {
	for _, unit := range unitlist {
		var vs []string
		for _, u := range unit {
			vs = append(vs, values[u])
		}
		for _, d := range digits {
			d1 := string(d)
			found := false
			for _, v := range vs {
				if v == d1 {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

const sample = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"

func loadfile(f string) (lines []string) {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func print(samples []string) {
	for _, s := range samples {
		fmt.Println("target:", s)
		vs := solve(s)
		if isSolved(vs) {
			fmt.Print("solved: ")
			for _, v := range squares {
				fmt.Print(vs[v])
			}
			fmt.Println("")
		} else {
			fmt.Println("!!! could not solved")
		}
	}
}

func main() {
	print(loadfile("./src/github.com/tacoo/sudoku-go/easy50.txt"))
	print(loadfile("./src/github.com/tacoo/sudoku-go/top95.txt"))
	print(loadfile("./src/github.com/tacoo/sudoku-go/hardest.txt"))
}
