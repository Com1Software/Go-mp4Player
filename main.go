package main

import (
	"bytes"
	"math/rand/v2"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fmt"
	"os"
)

type Tags struct {
	Tag string `dbase:"TAG"`
}

// ----------------------------------------------------------------
func main() {
	//	var openCmd *exec.Cmd
	fmt.Println("Go mp4 Player")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	exefile := ""
	exefilea := ""
	drive := "c"
	wdir := "/tunes/"
	//	tnfile := drive + ":" + wdir + "test.mp4"
	rfile := "./tmp.mp4"
	switch runtime.GOOS {
	case "windows":
		exefile = "c:/ffmpeg/bin/ffmpeg.exe"
		exefilea = "c:/ffmpeg/bin/ffprobe.exe"
		wdir = drive + ":/tunes/"

	case "linux":
		exefile = "ffmpeg"
		wdir = "/media/dave/Elements/dwhelper/"

	}

	subdir := true

	fmt.Println(exefile)
	fmt.Println(wdir)
	fmt.Println(subdir)
	switch {
	//-------------------------------------------------------------
	case len(os.Args) == 2:

		fmt.Println("Not")

		//-------------------------------------------------------------
	default:

		fmt.Println("Running....")

		fmt.Println("")

		// Read directory
		entries, err := os.ReadDir(wdir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		var playlist []string

		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".mp4") {
				fullpath := filepath.Join(wdir, e.Name())
				playlist = append(playlist, fullpath)
			}
		}

		if len(playlist) == 0 {
			fmt.Println("No mp4 files found.")
			return
		}

		min := 0
		max := len(playlist)
		matrix := make([][2]int, len(playlist))
		rtn := 0
		ttnfile := ""
		fmt.Println(matrix)
		for idx, tnfile := range playlist {

			fmt.Println(idx)
			fmt.Println(min + rand.IntN(max-min))
			rtn, matrix = updateMatrix(matrix, idx, max)
			ttnfile = playlist[rtn]
			// Delete old tmp.mp4
			if _, err := os.Stat(rfile); err == nil {
				fmt.Println("Deleting old tmp.mp4...")
				if err := os.Remove(rfile); err != nil {
					fmt.Println("Error deleting tmp.mp4:", err)
				}
			}

			// Copy with ffmpeg
			cmd := exec.Command(exefile, "-i", ttnfile, "-c", "copy", rfile)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Println("Converting:", tnfile)
			if err := cmd.Run(); err != nil {
				fmt.Printf("ffmpeg error: %v\n", err)
				return
			}

			openCmd := exec.Command("cmd", "/c", "start", "", "/wait", rfile)
			if err := openCmd.Run(); err != nil {
				fmt.Printf("Error opening file: %v\n", err)
			}

			// Optional: wait for user to close player before continuing
			// fmt.Println("Press ENTER for next video...")
			// fmt.Scanln()
			duration := getDuration(exefilea, tnfile)
			wait := time.Duration(duration) * time.Second
			wait += 10 * time.Second
			time.Sleep(wait)
		}

		if len(playlist) == 0 {
			fmt.Println("No mp4 files found.")
			return
		}

	}
}

func getDuration(ffprobePath, file string) int {
	cmd := exec.Command(ffprobePath,
		"-v", "quiet",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		file,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("Duration error:", err)
		return 0
	}

	raw := strings.TrimSpace(out.String())
	if raw == "" || raw == "N/A" {
		return 0
	}

	raw = strings.Split(raw, ".")[0]
	sec, _ := strconv.Atoi(raw)
	return sec
}

func updateMatrix(m [][2]int, idx int, max int) (int, [][2]int) {
	ii := 0
	for i := range m {
		m[i][0] = i
		m[i][1] = rand.IntN(len(m))
	}
	return ii, m
}
