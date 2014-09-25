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

func listPorts(port string) (list []string) {
	if len(port) > 0 {
		list = []string{port}
	} else {
		list = listSerialPorts()
	}
	return list
}

func openPort(p string) io.ReadWriteCloser {
	c := &serial.Config{Name: p, Baud: 57600}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Print(err)
		return nil
	}
	return port
}

func openPorts(list []string) (ports []io.ReadWriteCloser) {
	for _, port := range list {
		ports = append(ports, openPort(port))
	}
	return ports
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

func process(s []io.ReadWriteCloser, file *os.File) {
	r := bufio.NewReader(file)
	for {
        line, err := r.ReadString('\n')
        if len(line) > 0 {
            // log.Print(line)
            for _, port := range s {
            	send(port, []byte(line))
            	log.Printf("%s", read(port))
            }
        }
        if err != nil {
            break
        }
    }
}

func main() {

	var p = flag.String("p", "", "the USB port to use")
	var l = flag.Bool("l", false, "list all avliable serial ports")
	flag.Parse()

	if *l {
		log.Printf("%v", listSerialPorts())
		return
	}

	filepath := flag.Arg(0)

	if flag.Arg(0) == "" {
		log.Print("A path to a plotter file must be provided.\n")
		return
	}

	list := listPorts(*p)

	if len(list) == 0 {
		log.Print("No serial ports found.\n")
		return
	}

	ports := openPorts(list)

	if len(ports) == 0 {
		log.Print("No serial ports could be opened.\n")
		return
	}

	file, err := os.Open(filepath)

	if err != nil {
		log.Println(err)
		return
	}

	process(ports, file)
}
