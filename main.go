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
	cpu, workers, jobs, count int
	output, input, port       string
	hideDialErr, exists       bool
	scanner                   *bufio.Scanner
	timeout                   = 250 * time.Millisecond
	mutex                     sync.Mutex
	buf                       bytes.Buffer
	done                      chan bool
	outputFile                *os.File
	err                       error
)

func init() {
	fmt.Printf(logo+"\n", color(Logo, "PORTakal"))

	flag.IntVar(&cpu, "cpu", runtime.NumCPU(), "-cpu=8 (if empty use all cpus)")
	flag.StringVar(&port, "port", "", "-port=3386")
	flag.BoolVar(&hideDialErr, "errors", true, "-errors")
	flag.StringVar(&output, "output", "", "-output=output.txt")
	flag.StringVar(&input, "input", "", "-input=live.txt")
	flag.IntVar(&workers, "workers", 10, "-workers=100")
	flag.IntVar(&jobs, "jobs", 10, "-jobs=50")
	flag.IntVar(&count, "count", 10, "-count=10")

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

func writeLine(file *os.File, host string) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, err := file.WriteString(host + "\n")

	if err != nil {
		return err
	}
	return nil

}

func fileExists(filename string) (file *os.File, exists bool, err error) {
	if _, err = os.Stat(output); os.IsNotExist(err) {
		return nil, false, nil
	} else {
		file, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		return file, true, err
	}
}

func createFile(filename string) (outputFile *os.File, err error) {

	file, err := os.Create(filename)

	if err != nil {
		return
	}

	defer file.Close()

	outputFile, err = os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		return
	}

	return
}

func initOutputFile(output string) (outputFile *os.File, err error) {

	outputFile, exists, err = fileExists(output)

	if err != nil {
		return
	}

	if !exists {

		outputFile, err = createFile(output)

		if err != nil {
			return
		}

	} else {

		outputFile, err = os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			return
		}

	}

	return
}
