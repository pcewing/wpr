package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"

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

	for i, e := range data {
		fmt.Printf("Byte %d: %d\n", i, e)
	}

	xproto.ChangePropertyChecked(
		/* Display *xgb.Conn */ X,
		/* Mode byte */ xproto.PropModeReplace,
		/* Window Window */ root,
		/* Property Atom */ propAtom.Atom,
		/* Type Atom */ typAtom.Atom,
		/* Format byte */ 32,
		/* DataLen uint32 */ 1,
		/* Data []byte */ data,
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

func SetBackgroundX11(img image.Image) {
	X := ConnectToXServer()
	defer X.Conn().Close()

	root := GetRootWindow(X)

	ximg := ConvertImage(X, img)

	fmt.Printf("Pixmap Id: %d\n", ximg.Pixmap)

	err := ximg.CreatePixmap()
	if err != nil {
		log.Fatal("Failed to create pixmap")
	}

	fmt.Printf("Pixmap Id: %d\n", ximg.Pixmap)

	fmt.Printf("Drawing the image onto the pixmap\n")
	ximg.XDraw()

	SetPixmapProperty(X.Conn(), root.Id, "_XROOTPMAP_ID", ximg.Pixmap)
	SetPixmapProperty(X.Conn(), root.Id, "ESETROOT_PMAP_ID", ximg.Pixmap)
}
