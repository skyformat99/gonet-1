package main

import (
	"errors"
	"syscall"
	//"fmt"
)

type Network_Writer struct {
	fd       int
	sockAddr syscall.Sockaddr
}

func NewNetwork_Writer() (*Network_Writer, error) {
	fd, err := syscall.Socket(AF_PACKET, SOCK_RAW, HTONS_ETH_P_ALL)
	if err != nil {
		return nil, errors.New("Write's socket failed")
	}

	addr := &syscall.SockaddrLinklayer{
		// Family is automatically set to AF_PACKET
		Protocol: ETHERTYPE_IP, // should be inherited anyway
		Addr:     myMACAddr,    // sending to myself
		Halen:    ETH_ALEN,     // may not be correct
		Ifindex:  MyIfIndex,    // TODO: don't hard code this... fix it later
	}

	/*err = syscall.Sendto(fd, []byte{0x08, 0x00, 0x27, 0x9e, 0x29, 0x63, 0x08, 0x00, 0x27, 0x9e, 0x29, 0x63, 0x08, 0x00}, 0, addr) //Random bytes
	  if err != nil {
	      fmt.Println("ERROR returned by syscall.Sendto", err)
	  } else {
	      fmt.Println("Sent the test packet")
	  }*/

	return &Network_Writer{
		fd:       fd,
		sockAddr: addr,
	}, nil
}

func (nw *Network_Writer) write(data []byte) error {
	// build the ethernet header
	/*etherHead :=  append(append(
	    myMACSlice, // dst MAC
	    myMACSlice...), // src MAC
	    0x08, 0x00, // ethertype (IP)
	)*/
	// TODO: decide this dynamically
	etherHead := []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0,
	}
	//fmt.Println("My header:", etherHead)

	// add on the ethernet header
	newPacket := append(etherHead, data...)
	//fmt.Println("Full Packet with ethernet header:", newPacket)

	return syscall.Sendto(nw.fd, newPacket, 0, nw.sockAddr)
}

func (nw *Network_Writer) close() error {
	return syscall.Close(nw.fd)
}