package main

import (
	"strconv"
)

type Link struct {
	step int8
	// minimum amount of steps to reach the rootnode
	// is 0 if next contains the rootnode
	// the root node (A) is the only node that leaves no gabs before it
	// akak
	// minimaal aantal stappen dat je moet doen om in een gesloten situatie te komen
	linksToRoot uint8
	next        *[]*Link
}

func nTimesM(n uint8, m int8, linksToRoot uint8, last *[]*Link) *[]*Link {
	if n == 0 {
		return last
	}
	return &[]*Link{{
		step:        m,
		linksToRoot: linksToRoot,
		next:        nTimesM(n-1, m, linksToRoot-1, last),
	}}
}

var A = &[]*Link{}
var B = Link{step: -1,}
var C = &[]*Link{
	{step: +2, linksToRoot: 0, next: A},
	{step: +3, next: &[]*Link{&B}},
}

func D(depth uint8) *Link {
	l := 3 + (depth * 2)

	nextLinks := []*Link{
		{-1, l - 1,
			nTimesM(depth, -3, l-2,
				&[]*Link{{-1, l - 2 - depth,
					nTimesM(depth, +3, l-3-depth,
						&[]*Link{{+3, 0, A}})}})},
		{-2, l - 1,
			nTimesM(depth, -3, l-2,
				&[]*Link{{+1, l - 2 - depth,
					nTimesM(depth, +3, l-3-depth, C)}})},
		{+1, l,
			nTimesM(depth+1, -3, l-1,
				&[]*Link{{+1, l - depth - 2,
					nTimesM(depth+1, +3, l - depth - 3, A)}})},
	}
	if depth < 30 {
		nextLinks = append(nextLinks, D(depth+1))
	}

	return &Link{
		step:        +3,
		linksToRoot: l,
		next:        &nextLinks,
	}
}

func E(depth uint8) *Link {
	l := 2 + 2*depth

	nextLinks := []*Link{
		{-1, l - 1,
			nTimesM(depth, -3, l-2,
				&[]*Link{{+2, l - 2 - depth,
					nTimesM(depth, +3, l-3-depth, A)}})},
		{-2, l - 1,
			nTimesM(depth-1, -3, l-2,
				&[]*Link{{-2, l - 1 - depth,
					nTimesM(depth, +3, l - 2 - depth, C)}})},
		{+1, l,
			nTimesM(depth, -3, l-1,
				&[]*Link{{-2, l-1-depth,
					nTimesM(depth+1, +3, l-2-depth, A)}})},
	}

	if depth < 50 {
		nextLinks = append(nextLinks, E(depth+1))
	}
	return &Link{
		step:        +3,
		linksToRoot: l,
		next:        &nextLinks,
	}
}

func buildStarmap() *[]*Link {

	/////////////
	var rootPlus1 = Link{
		step:        +1,
		linksToRoot: 0,
		next:        A,
	}

	/////////////
	var e = E(1)
	var rootPlus2 = Link{
		step:        +2,
		linksToRoot: 2,
		next:        &[]*Link{
			&B,
			{+1, 2, &[]*Link{{-2, 1, &[]*Link{{+3, 0, A}}}}},
			{+2, 3, &[]*Link{{-3, 2, &[]*Link{{+2, 1, C}}}}},
			e,
		},
	}

	////////////
	B.next = C

	////////////
	var d = D(1)
	var rootPlus3 = Link{
		step:        +3,
		linksToRoot: 3,
		next: &[]*Link{
			{-1, 2, &[]*Link{{-1, 1, &[]*Link{{+3, 0, A}}}}},
			{-2, 2, &[]*Link{{+1, 1, C}}},
			{+1, 3, &[]*Link{{-3, 2, &[]*Link{{+1, 1, &[]*Link{{+3, 0, A}}}}}}},
			{+2, 4, &[]*Link{{-3, 3, &[]*Link{{-1, 2, &[]*Link{{+3, 1, C}}}}}}},
			d,
		},
	}

	*A = append(*A, &rootPlus1)
	*A = append(*A, &rootPlus2)
	*A = append(*A, &rootPlus3)

	//printLink(e, "")

	return A
}

func printLink(d *Link, indent string) {
	print(indent, d.step, "("+strconv.Itoa(int(d.linksToRoot))+")")
	if d.next == A {
		println(" A")
		return
	}
	if d.next == C {
		println(" C")
		return
	}
	if d.next == nil {
		println("NILL")
		return
	}
	if len(*d.next) < 2 {
		printLink((*d.next)[0], indent)
		return
	}
	println()
	for _, l := range *d.next {
		printLink(l, indent+"\t")
	}
}

var STARMAP = buildStarmap()

func CopyMap(m map[Resource]uint8) map[Resource]uint8 {
	cp := make(map[Resource]uint8)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
