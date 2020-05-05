package main

import (
	svd "github.com/iansmith/vsd"
	"log"
	"os"
	"flag"
)

var outfile = flag.String("o", "", "output filename")
var dump = flag.Bool("d", false, "dump human readable version")
var pkg = flag.String("p", "main", "package to emit generated code into")
var tags = flag.String("b", "", "build tags (copied verbatim to output)")


func main() {
	flag.Parse()
	if flag.NArg()==0 {
		log.Fatalf("usage svd2go -d -p <pkg> -o <outputfile> <input filename>")
	}
	fp, err:=os.Open(flag.Arg(0))
	if err!=nil {
		log.Fatal(err)
	}
	defer fp.Close()
	out:=os.Stdout
	if *outfile!="" {
		o, err:=os.Create(*outfile)
		if err!=nil {
			log.Printf("out is %+v",*outfile)
			log.Fatal(err)
		}
		out=o
	}
	defer out.Close()
	opts:=&svd.UserOptions{
		Out:out,
		Dump: *dump,
		Pkg: *pkg,
		InputFilename: flag.Arg(0),
	}
	svd.ProcessSVD(fp,opts)

}
