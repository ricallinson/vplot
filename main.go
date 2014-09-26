package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/goserial"
	"io"
	"io/ioutil"
	"log"
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

func write(port io.ReadWriteCloser, msg []byte) int {
	n, err := port.Write(msg)
	if err != nil {
		log.Print(err)
		return 0
	}
	return n
}

func processCommand(w io.Writer, r *bufio.Reader, cmd string) bool {
	// fmt.Print("CMD: " + cmd)
	// write(w, []byte(cmd))
	for {
		ack, err := r.ReadString('\n')
		fmt.Print(ack)
		if ack == "OK\n" {
			return true
		}
		if err != nil {
			return false
		}
	}
	return false
}

func processFile(port io.ReadWriteCloser, file *os.File) bool {
	p := bufio.NewReader(port)
	r := bufio.NewReader(file)
	processCommand(port, p, "\n") // boot up
	for {
		cmd, err := r.ReadString('\n')
		if processCommand(port, p, cmd) == false {
			return false
		}
		if err != nil {
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
		log.Printf("%v", listSerialPorts())
		return
	}

	serialport := flag.Arg(0)
	filepath := flag.Arg(1)

	if serialport == "" {
		log.Print("A serial port must be provided.\n")
		return
	}

	if filepath == "" {
		log.Print("A path to a vplotter file must be provided.\n")
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
		log.Print("The given serial port could not be opened.\n")
		return
	}

	defer port.Close()

	file, err := os.Open(filepath)

	if err != nil {
		log.Println(err)
		return
	}

	processFile(port, file)
}
