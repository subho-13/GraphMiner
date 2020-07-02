package main

import (
	"bufio"
	"fmt"
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

	tickerWriter := time.NewTicker(time.Minute * 5)
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

	for stop != true {
		index = index % set.numCollections
		if set.coagulate(index, graph) {
			index += 2*index + 1
		} else {
			index += 2
		}

		select {
		case <-tickerPrinter.C:
			fmt.Println(set.numCollections, set.modularity+set.regularization)
		case <-tickerWriter.C:
			{
				fmt.Println("Writing file")
				set.writeRes(path, outName, graph)
				index = randNum.Uint32()
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
