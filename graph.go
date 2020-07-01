package main

import (
	"bufio"
	"os"
	"strconv"
)

// Graph ... Contains the graph
type Graph struct {
	list      map[uint32]map[uint32]bool
	totEdges  uint32
	totVertex uint32
}

// addEdge ... Add edge to a graph
func (graph *Graph) addEdge(to, from uint32) {
	if _, existsTo := graph.list[to]; existsTo {
		if !graph.list[to][from] {
			graph.list[to][from] = true
		}
	} else {
		graph.list[to] = make(map[uint32]bool)
		graph.list[to][from] = true
	}
}

// ReadGraph ... Read graph from a file
func (graph *Graph) ReadGraph(path, name string) {
	graph.list = make(map[uint32]map[uint32]bool)

	filename := path + "/" + name
	file, err := os.Open(filename)
	check(err, "Problem opening file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	totEdges := 0

	for {

		if !scanner.Scan() {
			break
		}

		from, er1 := strconv.ParseUint(scanner.Text(), 10, 32)
		check(er1, "Couldn't read integer")

		if !scanner.Scan() {
			break
		}

		to, er2 := strconv.ParseUint(scanner.Text(), 10, 32)
		check(er2, "Couldn't read integer")

		graph.addEdge(uint32(to), uint32(from))
		graph.addEdge(uint32(from), uint32(to))

		totEdges++

	}

	graph.totEdges = uint32(totEdges)
	graph.totVertex = uint32(len(graph.list))
}
