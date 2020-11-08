package main

type Link struct {
	step int8
	// minimum amount of steps to reach the rootnode
	// is 0 if next contains the rootnode
	// the root node (A) is the only node that leaves no gabs before it
	// akak
	// minimaal aantal stappen dat je moet doen om in een gesloten situatie te komen
	linksToRoot uint8
	next *[]*Link
}

func buildStarmap() *[]*Link {
	var A = &[]*Link{}

	/////////////
	var rootPlus1 = Link{
		step: +1,
		linksToRoot: 0,
		next: A,
	}

	/////////////
	var B = Link{step: -1,}

	/////////////
	var rootPlus2 = Link{
		step: +2,
		linksToRoot: 2,
		next: &[]*Link{&B},
	}

	////////////
	var C = &[]*Link{
		{step: +2, linksToRoot: 0, next: A},
		{step: +3, next: &[]*Link{&B}},
	}
	B.next = C

	////////////
	var rootPlus3 = Link{
		step: +3,
		linksToRoot: 3,
		next: &[]*Link{
			{-1, 2,&[]*Link{{-1, 1,&[]*Link{{+3, 0, A}}}}},
			{-2, 2, &[]*Link{{+1,1,  C}}},
			{+1, 3, &[]*Link{{-3, 2, &[]*Link{{+1, 1, &[]*Link{{+3, 0, A}}}}}}},
			{+2, 4, &[]*Link{{-3, 3, &[]*Link{{-1, 2, &[]*Link{{+3, 1, C}}}}}}},
		},
	}

	*A = append(*A, &rootPlus1)
	*A = append(*A, &rootPlus2)
	*A = append(*A, &rootPlus3)

	return A
}

var STARMAP = buildStarmap()

func CopyMap(m map[Resource]uint8) map[Resource]uint8 {
	cp := make(map[Resource]uint8)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
