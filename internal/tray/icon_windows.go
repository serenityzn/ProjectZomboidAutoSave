//go:build windows

package tray

import (
	"bytes"
	"encoding/binary"

	"github.com/getlantern/systray"
)

func setIcon(icon []byte) {
	systray.SetIcon(pngToIco(icon))
}

// pngToIco wraps PNG bytes in a minimal ICO container.
// Windows Vista+ supports PNG images embedded directly inside ICO files.
func pngToIco(png []byte) []byte {
	var buf bytes.Buffer

	// ICONDIR header (6 bytes)
	binary.Write(&buf, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // type: 1 = icon
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // image count: 1

	// ICONDIRENTRY (16 bytes) — offset = 6 (header) + 16 (entry) = 22
	buf.WriteByte(0)                                                     // width  (0 = 256)
	buf.WriteByte(0)                                                     // height (0 = 256)
	buf.WriteByte(0)                                                     // color count
	buf.WriteByte(0)                                                     // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1))                  // planes
	binary.Write(&buf, binary.LittleEndian, uint16(32))                 // bit count
	binary.Write(&buf, binary.LittleEndian, uint32(len(png)))           // image size
	binary.Write(&buf, binary.LittleEndian, uint32(22))                 // image offset

	// PNG image data
	buf.Write(png)

	return buf.Bytes()
}
