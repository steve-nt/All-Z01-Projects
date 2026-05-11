package main

import (
	"lem-in/funcs"
)

func main() {

	funcs.ParseInput()
	funcs.BuildConnections()
	funcs.StartEndConnection()
	funcs.PrintFile()
	funcs.VertexDisjointPaths()
	funcs.OptimalAntDistribution()

}
