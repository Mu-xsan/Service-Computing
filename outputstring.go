package main

import (
    "fmt"
    "os"
    "bufio"
    "io"
)

func main() {
	var fin *bufio.Reader
	fin = bufio.NewReader(os.Stdin)
	fmt.Printf("output something in outputstring:\n")
	for true {
            crc,err := fin.ReadString('\n') 
            if err!=nil || io.EOF == err {
                break
            }
            fmt.Printf(crc)
        }
}