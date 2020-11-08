package main

type Link struct {
	step int8
	next []*Link
}

func buildStarmap() *Link {
	var A = Link {
		step: 0,
		next: []*Link{},
	}

	/////////////
	var rootPlus1 = Link{
		step: +1,
		next: nil,
	}
	rootPlus1.next = A.next

	/////////////
	var B = Link {step: -1,	}

	/////////////
	var rootPlus2 = Link {
		step: +2,
		next: []*Link{&B},
	}

	////////////
	var C = Link {step: 0,
		next: []*Link{
			{step: +2, next: []*Link{&A}},
			{step: +3, next: []*Link{&B}},
		},
	}
	B.next = []*Link{&C}

	////////////
	var rootPlus3 = Link{
		step: +3,
		next: []*Link{
			{-1, []*Link{{-1, []*Link{{+3, []*Link{&A}}}}}},
			{-2, []*Link{{+1, []*Link{&C}}}},
			{+1, []*Link{{-3, []*Link{{+1, []*Link{{+3, []*Link{&A}}}}}}}},
			{+2, []*Link{{-3, []*Link{{-1, []*Link{{+3, []*Link{&C}}}}}}}},
		},
	}
	rootPlus1.next = []*Link{&A}

	A.next = []*Link{&rootPlus1, &rootPlus2, &rootPlus3}
	return &A
}

var STARMAP = buildStarmap()



func CopyMap(m map[Resource]uint8) map[Resource]uint8 {
	cp := make(map[Resource]uint8)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}