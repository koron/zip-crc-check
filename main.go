package main

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func calcCRC32(zf *zip.File) (uint32, error) {
	r, err := zf.Open()
	if err != nil {
		return 0, err
	}
	defer r.Close()
	h := crc32.NewIEEE()
	if _, err := io.Copy(h, r); err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func check(name string) (bool, error) {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return false, err
	}
	defer zr.Close()
	result := true
	for _, zf := range zr.File {
		if zf.Mode().IsDir() {
			continue
		}
		v, err := calcCRC32(zf)
		if err != nil {
			return false, err
		}
		if v != zf.CRC32 {
			if result {
				fmt.Printf("%s NG\n", name)
			}
			result = false
			fmt.Printf("  %s CRC32 mismatch (expected:%08x actual:%08x)\n",
				zf.Name, zf.CRC32, v)
		}
	}
	return result, nil
}

func main() {
	for _, name := range os.Args[1:] {
		ok, err := check(name)
		if err != nil {
			fmt.Printf("%s ERROR %s\n", name, err)
			continue
		}
		if !ok {
			continue
		}
		fmt.Printf("%s OK\n", name)
	}
}
