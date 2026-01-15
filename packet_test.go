package raknet

import (
	"fmt"
	"testing"
)

func TestPacket(t *testing.T) {
	pk := &packet{
		reliability: reliabilityReliableOrdered,
	}
	data := []byte{0x60, 0x0, 0x54, 0x2, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x60, 0x2, 0x1e, 0x40}
	n, err := pk.read(data)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err)
	}
	if n != 21 {
		t.Fatalf("Packet read wrong size: expected %d got %d", 11, n)
	}
	fmt.Println(n)
}
