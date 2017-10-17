package main

import (
    //"flag"
    "bufio"  
    "fmt"
    "strings"
    "io"
    "os"
    "strconv"
    "os/exec"
)
type selpg_args struct 
{
    start_page int
    end_page int
    in_filename string
    page_len int
    page_type string 
    print_dest string
}
type sp_args selpg_args
const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
var (
    progname = "selpg"
)
func process_args(ac int, av []string, psa *selpg_args) {
    var s1 string
    var s2 string
    var i int
    var argno int
    var errmsg string
    errmsg = progname+": "
    if ac < 3  {
        usage()
        errmsg = errmsg+"not enough arguments\n"
        panic(errmsg)
    }
    s1 = av[1]
    if !strings.HasPrefix(s1,"-s") {
        usage();
        errmsg = errmsg+"1st arg should be -sstart_page\n"
        panic(errmsg)
    }
    i,err := strconv.Atoi(s1[2:len(s1)])
    _ = err
   
    if  err != nil || i < 1 || i > (MaxInt - 1) {
        usage()
        errmsg = errmsg+"invalid start page "+s1[2:len(s1)]+"\n"
        panic(errmsg)
    }
    (*psa).start_page = i;

    s1 = av[2]
    if !strings.HasPrefix(s1,"-e") {
        usage()
        errmsg = errmsg+"2nd arg should be -eend_page\n"
        panic(errmsg)
    }

    i,err = strconv.Atoi(s1[2:len(s1)])
    _ = err
   
    if  err != nil || i < 1 || i > (MaxInt - 1) || i < (*psa).start_page {
        usage()
        errmsg = errmsg+"invalid end page "+s1[2:len(s1)]+"\n"
        panic(errmsg)
    }
    (*psa).end_page = i;

    argno = 3
    for argno <= (ac - 1) && av[argno][0] == '-' {
        s1 = av[argno]
        if s1[1] == 'l' {
            s2 = s1[2:len(s1)]
            i,err = strconv.Atoi(s2)
            if  err != nil || i < 1 || i > (MaxInt - 1) {
                    usage()
                    errmsg = errmsg+"invalid page length "+s2+"\n"
                    panic(errmsg)
                }
                (*psa).page_len = i
                argno = argno + 1
        } else if s1[1] == 'f' {
            if !strings.HasPrefix(s1,"-f") {
                    usage()
                    errmsg = errmsg+"option should be \"-f\"\n"
                    panic(errmsg)
                }
                (*psa).page_type = "f";
                argno = argno + 1
        } else if s1[1] == 'd' {
            s2 = s1[2:len(s1)] 
                if (len(s2) < 1) {
                    usage()
                    errmsg = errmsg+"-d option requires a printer destination\n"
                    panic(errmsg)
                }
                (*psa).print_dest = s2
                argno = argno + 1
        } else {
                usage()
                errmsg = errmsg+"unknown option "+s1+"\n"
                panic(errmsg)
        }
    }

    if argno <= ac-1 {
        (*psa).in_filename = av[argno]
        _, ero := os.Open(av[argno])
        if ero != nil {
            usage()
            errmsg = errmsg+"input file \""+av[argno]+"\" does not exist\n"
            panic(errmsg)
        }
    }
}
func process_input(sa selpg_args) {
    var fin *bufio.Reader
    var fout *bufio.Writer
    var f *os.File
    var err error
    var errmsg string
    if sa.in_filename == "" {
        fmt.Printf("start\n")
        fin = bufio.NewReader(os.Stdin)
    } else {
        f, err = os.Open(sa.in_filename)
        if err != nil {
            usage()
            errmsg = errmsg+"could not open input file \""+sa.in_filename+"\"\n"
            panic(errmsg)
        }
        fin = bufio.NewReader(f)
    }
     another := exec.Command("cat","-n")
     in,_ := another.StdinPipe()
    if sa.print_dest != "" {
        f,err = os.OpenFile(sa.print_dest,os.O_RDWR|os.O_APPEND,0644)
        if (err != nil) {
           // usage()
           // errmsg = errmsg+"could not open pipe to \""+sa.print_dest+"\"\n"
           // panic(errmsg)
        }
        fout = bufio.NewWriter(f)
       

    } else {
        fout = bufio.NewWriter(os.Stdout)
    }
    
    if sa.page_type == "l" {
        line_ctr := 0
        page_ctr := 1
        for true {
            crc,err := fin.ReadString('\n') 
            if err!=nil || io.EOF == err {
                break
            }
            line_ctr = line_ctr+1
            if line_ctr > sa.page_len {
                page_ctr = page_ctr + 1
                line_ctr = 1;
            }
            if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
                if sa.print_dest != "" {
                    in.Write([]byte(crc))
                } else {
                fout.WriteString(crc)
                fout.Flush()
                }
            }
        }
        if sa.print_dest != "" {
            in.Close()
            another.Stdout = os.Stdout
            another.Start()
        }
    } else {
        page_ctr := 1;
        for true  {
            c, _, err := fin.ReadRune()
            if err != nil || err == io.EOF { /* error or EOF */
                break
            }
            if (c == '\f') { /* form feed */
                page_ctr = page_ctr+1
            }
            if  (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
                fmt.Printf("%c",c)
            }
        }
        fmt.Printf("\n")
    }
    return
}
func usage() {
    fmt.Printf("\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progname)
}
func main() {
    av := os.Args
    ac := len(av)

    var sa selpg_args
    sa = selpg_args{start_page:-1,end_page:-1,in_filename:"",page_len:72,page_type:"l",print_dest:""}
    
    progname = av[0]
    _ = progname
    
    process_args(ac, av, &sa)
    process_input(sa)

}
