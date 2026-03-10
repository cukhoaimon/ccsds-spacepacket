package spacepacket

import (
	"bytes"
	"testing"
)

func TestDefaultSpacePacketCodec_EncodeDecode(t *testing.T) {
	codec := DefaultSpacePacketCodec{}
	
	packet := SpacePacket{
		Header: SpacePacketHeader{
			Version:             0,
			Type:                TMPacket,
			SecondaryHeaderFlag: 0,
			APID:                0x567,
			SeqFlags:            SeqFlagUnsegmented,
			SeqCounts:           42,
		},
		Data: []byte{0xDE, 0xAD, 0xBE, 0xEF},
	}
	
	encoded, err := codec.Encode(packet)
	if err != nil {
		t.Fatalf("unexpected error during encode: %v", err)
	}
	
	// Check the length of encoded data
	if len(encoded) != 6 + 4 {
		t.Errorf("expected length 10, got %d", len(encoded))
	}
	
	// Test decoding
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("unexpected error during decode: %v", err)
	}
	
	if decoded.Header.APID != packet.Header.APID {
		t.Errorf("expected APID %d, got %d", packet.Header.APID, decoded.Header.APID)
	}
	
	if decoded.Header.DataLen != 3 {
		t.Errorf("expected DataLen 3, got %d", decoded.Header.DataLen)
	}
	
	if !bytes.Equal(decoded.Data, packet.Data) {
		t.Errorf("expected Data %v, got %v", packet.Data, decoded.Data)
	}
}

func TestDefaultSpacePacketCodec_Decode_BoundsCheck(t *testing.T) {
	codec := DefaultSpacePacketCodec{}
	
	// Create a header defining DataLen = 100 (101 bytes), but provide only 2 bytes
	encoded := []byte{
		0x00, 0x01, 0x00, 0x01, 0x00, 100, // Header
		0xFF, 0xFF, // Just 2 bytes of data
	}
	
	_, err := codec.Decode(encoded)
	if err == nil {
		t.Errorf("expected error due to insufficient data length, got none")
	} else if err.Error() != "packet data length mismatch: not enough data" {
		t.Errorf("got unexpected error: %v", err)
	}
}
