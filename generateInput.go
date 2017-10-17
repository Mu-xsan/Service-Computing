package main

import (
    "flag"
    "fmt"
    "strconv"
    "os"
    "bufio"
)

func main() {
	var filename string
	var prefix string
	var num int
	var kind int
	flag.StringVar(&filename,"n","input_file","The name of the input file. Defaults to input_file.")
	flag.StringVar(&prefix,"p","input_file's","The prefix of every line in the input file. Defaults to input_file's.")
	flag.IntVar(&num,"c",100,"The num of character. Defaults to 100.")
	flag.IntVar(&kind,"t",1,"1 for \\n, 0 for \\f. Defaults to 1.")
	flag.Parse()
    /*
	fmt.Printf("filename = %s\n", filename)
	fmt.Printf("prefix = %s\n", prefix)
	fmt.Printf("num = %d\n", num)
	fmt.Printf("kind = %d\n", kind)
    */
	var fout *bufio.Writer
	var f *os.File
    var err error
    f,err = os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
    if (err != nil) {
        fmt.Printf("%s: could not open pipe to \"%s\"\n",os.Args[0], filename);
        os.Exit(1);
    }
    _ = f
    //fout = bufio.NewWriter(f)
    fout = bufio.NewWriter(os.Stdout)
    fmt.Printf("The generateInput's Stdout is:\n")
    var i int
    var temp string
    var charac string
    if kind == 1{
    	charac = "\n"
    } else {
    	charac = "\f"
    }
    i = 1
    for true {
    	if i > num {
    		break
    	}
    	temp = prefix + strconv.Itoa(i) + charac
    	fout.WriteString(temp)
    	i = i+1
    }
    fout.Flush()

}