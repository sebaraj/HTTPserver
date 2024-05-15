package socket

import (
	"syscall"
)

type Socket struct {
	fd       int
	sockaddr syscall.SockaddrInet4
	backlog  int
}

func (sock *Socket) CreateSocket(dom int, serv int, protoc int, port uint16, addr uint32, bcklog int) error {
	var err error

	// socket address/port
	sock.sockaddr.Addr = [4]byte{byte(addr >> 24 & 0xFF), byte(addr >> 16 & 0xFF), byte(addr >> 8 & 0xFF), byte(addr & 0xFF)}
	sock.sockaddr.Port = int(port)
	sock.backlog = bcklog

	// domain needs to be AF_INET, not AF_INET6 bc Socket cant handle https://cs.opensource.google/go/go/+/refs/tags/go1.22.3:src/syscall/syscall_unix.go;l=497
	// establish connection to socket (fd)
	sock.fd, err = syscall.Socket(dom, serv, protoc)

	return err

}

func (sock *Socket) BindSocket() error {
	return syscall.Bind(sock.fd, &sock.sockaddr)
}

func (sock *Socket) ConnectSocket() error {
	return syscall.Connect(sock.fd, &sock.sockaddr)
}

func (sock *Socket) ListenSocket() error {
	return syscall.Listen(sock.fd, sock.backlog)
}

func (sock *Socket) AcceptSocket() (int, syscall.Sockaddr, error) {
	return syscall.Accept(sock.fd)
}

func (sock *Socket) CloseSocket() error {
	return syscall.Close(sock.fd)
}

/*
func (sock *Socket) GetFD() int {
	return sock.fd
}

func (sock *Socket) GetSockaddr() syscall.SockaddrInet4 {
	return sock.sockaddr
}
*/
