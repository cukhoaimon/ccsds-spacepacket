package spacepacket

import "errors"

type SpacePacketCodec interface {
	Encode(packet SpacePacket) ([]byte, error)
	Decode(data []byte) (SpacePacket, error)
}

type DefaultSpacePacketCodec struct {
}

func (c DefaultSpacePacketCodec) Encode(packet SpacePacket) ([]byte, error) {
	if len(packet.Data) == 0 {
		return nil, errors.New("packet data field cannot be empty")
	}
	if len(packet.Data) > 65536 {
		return nil, errors.New("packet data exceeds maximum length of 65536 bytes")
	}
	packet.Header.DataLen = uint16(len(packet.Data) - 1)

	packetSizeInOctets := PacketPrimaryHeaderSizeInByte + packet.GetDataLen()
	data := make([]byte, packetSizeInOctets)

	// octet 1 = 3 bits version + 1 bit type + 1 bit secondary header flag + first 3 bits of APID
	data[0] = (packet.Header.Version << 5) | (packet.Header.Type << 4) | (packet.Header.SecondaryHeaderFlag << 3) | uint8(packet.Header.APID>>8)
	// octet 2 = remaining 8 byte of APID
	data[1] = uint8(packet.Header.APID)
	// octet 3 = 2 bits sequence flags + 6 bits packet sequence count
	data[2] = (packet.Header.SeqFlags << 6) | (uint8(packet.Header.SeqCounts >> 8))
	// octet 4 = remaining 8 byte of SeqCounts
	data[3] = uint8(packet.Header.SeqCounts)
	// octet 5 + 6 = packet data length
	data[4] = uint8(packet.Header.DataLen >> 8)
	data[5] = uint8(packet.Header.DataLen)

	for i := range packet.Data {
		data[i+6] = packet.Data[i]
	}

	return data, nil
}

func (c DefaultSpacePacketCodec) Decode(packet []byte) (SpacePacket, error) {
	header := SpacePacketHeader{}

	if len(packet) < 7 {
		return SpacePacket{}, errors.New("packet size must be at least 7 bytes")
	}

	header.Version = packet[0] >> 5
	header.Type = (packet[0] >> 4) & 0b1
	header.SecondaryHeaderFlag = (packet[0] >> 3) & 0b1
	header.APID |= uint16(packet[0]&0b111) << 8
	header.APID |= uint16(packet[1])
	header.SeqFlags = packet[2] >> 6
	header.SeqCounts |= uint16(packet[2]&0b111111) << 8
	header.SeqCounts |= uint16(packet[3])
	header.DataLen |= uint16(packet[4]) << 8
	header.DataLen |= uint16(packet[5])

	if len(packet) < int(6+header.DataLen+1) {
		return SpacePacket{}, errors.New("packet data length mismatch: not enough data")
	}

	data := packet[6 : 6+header.DataLen+1]
	return SpacePacket{header, data}, nil
}
