package main

import (
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"fmt"
	"os"

	asciistring "github.com/Com1Software/Go-ASCII-String-Package"
)

type Tags struct {
	Tag string `dbase:"TAG"`
}

// ----------------------------------------------------------------
func main() {
	fmt.Println("Go mp4 Player")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	exefile := ""

	drive := "f"
	wdir := "/tmp/"
	switch runtime.GOOS {
	case "windows":
		exefile = "c:/ffmpeg/bin/ffmpeg.exe"
		wdir = drive + ":/dwhelper/"

	case "linux":
		exefile = "ffmpeg"
		wdir = "/media/dave/Elements/dwhelper/"

	}

	subdir := true

	fmt.Print(exefile)
	fmt.Print(wdir)
	fmt.Print(subdir)
	switch {
	//-------------------------------------------------------------
	case len(os.Args) == 2:

		fmt.Println("Not")

		//-------------------------------------------------------------
	default:

		fmt.Println("Running....")

		fmt.Println("")
	}
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func fixFileName(fileName string) string {
	newName := ""
	chr := ""
	ascval := 0

	for x := 0; x < len(fileName); x++ {
		chr = fileName[x : x+1]
		ascval = asciistring.StringToASCII(chr)
		switch {
		case ascval < 45:
		case ascval == 64:
		case ascval == 92:
		case ascval == 96:
		case ascval > 122:
		default:
			newName = newName + chr
		}
	}

	err := os.Rename(fileName, newName)
	if err != nil {
		fmt.Println("Error renaming file:", err)
	}
	return newName
}

func ValidFileType(fileExt string) bool {
	rtn := false
	switch {
	case fileExt == ".mp4":
		rtn = true
	}
	return rtn
}
