/*
 * Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */

package jag

import (
	"bytes"
	"compress/bzip2"
	"encoding/binary"
	"io"
	"io/ioutil"
	"strings"

	"github.com/spkaeros/rscgo/pkg/log"
)

// TODO: Parse metadata into meaningful data structures for easier reading and to provide a better API

//Archive Represents a JAG archive, which is a bzip2 compressed file format to hold many files in one more easily.
type Archive struct {
	//FileCount How many files this JAG archive contains
	FileCount int
	//MetaData The meta data for each file.  4-byte int nameHash, 3-byte decompLen, 3-byte compLen
	MetaData []byte
	//FileData The raw, consecutive file data.
	FileData []byte
}

//decoder Returns an io.Reader that reads JAG archive file data and turns it into the raw, decompressed data that
// made it. In order to get the standard library to decode this strange antiquated file format, I had to remove the JAG
// archive header (2x3-byte ints, decompressed, then compressed length) then manually insert a BZ2 header
// ('B','Z','h','[1-9]', the last byte is the compression level, default 1) before the compressed payload.
func decoder(data []byte) io.Reader {
	return bzip2.NewReader(bytes.NewReader(append([]byte{'B', 'Z', 'h', '1'}, data[6:]...)))
}

//New Returns a new JAG archive, with the entry count, file metadata, and file data parsed to make reading the data much simpler.
func New(file string) *Archive {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		log.Warn("Problem occurred attempting to read the JAG archive:", err)
		return nil
	}
	buf, err := ioutil.ReadAll(decoder(fileData))
	if err != nil && !strings.HasSuffix(err.Error(), "continuation file") {
		log.Warn("Problem occurred attempting to decompress the JAG archive:", err)
		return nil
	}
	count := int(binary.BigEndian.Uint16(buf[:]))
	return &Archive{count, buf[2 : count*10+2], buf[count*10+2:]}
}
