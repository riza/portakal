package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	logo = `
	:ms-                 
    mdmd/              
     /mddmd:             
      +mdddm/ .-//+++++:
   -+ssddmddd+hddddddmdy-
 -yy/--:+ohNMmmddddhs/.  
:m+-/+++++hysyN:..     
ms-:+++++++++++hy        
Ns+++++++++++++hy        
/mo+++++++++++sm-        
 :dho+++++++shh.         
   -oyhhhhhyo. %s - Bulk port checker
	`
)

var (
	cpu                 int
	output, input, port string
	hideDialErr, exists bool
	scanner             *bufio.Scanner
	timeout             = 250 * time.Millisecond
	mutex               sync.Mutex
	buf                 bytes.Buffer
	done                chan bool
	outputFile          *os.File
	err                 error
)

func init() {
	fmt.Printf(logo+"\n", color(Logo, "PORTakal"))

	flag.IntVar(&cpu, "cpu", runtime.NumCPU(), "-cpu=8 (if empty use all cpus)")
	flag.StringVar(&port, "port", "", "-port=3386")
	flag.BoolVar(&hideDialErr, "errors", true, "-errors")
	flag.StringVar(&output, "output", "", "-output=output.txt")
	flag.StringVar(&input, "input", "", "-input=live.txt")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "USAGE:")
	fmt.Fprintln(os.Stderr, "./portakal -cpu=8 -port=3389 -input=list.txt -output=live.txt")
	fmt.Fprintln(os.Stderr, "")
}

func main() {

	flag.Parse()

	if input == "" {
		usage()
		os.Exit(1)
	}

	outputFile, err = initOutputFile(output)

	if err != nil {
		errMsg(err, true)
	}

	defer outputFile.Close()

	inputFile, err := os.Open(input)

	if err != nil {
		errMsg(err, true)
	}

	tee := io.TeeReader(inputFile, &buf)
	count, err := lineCounter(tee)

	if err != nil {
		errMsg(err, true)
	}

	info("PORTakal scanning your list", false)
	info(fmt.Sprintf("%d host ready for scan", count), false)
	info(fmt.Sprintf("%s ETA Scan time\n", time.Duration(count)*(timeout)), false)

	scanner = bufio.NewScanner(&buf)
	done = make(chan bool)

	start := time.Now()
	go scanAndCheck(done, outputFile)
	<-done

	elapsed := time.Since(start)
	info(fmt.Sprintf("Checker took %s", elapsed), true)

}

func scanAndCheck(done chan bool, outputFile *os.File) {

	for scanner.Scan() {
		checkHost(outputFile, scanner.Text())
	}

	done <- true
}

func checkHost(outputFile *os.File, addr string) {

	if port != "" {
		addr = net.JoinHostPort(addr, port)
	}

	ok, err := dial(addr)

	if err != nil && !hideDialErr {
		errMsg(err, true)
	}

	if ok {

		live(addr, false)

		err = writeLine(outputFile, addr)

		if err != nil {
			errMsg(err, true)
		}

	} else {
		dead(addr, false)
	}
}

func dial(host string) (bool, error) {

	conn, err := net.DialTimeout("tcp", host, timeout)

	if err != nil {
		return false, err
	}

	if conn != nil {
		conn.Close()
		return true, nil
	}

	return false, nil
}
