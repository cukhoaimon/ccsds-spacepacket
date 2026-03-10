package spacepacket

const (
	APIDIdle uint16 = 0x07FF // '11111111111'

	PacketPrimaryHeaderSizeInByte uint32 = 6 // 6 octets

	SeqFlagContinuation uint8 = 0x00
	SeqFlagFirst        uint8 = 0x01
	SeqFlagLast         uint8 = 0x02
	SeqFlagUnsegmented  uint8 = 0x03

	TMPacket uint8 = 0x00
	TCPacket uint8 = 0x01

	SecondaryHeaderFlagPresent uint8 = 0x01
)

type SpacePacketHeader struct {
	// PACKET VERSION NUMBER
	/*
		contain the binary encoded Packet Version Number
	*/
	Version uint8 // 3 bits

	// PACKET IDENTIFICATION
	/*
		'0' for telemetry, '1' for telecommand
	*/
	Type uint8 // 1 bit

	/*
		'1' if a packet Secondary header is present, '0' for not present or Idle packets.
	*/
	SecondaryHeaderFlag uint8 // 1 bit

	/*
		Application process identifier, unique for its naming domain, idle packet shall be '11111111111'
	*/
	APID uint16 // 11 bits

	// PACKET SEQUENCE CONTROL
	/*
		'00' for continuation segment of User Data
		'01' for the first segment of User Data
		'10' for the last segment of User Data
		'11' for unsegmented of User Data
	*/
	SeqFlags uint8 // 2 bits

	/*
		For telemetry packet, this field shall contains Packet Sequence Count
		For telecommand packet, this field shall contain either the Packet Sequence Count or Packet Name
	*/
	SeqCounts uint16 // 14 bits, seq count or name

	// PACKET DATA LENGTH
	/*
		Length in octets of the Packet Data Field.
		DataLen = (total number of Octets in the Packet Data Field) - 1
	*/
	DataLen uint16 // 16 bits
}

type SpacePacket struct {
	Header SpacePacketHeader
	Data   []byte
}

// GetDataLen return size of the actual data packet in octets
// In CCSDS specification, the DataLen = (total number of Octets in the Packet Data Field) - 1,
// to get the actual data we have to plus 1.
func (p SpacePacket) GetDataLen() uint32 {
	return uint32(p.Header.DataLen + 1)
}
