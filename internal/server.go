package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	listner net.Listener
	db      Database
}

func NewServer() *Server {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("DB Error")
	}

	ser := &Server{
		listner: ln,
		db:      NewDatabase(),
	}
	go ser.Serve()
	return ser
}

func (s *Server) Serve() {
	for {
		conn, err := s.listner.Accept()
		if err != nil {
			log.Fatal()
		}
		s.write(conn, "Welcome to memory DB")
		go s.handleConnection(conn)
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
		querytype := tokens[0]
		var res string
		switch strings.ToUpper(querytype) {
		case "SET":
			if (len(tokens)) != 3 {
				res = invalidquery
			} else {

				res = s.db.SetValue(tokens[1], tokens[2])
			}
		case "GET":
			if len(tokens) != 2 {
				res = invalidquery
			} else {

				res = s.db.GetValue(tokens[1])
			}
		case "DEL":
			if len(tokens) != 2 {
				res = invalidquery
			} else {
				res = s.db.DeleteValue(tokens[1])
			}
		default:
			res = invalidquery
		}
		s.write(conn, res)
	}
}

func (s *Server) Stop() {
	fmt.Println("Stop thes server")
}
