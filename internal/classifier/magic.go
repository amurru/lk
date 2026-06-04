package classifier

import "github.com/gabriel-vasile/mimetype"

func mimetypeDetect(data []byte) string {
	return mimetype.Detect(data).String()
}
