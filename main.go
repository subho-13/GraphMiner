package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const path = "/home/subho"

var randNum = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	fmt.Printf("Enter input filename  :: ")
	scanner.Scan()
	inName := scanner.Text()
	fmt.Printf("Enter output filename :: ")
	scanner.Scan()
	outName := scanner.Text()

	tickerWriter := time.NewTicker(time.Second * 11)
	tickerPrinter := time.NewTicker(time.Second * 10)
	defer tickerWriter.Stop()
	defer tickerPrinter.Stop()

	sigs := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	}()

	graph := new(Graph)
	graph.ReadGraph(path, inName)

	set := new(Set)
	set.initialize(graph, path, outName)

	var index uint32 = 0
	stop := false

	var oldVal, oldValWrite, newVal float64

	for stop != true {
		index = index % set.numCollections
		set.coagulate(index, graph)
		index++

		newVal = set.modularity + set.regularization

		select {
		case <-tickerPrinter.C:
			{
				fmt.Println(set.numCollections, newVal)
			}
		case <-tickerWriter.C:
			{
				if newVal > oldValWrite {
					fmt.Println("Writing file")
					set.writeRes(path, outName, graph)
					oldValWrite = newVal
				}
				if math.Abs(oldVal-newVal) < 0.000001 {
					min := 100.0
					var ind uint32 = 0
					var i uint32
					n := set.numCollections - randNum.Uint32()%set.numCollections
					m := n - randNum.Uint32()%(n/4)
					for i = m; i < n; i++ {
						c := set.collections[i]
						val := (c.modularity + c.density) / math.Log10(float64(len(c.nodes)))
						if len(c.nodes) > 1 && val < min {
							ind = i
							min = val
						}
					}

					set.split(ind, graph)
				}
				oldVal = newVal
			}

		default:
		}

		select {
		case <-sigs:
			{
				fmt.Println("Stopping Execution. Writing files")
				set.writeRes(path, outName, graph)
				stop = true
			}
		default:
		}
	}

}
