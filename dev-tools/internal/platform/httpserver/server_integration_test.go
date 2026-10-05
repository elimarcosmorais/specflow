//go:build integration

package httpserver

import (
	"errors"
	"net"
	"testing"
)

// CT-21
func TestListen_PortInUse(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	port := occupied.Addr().(*net.TCPAddr).Port

	ln, err := Listen(port)

	if ln != nil {
		ln.Close()
		t.Fatal("Listen devolveu um listener numa porta ocupada")
	}
	if !errors.Is(err, ErrPortInUse) {
		t.Errorf("err = %v, quer errors.Is(err, ErrPortInUse)", err)
	}
}

func TestListen_OtherFailureIsNotPortInUse(t *testing.T) {
	ln, err := Listen(-1)

	if ln != nil {
		ln.Close()
		t.Fatal("Listen devolveu um listener para a porta -1")
	}
	if err == nil || errors.Is(err, ErrPortInUse) {
		t.Errorf("err = %v, quer erro que não seja ErrPortInUse", err)
	}
}

// CT-22
func TestListen_LoopbackOnly(t *testing.T) {
	ln, err := Listen(0)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("Addr() = %T, quer *net.TCPAddr", ln.Addr())
	}
	if got := addr.IP.String(); got != "127.0.0.1" {
		t.Errorf("IP = %s, quer 127.0.0.1", got)
	}
	if addr.Port == 0 {
		t.Error("porta efetiva = 0")
	}
}
