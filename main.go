package main

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"bufio"
	"io"
	"regexp"
	"errors"
)

// todo: channel-fy this. Like Unix pipes. Communicate data through channels:
// split paths into directory names | list directory contents | 
// how do I make sure this is deterministically ordered? I have to block/buffer at the end until the sequence of things arrive
// will be an array with a linked list e.g.
// [ { (head node) value is max, next pointer is to highest value or null, prev is to tail }
//   { (tail node) value is min, next pointer is to head, prev is to lowest value or null }
//   { 3, 5, 0 }
//   { 1, 1, 5 }
//   { 2, 4, 3 }
// ]
// which is  sequence of three directories having come in, out of order, coming
// as "3" then "1" then "2". We'll also store an index into a buffer array that
// we'll set according the the segment that came in.
func main() {
	path := os.Getenv("PATH")
	paths := strings.SplitSeq(path, string(os.PathListSeparator))

	for dirpath := range paths {
		dstat, err := os.Stat(dirpath)

		if os.IsNotExist(err) {
			// part of PATH, but doesn't exist
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s: os statdir(%s)\n", err, dirpath)
			os.Exit(1)
		}

		if !dstat.IsDir() {
			// PATH should only contain path prefixes, aka directories
			continue
		}

		matchCount := 0

		entries, err := os.ReadDir(dirpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s: os readdir(%s)\n", err, dirpath)
			os.Exit(1)
		}

		for _, entry := range entries {
			if entry.Type().IsRegular() {
				fpath := filepath.Join(dirpath, entry.Name())

				fstat, err := os.Stat(fpath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "os stat(%s)\n", fpath)
					os.Exit(1)
				}

				if fstat.Mode() & 0111 == 0 {
					// not an executable
					continue
				}

				buffer := make([]byte, 256)

				file, err := os.Open(fpath)

				if errors.Is(err, os.ErrPermission) {
					// not allowed to open file
					continue
				}

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error %s\n", os.ErrPermission)
					fmt.Fprintf(os.Stderr, "Error %s: os open(%s)\n", err, fpath)
					os.Exit(1)
				}

				defer file.Close()

				reader := bufio.NewReader(file)

				bytesRead, err := reader.Read(buffer)

				if err == io.EOF {
					// empty file
					continue
				}

				if err != nil {
					fmt.Fprintf(os.Stderr, "Error %s: buf read(%s) #%d\n", err, fpath, bytesRead)
					os.Exit(1)
				}

				if bytesRead >= 2 && buffer[0] == '#' && buffer[1] == '!' {
					match, err := regexp.Match(`#!.*?.use.*?\n`, buffer)

					if err != nil {
						fmt.Fprintf(os.Stderr, "Error %s: regexp match\n", err)
					}

					if match {
						if matchCount > 0 {
							fmt.Println("")
						}

						matchCount += 1

						fmt.Printf("%s\t(%s)\n", filepath.Base(fpath), fpath)

						// Now I want to parse the usage script out of it.
						file.Seek(0, 0)
						scanner := bufio.NewScanner(file)
						for scanner.Scan() {
							line := scanner.Text()
							// discard shebang line
							if strings.HasPrefix(line, "#!") {
								continue
							} else if strings.HasPrefix(line, "#") {
								fmt.Println(line)
							} else {
								break
							}
						}

					}
				}
			}
		}
	}
}
