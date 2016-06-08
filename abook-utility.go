package main

import (
	"flag"
	"fmt"
	"github.com/ichigo663/abook-utility/abook-lib"
	"os"
)

func main() {
	var outFile *os.File
	var usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... abook-file\n", os.Args[0])
		flag.PrintDefaults()
	}
	removeMailOrphan := flag.Bool("mail-only", false, "Remove contacts without a mail address")
	removeWith := flag.String("remove-with", "", "Remove contacts with a mail address containing \"string\"")
	outFilePath := flag.String("out", "", "Path of the output file (Default is stdout)")
	showHelp := flag.Bool("help", false, "Show this help")
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	flag.Parse()
	if *showHelp {
		usage()
		os.Exit(0)
	}
	abookFile, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//init the Abook struct
	abookDB := abook.NewAbook(abookFile)
	err = abookDB.Parse()
	abookFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if *removeMailOrphan || *removeWith != "" {
		abookDB.Remove(*removeMailOrphan, *removeWith)
	}
	if *outFilePath != "" {
		outFile, err = os.Create(*outFilePath)
		defer outFile.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		outFile = os.Stdout
	}
	abookDB.WriteTo(outFile)
}
