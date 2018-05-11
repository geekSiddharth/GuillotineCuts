// Plane permutations are those which avoid the pattern 213'54
// Avoiding 213'54 is equivalent to avoiding 2-14-3
// This generating tree based enumeration is based on the following paper:
// https://arxiv.org/pdf/1702.04529.pdf

package main

import (
    "fmt"
    "runtime"
    "sync"
)

// struct to mimick set in Python | NOT FUNCTIONAL
// code from https://play.golang.org/p/_FvECoFvhq
type SliceSet struct {
	set map[[]int]bool
}

func NewSliceSet() *SliceSet {
	return &SliceSet{make(map[[]int]bool)}
}

func (set *SliceSet) Add(p []int) bool {
	_, found := set.set[p]
	set.set[p] = true
	return !found	//False if it existed already
}

// func to mimick min in Python
func min (a, b int) (int) {
	if a<=b {return a}
	return b
}

func localExp (perm []int, a int, p chan []int, c chan int) {
	c <- 1
	defer <- c

	// Local expansion as described in the paper
	newPerm := make([]int, len(perm)+1)
	for i, k := range perm {
		if k < a {
			newPerm[i] = k
		} else {
			newPerm[i] = k+1
		}
	}
	newPerm[len(perm)] = a
	p <- newPerm
}

func isPlane (perm []int) (bool) {
	n := len(perm)
	steps := make([]int, 0, n)
	for k:=0; k<n-1; k++ {
		if perm[k] < perm[k+1] - 1 {steps = append(steps, k)}
	}

	for _,s := range steps {
		m, M := perm[s], perm[s+1]
		two, three := 1000, 0
		prefix, suffix := perm[:s], perm[s+2:]
		for _,k := range prefix {
			if (k > m) && (k < M - 1) {
				two = min(k, two)
			}
		}

		for _,k := range suffix {
			if (k > two) && (k < M) {
				three = k
				return false
			}
		}
	}

	return true
}

func expand()

func main() {
	procs := 8
	runtime.GOMAXPROCS(procs)
	// Implement set
	curLevel := NewSliceSet()
	// Implement Add
	curLevel.Add([]int{1,2,3})
	curLevel.Add([]int{1,3,2})
	curLevel.Add([]int{2,1,3})
	curLevel.Add([]int{3,1,2})
	curLevel.Add([]int{2,3,1})
	curLevel.Add([]int{3,2,1})
	level := 3

	c := make(chan int, procs)
	for level < 20 {
		newLevel := NewSliceSet()
		p := make(chan []int, 100)
		go funct() {
			for _,perm := range curLevel {
				for a:=1; a<=level+1; a++ {
					go localExp(perm, a, p, c)
				}
			}
		}

		for newPerm := range p {
			go func() {
				c <- 1
				defer <- c
				// Implement Add
				if isPlane(newPerm) {newLevel.Add(newPerm)}
			}
		}

		fmt.Println(len(newLevel))
		curLevel = newLevel
		level += 1
	}
}
