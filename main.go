package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Endianness int

const (
	BigEndian Endianness = iota
	LittleEndian
)

func PrintRectangle(rect image.Rectangle) {
	fmt.Printf("Bounds\n  Min: (%d,%d)\n  Max: (%d,%d)\n", rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}

func ReadImageFromFile(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()

	PrintRectangle(bounds)

	return m
}

func ConnectToXServer() *xgbutil.XUtil {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal("Failed to establish an X connection")
	}
	return X
}

func GetRootWindow(X *xgbutil.XUtil) *xwindow.Window {
	root := xwindow.New(X, X.RootWin())
	fmt.Printf("Root Window Id: %d\n", root.Id)
	return root
}

func ConvertImage(X *xgbutil.XUtil, img image.Image) *xgraphics.Image {
	ximg := xgraphics.NewConvert(X, img)
	if ximg == nil {
		log.Fatal("Failed to convert the image to an XImage")
	}
	return ximg
}

func SetPixmapProperty(X *xgb.Conn, root xproto.Window, propName string, pixmap xproto.Pixmap) {
	propAtom, err := xproto.InternAtom(X, false, uint16(len(propName)), propName).Reply()
	if err != nil {
		log.Fatal("Failed to get the %s property atom", propName)
	}

	typName := "PIXMAP"
	typAtom, err := xproto.InternAtom(X, false, uint16(len(typName)), typName).Reply()
	if err != nil {
		log.Fatal("Failed to get the %s type atom", typName)
	}

	data := Int32ToByteSlice(int32(pixmap), BigEndian)
	//data := Int32ToByteSlice(int32(pixmap), LittleEndian)

	for i, e := range data {
		fmt.Printf("Byte %d: %d\n", i, e)
	}

	xproto.ChangePropertyChecked(
		X, // Display *xgb.Conn
		xproto.PropModeReplace, // Mode byte
		root,          // Window Window
		propAtom.Atom, // Property Atom
		typAtom.Atom,  // Type Atom
		32,            // Format byte
		1,             // DataLen uint32
		data,          // Data []byte
	).Check()

	if err != nil {
		log.Fatal("Failed to change property: %s", err)
	}
}

func Int32ToByteSlice(input int32, endian Endianness) []byte {
	var byteA, byteB, byteC, byteD byte

	switch e := endian; e {
	case LittleEndian:
		byteD = byte(input & 0xff)
		byteC = byte((input >> 8) & 0xff)
		byteB = byte((input >> 16) & 0xff)
		byteA = byte((input >> 24) & 0xff)
	case BigEndian:
		byteA = byte(input & 0xff)
		byteB = byte((input >> 8) & 0xff)
		byteC = byte((input >> 16) & 0xff)
		byteD = byte((input >> 24) & 0xff)
	}

	return []byte{byteA, byteB, byteC, byteD}
}

func main() {
	dir := "/home/paul/go/src/github.com/pcewing/wpr2/data"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Failed to read files")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	fileInfo := files[rand.Intn(len(files))]

	filePath := path.Join(dir, fileInfo.Name())
	fmt.Printf("%s\n", filePath)

	img := ReadImageFromFile(filePath)

	X := ConnectToXServer()
	defer X.Conn().Close()

	root := GetRootWindow(X)

	ximg := ConvertImage(X, img)

	fmt.Printf("Pixmap Id: %d\n", ximg.Pixmap)

	err = ximg.CreatePixmap()
	if err != nil {
		log.Fatal("Failed to create pixmap")
	}

	fmt.Printf("Pixmap Id: %d\n", ximg.Pixmap)

	fmt.Printf("Drawing the image onto the pixmap\n")
	ximg.XDraw()

	SetPixmapProperty(X.Conn(), root.Id, "_XROOTPMAP_ID", ximg.Pixmap)
	SetPixmapProperty(X.Conn(), root.Id, "ESETROOT_PMAP_ID", ximg.Pixmap)
}
