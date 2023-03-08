package nets

import (
	"github.com/chaniqure/chan_tool/files"
	"testing"
)

func TestDownload(t *testing.T) {
	bytes, err := GetBytes("https://pic.biedoul.com/Uploads/Images/44/595a1ca497c9c.jpg")
	if err != nil {
		return
	}
	files.WriteBytesToFile("./1.jpg", bytes)
}
