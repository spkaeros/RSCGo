package rsdatabase

import (
	"fmt"
	"io/ioutil"
)

//EntrySize Size of a single entry, in bytes.
const EntrySize = 140

var NilEntry = RSDEntry{-1, 0, "", -1, -1}

//RSDatabase Represents a single player profile database.
type RSDatabase struct {
	buffer   []byte
	position int
}

type RSDEntry struct {
	Index int
	// base37
	Username uint64
	// SHAKE256-64, len=128
	Password string
	X, Y     int
}

//New Returns a new RSDatabase reference populated with the RSDatabase file at path.
func New(path string) *RSDatabase {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Problem reading RSD file:", err)
		return nil
	}
	if len(buf)%EntrySize != 0 {
		fmt.Println("WARNING: RSDatabase buffer length incorrect.  Not evenly divisible by entry size.")
	}
	return &RSDatabase{buffer: buf}
}

//Size Returns the number of entries in this database.
func (db *RSDatabase) Size() int {
	return len(db.buffer) / EntrySize
}

//ReadInt Reads the next integer from the buffer and return it.
func (db *RSDatabase) ReadInt() int {
	defer func() {
		db.position += 4
	}()
	return int(int(db.buffer[db.position])<<24 | int(db.buffer[db.position+1])<<16 | int(db.buffer[db.position+2])<<8 | int(db.buffer[db.position+3]))
}

//ReadLong Reads the next long integer from the buffer and return it.
func (db *RSDatabase) ReadLong() uint64 {
	defer func() {
		db.position += 4
	}()
	return uint64(int(db.buffer[db.position])<<56 | int(db.buffer[db.position+1])<<48 | int(db.buffer[db.position+2])<<40 | int(db.buffer[db.position+3])<<32 | int(db.buffer[db.position+4])<<24 | int(db.buffer[db.position+5])<<16 | int(db.buffer[db.position+6])<<8 | int(db.buffer[db.position+7]))
}

//ReadShort Return the next short integer from the buffer and return it.
func (db *RSDatabase) ReadShort() int {
	defer func() {
		db.position += 2
	}()
	return int(int(db.buffer[db.position])<<8 | int(db.buffer[db.position+1]))
}

//ReadPass Returns next 128 bytes as a string representation of a password hash.
func (db *RSDatabase) ReadPass() string {
	defer func() {
		db.position += 128
	}()
	return string(db.buffer[db.position : db.position+128])
}

//Reset Resets the position offset in the buffer to 0
func (db *RSDatabase) Reset() {
	db.position = 0
}

//HasNext Returns true if there's any more entries in the buffer, otherwise returns false.
func (db *RSDatabase) HasNext() bool {
	return db.position+EntrySize < len(db.buffer)
}

//Next Returns the next entry in the database.  If there isn't any entries left, returns NilEntry.
func (db *RSDatabase) Next() RSDEntry {
	if !db.HasNext() {
		fmt.Println("Buffer overflow: ")
		return NilEntry
	}
	return RSDEntry{db.ReadInt(), db.ReadLong(), db.ReadPass(), db.ReadShort(), db.ReadShort()}
}

//Find Returns the player in the buffer with the specified index.  If no player is found, returns a RSDEntry filled with -1s
func (db *RSDatabase) Find(index int) RSDEntry {
	defer db.Reset()
	db.Reset()
	for db.HasNext() {
		entry := db.Next()
		if entry.Index == index {
			return entry
		}
	}
	/*	for db.Reset(); db.position < len(db.buffer); {
			nextIndex := db.ReadInt()
			if nextIndex == index {
				return RSDEntry{Index: nextIndex, X: db.ReadShort(), Y: db.ReadShort()}
			} else {
				db.position += 4
			}
		}
	*/
	return NilEntry
}

//Save Saves the specified information to the player profile at the specified index
func (db *RSDatabase) Save(index int, userHash uint64, passHash string, x, y int) {
	if index > db.Size() {
		fmt.Println("rsd.Save: Index out of bounds")
		return
	}
	pos := index * EntrySize
	// TODO: Omit this?
	db.buffer[pos] = byte(index >> 24)
	db.buffer[pos+1] = byte(index >> 16)
	db.buffer[pos+2] = byte(index >> 8)
	db.buffer[pos+3] = byte(index)

	db.buffer[pos+4] = byte(userHash >> 56)
	db.buffer[pos+5] = byte(userHash >> 48)
	db.buffer[pos+6] = byte(userHash >> 40)
	db.buffer[pos+7] = byte(userHash >> 32)
	db.buffer[pos+8] = byte(userHash >> 24)
	db.buffer[pos+9] = byte(userHash >> 16)
	db.buffer[pos+10] = byte(userHash >> 8)
	db.buffer[pos+11] = byte(userHash)

	for i, c := range passHash {
		db.buffer[pos+12+i] = byte(c)
	}

	db.buffer[pos+140] = byte(x >> 8)
	db.buffer[pos+141] = byte(x)

	db.buffer[pos+142] = byte(y >> 8)
	db.buffer[pos+143] = byte(y)
}
