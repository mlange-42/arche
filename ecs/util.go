package ecs

func capacity(size, increment int) int {
	cap := increment * (size / increment)
	if size%increment != 0 {
		cap += increment
	}
	return cap
}

func capacityU32(size, increment uint32) uint32 {
	cap := increment * (size / increment)
	if size%increment != 0 {
		cap += increment
	}
	return cap
}
