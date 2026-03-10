# ccsds-spacepacket

A Go library providing an implementation of the Space Packet Protocol (SPP) according to the Consultative Committee for Space Data Systems (CCSDS) standard. This library allows you to easily encode and decode space telemetry and telecommand packets.

## Installation

To use this library in your Go project, use `go get`:

```bash
go get github.com/cukhoaimon/ccsds-spacepacket
```

## Usage

The library provides a main `SpacePacket` struct that encapsulates both the header information and the packet data payload, as well as a `SpacePacketCodec` interface for encoding and decoding these packets.

### Encoding a Space Packet

You can construct a `SpacePacket` and encode it into a byte slice using the `DefaultSpacePacketCodec`. 

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/cukhoaimon/ccsds-spacepacket"
)

func main() {
	codec := spacepacket.DefaultSpacePacketCodec{}
	
	packet := spacepacket.SpacePacket{
		Header: spacepacket.SpacePacketHeader{
			Version:             0,
			Type:                spacepacket.TMPacket, // Telemetry Packet
			SecondaryHeaderFlag: 0,
			APID:                0x567,
			SeqFlags:            spacepacket.SeqFlagUnsegmented,
			SeqCounts:           42,
		},
		Data: []byte{0xDE, 0xAD, 0xBE, 0xEF},
	}
	
	encoded, err := codec.Encode(packet)
	if err != nil {
		log.Fatalf("Failed to encode packet: %v", err)
	}
	
	fmt.Printf("Encoded Packet (%d bytes): %x\n", len(encoded), encoded)
}
```

### Decoding a Space Packet

To decode a received byte slice back into a `SpacePacket`:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/cukhoaimon/ccsds-spacepacket"
)

func main() {
	codec := spacepacket.DefaultSpacePacketCodec{}
	
	encodedData := []byte{ /* ... bytes received ... */ }
	
	decodedPacket, err := codec.Decode(encodedData)
	if err != nil {
		log.Fatalf("Failed to decode packet: %v", err)
	}
	
	fmt.Printf("Decoded APID: %d\n", decodedPacket.Header.APID)
	fmt.Printf("Decoded Sequence Count: %d\n", decodedPacket.Header.SeqCounts)
	fmt.Printf("Decoded Data Length: %d octets\n", decodedPacket.GetDataLen())
	fmt.Printf("Decoded Data Payload: %x\n", decodedPacket.Data)
}
```

## Features

- Complete support for the 6-byte CCSDS Space Packet Primary Header format.
- Constants for flags, definitions, and packet types (Telemetry and Telecommand).
- Safely encodes and decodes packets with comprehensive length checks to prevent index out of bounds errors.

## Future Enhancements
The following features are planned to be added to the library:
- **Framer**: Add standard synchronization markers (e.g., ASM) and CADU framing capabilities.
- **Stream Processing**: Support for parsing continuous streams of data containing consecutive space packets to reliably extract packet boundaries. 
- **Secondary Header Handling**: Add standard mechanisms for decoding and encoding the secondary header data fields like telemetry timestamps and telecommand functions.

## Testing

To run the test suite, ensure you are in the project root directory and execute:

```bash
go test ./...
```
