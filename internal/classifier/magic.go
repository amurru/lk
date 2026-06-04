package classifier

import (
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func mimetypeDetect(data []byte) string {
	return mimetype.Detect(data).String()
}

// iso9660Signature is the standard identifier at offset 1 in the primary
// volume descriptor (sector 16 = byte offset 32768).
const iso9660Signature = "CD001"

// detectISO9660 checks for the ISO 9660 signature at the well-known offset.
// Returns true if the file is an ISO 9660 image.
func detectISO9660(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, len(iso9660Signature))
	n, err := f.ReadAt(buf, 32769)
	if err != nil || n < len(iso9660Signature) {
		return false
	}
	return string(buf) == iso9660Signature
}
