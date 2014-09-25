package main

import (
	"bufio"
	"flag"
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

func openPort(p string) (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: p, Baud: 57600}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return port, err
}

func send(s io.ReadWriteCloser, msg []byte) int {
	n, err := s.Write(msg)
	if err != nil {
		log.Print(err)
		return 0
	}
	return n
}

func read(s io.ReadWriteCloser) []byte {
	buf := make([]byte, 256)
	n, err := s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf[:n]
}

func processLine(port io.ReadWriteCloser, line string) bool {
	send(port, []byte(line))
	log.Printf("%s", read(port))
	return true
}

func processFile(s io.ReadWriteCloser, file *os.File) bool {
	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		if processLine(s, line) == false {
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
		log.Print("A path to a plotter file must be provided.\n")
		return
	}

	port, err := openPort(serialport)

	if err != nil {
		log.Print("The given serial port could be opened.\n")
		return
	}

	file, err := os.Open(filepath)

	if err != nil {
		log.Println(err)
		return
	}

	processFile(port, file)
}
