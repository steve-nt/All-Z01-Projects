package core

import (
	"log"
	"os"
)

func Run() {
	if len(os.Args) != expArgc {
		log.Fatal("invalid argument count. expect 2 arguments corresponding to input/output file names")
	}
	inpath := os.Args[1]
	outpath := os.Args[2]
	buf, err := os.ReadFile(inpath)
	if err != nil {
		log.Fatal(err.Error())
	}
	tokens := tokenizeInput(string(buf))
	exe := executeCommands(tokens)
	output := buildOutput(exe)
	err = os.WriteFile(outpath, []byte(output), 0664)
	if err != nil {
		log.Fatal(err.Error())
	}
}
