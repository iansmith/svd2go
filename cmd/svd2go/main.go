package main

import (
	"flag"
	"github.com/iansmith/svd"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var outfile = flag.String("o", "", "output filename")
var dump = flag.Bool("d", false, "dump human readable version (debugging use only")
var pkg = flag.String("p", "main", "package to emit generated code into")
var tags = flag.String("b", "", "build tags (copied verbatim to output)")
var imp = flag.String("i", "runtime/volatile", "package name that has volatile.Register")


func main() {
	flag.Parse()
	if flag.NArg()==0 {
		log.Fatalf("usage svd2go -d -p <pkg> -o <outputfile> <input filename, either .csvd or .svd>")
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
		Tags: *tags,
		Import: *imp,
	}
	if strings.HasSuffix(flag.Arg(0),".csvd") {
		comp:=svd.ProcessCSVD(fp,opts)
		for _, part:=range comp.PeripheralFiles {
			fp, err:=os.Open(part)
			if err!=nil {
				log.Fatalf("unable to open %s (part of %s): %v", part, opts.InputFilename, err)
			}
			peripheral:=svd.ProcessPeripheral(fp,opts,part)
			comp.Add(peripheral)
		}
	}

	//it was _just_ an svd
	svd.ProcessSVD(fp, opts, opts.InputFilename)
}
