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
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Print(err)
		return nil
	}
	return s
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

func read(s io.ReadWriteCloser) {
	buf := make([]byte, 128)
    n, err := s.Read(buf)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("%s", buf[:n])
}

func listen(ports []io.ReadWriteCloser) {
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, _ := reader.ReadString('\n')
		for _, port := range ports {
			send(port, []byte(msg))
			read(port)
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

	listen(ports)
}
