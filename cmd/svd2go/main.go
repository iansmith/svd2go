package main

import (
	svd "github.com/iansmith/vsd"
	"log"
	"os"
	"flag"
)

var outfile = flag.String("o", "", "output filename")
func main() {
	flag.Parse()
	if flag.NArg()==0 {
		log.Fatalf("usage svd2go -o <outputfile> <input filename>")
	}
	fp, err:=os.Open(os.Args[1])
	if err!=nil {
		log.Fatal(err)
	}
	defer fp.Close()
	out:=os.Stdout
	if *outfile!="" {
		o, err:=os.Create(*outfile)
		if err!=nil {
			log.Fatal(err)
		}
		out=o
	}
	svd.ProcessSVD(fp,out)

}
