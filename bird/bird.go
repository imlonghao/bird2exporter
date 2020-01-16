package bird

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
)

var regexCommand = regexp.MustCompile(`^\d\d\d\d`)

type Bird struct {
	conn  net.Conn
	mutex *sync.Mutex
}

func New(path string) *Bird {
	birdConn, err := net.Dial("unix", path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewReader(birdConn)
	line, _, err := scanner.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Bird loaded, %s", line)
	return &Bird{
		conn:  birdConn,
		mutex: &sync.Mutex{},
	}
}

func (b *Bird) write(command string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if _, err := b.conn.Write([]byte(command + "\n")); err != nil {
		log.Fatal("fail to write command, err: %s", err)
	}
}

func (b *Bird) read() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	scanner := bufio.NewScanner(b.conn)
	txt := ""
	for scanner.Scan() {
		line := scanner.Text()
		if regexCommand.MatchString(line) {
			if line[4] == 45 {
				line = line[5:]
			} else {
				break
			}
		}
		txt += strings.TrimSpace(line) + "\n"
	}
	return strings.TrimSuffix(txt, "\n")
}
