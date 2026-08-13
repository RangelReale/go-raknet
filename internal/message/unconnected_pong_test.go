package message

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// A Videoke VSK 5.0 (RakNet protocol version 6) replies to an unconnected ping
// with the offline ping response appended raw, exactly as stock RakNet's
// RakPeer::OnUnconnectedPing does — no uint16 length prefix. Captured live from
// 192.168.68.63:35379.
func TestUnconnectedPongUnmarshalRawPayload(t *testing.T) {
	const motd = "pxP24QU5JH0tCuLmHlHHaGiHe4qPS2UneDQBSJUW2p6z5Ja3sZUMeIENIwnAsH/rwZhw4kXI\n" +
		"xc1MWBPHOCPRNA==\n"

	var b []byte
	b = binary.BigEndian.AppendUint64(b, 10000)     // ping time echoed back
	b = binary.BigEndian.AppendUint64(b, 748559125) // server GUID 0x2C9D2D15
	b = append(b, unconnectedMessageSequence[:]...)
	b = append(b, motd...)

	pk := &UnconnectedPong{}
	if err := pk.UnmarshalBinary(b); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}
	if pk.PingTime != 10000 {
		t.Errorf("PingTime = %d, want 10000", pk.PingTime)
	}
	if pk.ServerGUID != 748559125 {
		t.Errorf("ServerGUID = %#x, want 0x2c9d2d15", pk.ServerGUID)
	}
	// The first two bytes ("px") are a length prefix only by coincidence of
	// position; reading them as one yields 0x7078 and truncates the payload.
	if string(pk.Data) != motd {
		t.Errorf("Data  = %q\nwant   = %q", pk.Data, motd)
	}
}

// Minecraft-style servers do length-prefix the MOTD, and that must keep working.
func TestUnconnectedPongUnmarshalLengthPrefixedPayload(t *testing.T) {
	motd := []byte("MCPE;Dedicated Server;800;1.21.0;0;10")

	var b []byte
	b = binary.BigEndian.AppendUint64(b, 42)
	b = binary.BigEndian.AppendUint64(b, 7)
	b = append(b, unconnectedMessageSequence[:]...)
	b = binary.BigEndian.AppendUint16(b, uint16(len(motd)))
	b = append(b, motd...)

	pk := &UnconnectedPong{}
	if err := pk.UnmarshalBinary(b); err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}
	if !bytes.Equal(pk.Data, motd) {
		t.Errorf("Data = %q, want %q", pk.Data, motd)
	}
}
