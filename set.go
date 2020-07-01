package main

import (
	"fmt"
	"os"
)

// Set ... Set of collections
type Set struct {
	collections    []*Collection
	numCollections uint32
	regularization float64
	modularity     float64
}

func (set *Set) initialize(graph *Graph) {
	set.numCollections = graph.totVertex
	set.collections = make([]*Collection, set.numCollections)

	count := 0

	for node := range graph.list {
		collection := new(Collection)
		collection.initialize(node, graph)
		set.modularity += collection.modularity
		set.regularization += collection.density
		set.collections[count] = collection
		count++
	}
	set.regularization *= 1 / float64(set.numCollections)
	set.regularization -= float64(set.numCollections) / float64(graph.totVertex)
	set.regularization *= 0.5
}

// IndexIncrease ... Send (index, increase) in a channel
type IndexIncrease struct {
	index    uint32
	increase float64
}

func (set *Set) coagulate(i uint32, graph *Graph) bool {
	indexIncreaseChan := make(chan IndexIncrease, set.numCollections-1)

	var j uint32
	for j = 0; j < set.numCollections; j++ {
		if i != j {
			go func() {
				newMod := costNewMod(
					set.collections[i], set.collections[j], graph.totEdges)
				newReg := costNewReg(
					set.collections[i], set.collections[j], set.regularization,
					set.numCollections, graph.totVertex)

				increaseMod := newMod - (set.collections[i].modularity + set.collections[j].modularity)
				increaseReg := newReg - set.regularization
				indexIncreaseChan <- IndexIncrease{index: j,
					increase: increaseMod + increaseReg}
			}()
		}
	}

	max := 0.0
	var index uint32

	for j = 0; j < set.numCollections; j++ {
		if i != j {
			inInc := <-indexIncreaseChan
			if inInc.increase > max {
				index = inInc.index
				max = inInc.increase
			}
		}
	}

	if max > 0 {
		set.regularization = costNewReg(
			set.collections[i], set.collections[index], set.regularization,
			set.numCollections, graph.totVertex)
		set.modularity -= (set.collections[i].modularity + set.collections[index].modularity)

		merge(set.collections[i], set.collections[index], graph.totEdges)
		set.modularity += set.collections[i].modularity
		set.collections[index], set.collections[set.numCollections-1] = set.collections[set.numCollections-1], set.collections[index]
		set.numCollections--
		return true
	}

	return false
}

func merge(c1, c2 *Collection, totalEdges uint32) {
	c1.nodes = append(c1.nodes, c2.nodes...)

	commOutEdge := 0

	for edge := range c2.outEdges {
		if c1.outEdges[edge] == true {
			commOutEdge++
		} else {
			c1.outEdges[edge] = true
		}
	}

	c1.intEdges += c2.intEdges + uint32(commOutEdge)
	c1.totEdges += c2.totEdges - uint32(commOutEdge)

	c1.density = calcDensity(uint32(len(c1.nodes)), c1.intEdges)
	c1.modularity = calcMod(c1.intEdges, c1.totEdges, totalEdges)
}

func costNewMod(c1, c2 *Collection, totalEdges uint32) float64 {
	numCom := calcComEdge(c1.outEdges, c2.outEdges)

	newIntEdges := c1.intEdges + c2.intEdges + numCom
	newTotEdges := c1.totEdges + c2.totEdges - numCom
	newMod := calcMod(newIntEdges, newTotEdges, totalEdges)

	return newMod
}

func costNewReg(c1, c2 *Collection, oldRVal float64, n, v uint32) float64 {
	densitySum := oldRVal*2 + (float64(n) / float64(v))
	densitySum *= float64(n)

	newDensitySum := densitySum - (c1.density + c2.density)

	numCom := calcComEdge(c1.outEdges, c2.outEdges)
	newNodes := uint32(len(c1.nodes) + len(c2.nodes))
	newIntEdges := c1.intEdges + c2.intEdges + numCom

	newDensitySum += calcDensity(newNodes, newIntEdges)

	newRVal := 1/float64(n-1)*newDensitySum - float64(n-1)/float64(v)
	newRVal *= 0.5

	return newRVal
}

func (set *Set) writeRes(path, name string) {
	filename := path + "/" + name
	file, err := os.Create(filename)
	check(err, "File cannot be created")
	defer file.Close()

	var i uint32

	fmt.Fprintf(file, "Q = %f \n", set.modularity+set.regularization)


	for i = 0; i < set.numCollections; i++ {
		for _, node := range set.collections[i].nodes {
			fmt.Fprintf(file, "%d ", node)
		}
		fmt.Fprintf(file, "\n")
	}
}
