package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/goserial"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Looks into the "/dev" directory and returns all the files that maybe serial ports.
func listSerialPorts() (list []string) {
	dir := "/dev"
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "cu.") && (strings.Index(f.Name(), "usbserial") >= 0 || strings.Index(f.Name(), "usbmodem") >= 0 || strings.Index(f.Name(), "ttyUSB") >= 0) {
			list = append(list, path.Join(dir, f.Name()))
		}
	}
	return list
}

func openPortMock(p string) (io.ReadWriteCloser, error) {
	port, _ := os.Open(p)
	return port, nil
}

func openPort(p string) (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: p, Baud: 57600}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return port, err
}

func write(port io.Writer, cmd []byte) int {
	n, err := port.Write(cmd)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func processCommand(w io.Writer, r *bufio.Scanner, cmd string) bool {
	// fmt.Print("CMD: " + cmd)
	write(w, []byte(cmd + "\n"))
	for r.Scan() {
		ack := r.Text()
		if ack == "OK" {
			return true
		}
		fmt.Println(ack)
	}
	return false
}

func processFile(port io.ReadWriteCloser, file *os.File) bool {
	p := bufio.NewScanner(port)
	r := bufio.NewScanner(file)
	processCommand(port, p, "\n") // Read the boot up text.
	for r.Scan() {
		cmd := r.Text()
		if processCommand(port, p, cmd) == false {
			return false
		}
	}
	return true
}

func main() {

	var l = flag.Bool("l", false, "list all available serial ports")
	var m = flag.Bool("m", false, "use the given serial port value as a mock vplotter")
	flag.Parse()

	if *l {
		fmt.Printf("%v\n", listSerialPorts())
		return
	}

	serialport := flag.Arg(0)
	filepath := flag.Arg(1)

	if serialport == "" {
		fmt.Print("A serial port must be provided.\n")
		return
	}

	if filepath == "" {
		fmt.Print("A path to a vplotter file must be provided.\n")
		return
	}

	var port io.ReadWriteCloser
	var err error

	if *m {
		port, err = openPortMock(serialport)
	} else {
		port, err = openPort(serialport)
	}

	if err != nil {
		fmt.Print("The given serial port could not be opened.\n")
		return
	}

	defer port.Close()

	file, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
		return
	}

	if processFile(port, file) == false {
		fmt.Print("Error.\n")
	}

	fmt.Print("Completed.\n")
}
