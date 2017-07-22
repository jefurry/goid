package goid

import (
	"github.com/speps/go-hashids"
	"sync/atomic"
	"time"
)

/**
*
* 128bit
*
* --------------------------------------------------------------------------------------------------------------------------------
* |                  |                |                      |                  |                  |               |             |
* | 64 bit timestamp | 32 bit counter | 6 bit server room ID | 8 bit cluster ID | 9 bit machine ID | 3 bit work ID | 6 bit op ID |
* |                  |                |                      |                  |                  |               |             |
* --------------------------------------------------------------------------------------------------------------------------------
*
 */

const (
	VERSION            = "0.0.1"
	DEFAULT_ALPHABET   = hashids.DefaultAlphabet
	DEFAULT_MIN_LENGTH = 28
	DEFAULT_SALT       = `aC%eNsc8@!03Yb#69d`
	RAW_LEN            = 4 // 4 * 4 bytes
)

var (
	objectIDCounter = randInt()
)

type ID struct {
	alphabet                                         string
	salt                                             string
	minLength                                        int
	timestamp                                        int64
	counter                                          uint32
	serverRoomID, clusterID, machineID, workID, opID uint32
}

// New generates a globaly unique ID
func New() *ID {
	return &ID{
		alphabet:  DEFAULT_ALPHABET,
		salt:      DEFAULT_SALT,
		minLength: DEFAULT_MIN_LENGTH,
	}
}

// SetAlphabet set alphabet of hashids
func (id *ID) SetAlphabet(alphabet string) *ID {
	id.alphabet = alphabet

	return id
}

// SetSalt set salt as key secret of hashids
func (id *ID) SetSalt(salt string) *ID {
	id.salt = salt

	return id
}

// SetMinLength set minlength of hashids
func (id *ID) SetMinLength(minLength int) *ID {
	id.minLength = minLength

	return id
}

// Encode generates globally unique id
// serverroomID: [0, 64)
// clusterID: [0, 256)
// machineID: [0, 512)
// workID: [0, 8)
// opID: [0, 64)
func (id *ID) Encode(serverRoomID, clusterID, machineID, workID, opID uint32) (string, error) {
	nums64 := make([]int64, 0, RAW_LEN)

	// big endian
	timestamp := time.Now().UnixNano()
	c := atomic.AddUint32(&objectIDCounter, 1)

	nums64 = append(nums64, int64(timestamp>>32))
	nums64 = append(nums64, int64(timestamp&0xffffffff))
	nums64 = append(nums64, int64(c))

	last := ((serverRoomID & 0x3f) << 26) |
		((clusterID & 0xff) << 18) |
		((machineID & 0x1ff) << 9) |
		((workID & 0x7) << 6) |
		(opID & 0x3f)

	nums64 = append(nums64, int64(last))

	hd := id.getHD()
	s, err := hd.EncodeInt64(nums64)
	if err != nil {
		return "", err
	}

	id.timestamp = timestamp
	id.counter = c
	id.serverRoomID = serverRoomID
	id.clusterID = clusterID
	id.machineID = machineID
	id.workID = workID
	id.opID = opID

	return s, nil
}

// FromString reads an ID from its string representation
func (id *ID) FromString(s string) error {
	hd := id.getHD()
	nums64, err := hd.DecodeInt64WithError(s)
	if err != nil {
		return err
	}

	id.timestamp = int64((uint64(nums64[0]) << 32) | uint64(nums64[1]))
	id.counter = uint32(nums64[2])
	id.serverRoomID = (uint32(nums64[3]) >> 26) & 0x3f
	id.clusterID = (uint32(nums64[3]) >> 18) & 0xff
	id.machineID = (uint32(nums64[3]) >> 9) & 0x1ff
	id.workID = (uint32(nums64[3]) >> 6) & 0x7
	id.opID = uint32(nums64[3]) & 0x3f

	return nil
}

// Timestamp returns the timestamp part of the id.
func (id *ID) Timestamp() int64 {
	return id.timestamp
}

// Counter returns the counter part of the id.
func (id *ID) Counter() uint32 {
	return id.counter
}

// Serverroomid returns the serverroomid part of the id.
func (id *ID) ServerRoomID() uint32 {
	return id.serverRoomID
}

// Clusterid returns the clusterid part of the id.
func (id *ID) ClusterID() uint32 {
	return id.clusterID
}

// Machineid returns the machineid part of the id.
func (id *ID) MachineID() uint32 {
	return id.machineID
}

// WorkID returns the workID part of the id.
func (id *ID) WorkID() uint32 {
	return id.workID
}

// OpID returns the opID part of the id.
func (id *ID) OpID() uint32 {
	return id.opID
}

func (id *ID) getHD() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = id.salt
	hd.MinLength = id.minLength
	hd.Alphabet = id.alphabet

	return hashids.NewWithData(hd)
}
