package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	listner           net.Listener
	db                Database
	connections       map[int]net.Conn
	totalConnnections int
	quit              chan int
	exited            chan int
}

func NewServer() *Server {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("DB Error")
	}

	ser := &Server{
		listner:           ln,
		db:                NewDatabase(),
		connections:       make(map[int]net.Conn),
		totalConnnections: 0,
		quit:              make(chan int),
		exited:            make(chan int),
	}
	go ser.Serve()
	return ser
}

func (s *Server) Serve() {
	for {
		select {
		case <-s.quit:
			fmt.Println("Shutting down the server")
			err := s.listner.Close()
			if err != nil {
				fmt.Println("Some issue in shutdown")
			}

			for _, conn := range s.connections {
				s.write(conn, "server closing in 10sec")
				<-time.After(10 * time.Second)
				conn.Close()
			}

			close(s.exited)
			fmt.Println("Handle server close here")
		default:
			fmt.Println("Accept new connect here")
			listener, _ := s.listner.(*net.TCPListener)
			listener.SetDeadline(time.Now().Add(2 * time.Second))
			conn, err := listener.Accept()
			if oppErr, ok := err.(*net.OpError); ok && oppErr.Timeout() {
				continue
			}

			if err != nil {
				log.Fatal()
			}
			s.connections[s.totalConnnections] = conn
			s.totalConnnections += 1
			fmt.Println("Connect to memoryDB", s.totalConnnections)
			s.write(conn, "Welcome to memory DB")
			go s.handleConnection(conn)

		}

	}
}

func (s *Server) write(conn net.Conn, st string) {
	_, err := fmt.Fprintf(conn, "%s\n-> ", st)
	if err != nil {
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		query := scanner.Text()
		const invalidquery = "INVLALID QUERY"
		query = strings.Trim(query, " ")
		query = strings.ReplaceAll(query, "\n", "")
		tokens := strings.Split(query, " ")
		var res string
		switch {
		case tokens[0] == "set" && len(tokens) == 3:
			res = s.db.SetValue(tokens[1], tokens[2])
		case tokens[0] == "get" && len(tokens) == 2:
			res = s.db.GetValue(tokens[1])
		case tokens[0] == "del" && len(tokens) == 2:
			res = s.db.DeleteValue(tokens[1])
		default:
			res = invalidquery
		}
		s.write(conn, res)
	}
}

func (s *Server) Stop() {
	fmt.Println("Closing db connection")
	close(s.quit)
	<-s.exited

}
