package main

func inExplored(needle Point, haystack []Point) bool {
	for _, n := range haystack {
		if n.Col == needle.Col && n.Row == needle.Row {
			return true
		}
	}
	return false
}
