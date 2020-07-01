package main

import "math"

// Edge ... edge in the graph
type Edge struct {
	to, from uint32
}

func getEdge(to, from uint32) Edge {
	var edge Edge
	if to > from {
		edge.to, edge.from = from, to
	} else if to < from {
		edge.to, edge.from = to, from
	} else {
		var e error
		check(e, "Loop found")
	}

	return edge
}

// Collection ... collection of nodes
type Collection struct {
	nodes      []uint32
	intEdges   uint32
	outNodes   map[uint32]bool
	outEdges   map[Edge]bool
	totEdges   uint32
	density    float64
	modularity float64
}

func calcDensity(numNodes, numIntEdges uint32) float64 {
	if numNodes == 1 {
		return 1
	}

	density := float64(numIntEdges)
	density /= (0.5 * float64(numNodes) * float64(numNodes-1))

	return density
}

func calcMod(intEdges, totEdges, graphTotalEdges uint32) float64 {
	modularity := float64(intEdges) / float64(graphTotalEdges)

	modularity -= (0.25) * math.Pow(
		float64(intEdges+totEdges)/float64(graphTotalEdges), 2)

	return modularity
}

func (collection *Collection) initialize(id uint32, graph *Graph) {
	collection.nodes = make([]uint32, 1)
	collection.nodes[0] = id

	collection.intEdges = 0

	collection.outEdges = make(map[Edge]bool)
	collection.outNodes = make(map[uint32]bool)

	for to := range graph.list[id] {
		collection.outEdges[getEdge(to, id)] = true
		collection.outNodes[to] = true
	}

	collection.totEdges = uint32(len(collection.outEdges))

	collection.density = calcDensity(1, 0)
	collection.modularity = calcMod(0,
		collection.totEdges, graph.totEdges)
}

func calcComEdge(nodes2 []uint32, outNodes1 map[uint32]bool, outEdges1, outEdges2 map[Edge]bool) uint32 {
	found := false

	for _, node := range nodes2 {
		if outNodes1[node] == true {
			found = true
			break
		}
	}

	if !found {
		return 0
	}

	var comEdge uint32 = 0

	for edge2 := range outEdges2 {
		if outEdges1[edge2] {
			comEdge++
		}
	}

	return comEdge
}
