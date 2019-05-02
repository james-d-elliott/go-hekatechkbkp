package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//Executable is a string from ldflags, should be the executable name
var Executable string

//Version is a string from ldflags, should be the describe --tag --always value
var Version string

//Commit is a string from ldflags, should be the git commit in the long form
var Commit string

//Build is a string from ldflags, should be the travis build number or nothing
var Build string

func main() {
	//Start Time
	start := time.Now()

	//Pure Declarations
	var logFile *os.File
	var logOutput io.Writer
	var err error
	var chunkSum []byte

	//Declare Set
	chunkHasher := sha1.New()
	chunksValid, chunksInvalid, filesValid, filesInvalid, chunksize := 0, 0, 0, 0, 1024*4096

	//Regexes
	validSumsFile, _ := regexp.Compile("(?i).*?.sha1sums$")
	validSHA1, _ := regexp.Compile("(?i)^[0-9a-f]{40}$")

	//Flags for CLI
	var logging, logConsole bool
	var directory, logFileName string
	var debug int

	//Setup and Parse Flags
	flag.BoolVar(&logging, "log", true, "turns the log file on and off")
	flag.StringVar(&logFileName, "log-file", "verification.log", "changes the log file name")
	flag.BoolVar(&logConsole, "log-console", true, "turns the console logging on or off")
	flag.IntVar(&debug, "debug", 0, "sets the debug level")
	flag.StringVar(&directory, "dir", "./", "set the directory for the rawnand files")
	flag.Parse()

	//Fix dir if it doesn't end with a forward slash or contains windows style dirs (we use filepath.FromSlash() to generate our OS independent slashes)
	directory = strings.Replace(directory, "\\", "/", -1)
	if directory[len(directory)-1:] != "/" {
		directory = directory + "/"
	}

	//Logging Setup
	if logging {
		logFile, err = os.OpenFile(filepath.FromSlash(directory+"verification.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		check(err, true)
	}
	if logConsole && logging {
		logOutput = io.MultiWriter(logFile, os.Stdout)
	} else if logConsole {
		logOutput = io.MultiWriter(os.Stdout)
	} else if logging {
		logOutput = io.MultiWriter(logFile)
	}

	if logConsole || logging {
		log.SetOutput(logOutput)
	} else {
		log.SetFlags(0)
	}
	log.Printf("Logging started for %s %s commit %s build %s", Executable, Version, Commit, Build)

	//Read the list of files in the curent directory
	files, err := ioutil.ReadDir(directory)
	check(err, true)

	//Iterate files finding .sha1sums files
	for _, f := range files {
		fileNameSums := f.Name()
		if validSumsFile.MatchString(fileNameSums) {
			fileName := fileNameSums[:len(fileNameSums)-9]
			fileInfo, err := os.Stat(filepath.FromSlash(directory + fileName))
			if os.IsNotExist(err) {
				log.Printf("MISSING: found checksum file but the expected file %s does not exist", fileName)
			} else if err == nil {
				startFile := time.Now()
				chunksValidFile := 0
				chunksInvalidFile := 0

				var chunksInvalidSlice []int
				file, err := os.OpenFile(filepath.FromSlash(directory+fileName), os.O_RDONLY, 0666)
				check(err, true)
				defer file.Close()

				fileSums, err := os.OpenFile(filepath.FromSlash(directory+fileNameSums), os.O_RDONLY, 0666)
				check(err, true)
				defer fileSums.Close()

				//Read sums into string Array
				var sums []string
				scanner := bufio.NewScanner(fileSums)
				for scanner.Scan() {
					if validSHA1.MatchString(scanner.Text()) {
						if debug >= 1 {
							log.Printf("VERBOSE: found valid sha1 checksum %s", scanner.Text())
						}
						sums = append(sums, strings.ToLower(scanner.Text()))
					} else {
						wstr := removeWhitespace(scanner.Text())
						if strings.HasPrefix(wstr, "#chunksize:") {
							chunksize, err = strconv.Atoi(wstr[11:])
							check(err, false)
							if debug >= 1 {
								log.Printf("VERBOSE: chunk size set by sum file")
							}
						} else if debug >= 1 {
							log.Printf("VERBOSE: found invalid sha1 checksum %s", scanner.Text())
						}
					}
				}

				quotient, remainder := divmod(fileInfo.Size(), int64(chunksize))
				chunks := quotient
				if remainder != 0 {
					chunks++
				}
				if int64(len(sums)) != chunks {
					log.Printf("INVALID: file %s does not have the correct number of checksums; has %d should have %d", file.Name(), len(sums), chunks)
					continue
				}

				//Create the reading buffer and get ready to read
				log.Printf("CHECKING: file %s is being checked against %d checksums", fileName, len(sums))
				buf := make([]byte, chunksize)

				for chunk := 0; true; chunk++ {
					n, err := file.Read(buf)
					//break the loop if the the end of the file has been reached
					if err == io.EOF {
						break
					}
					check(err, true)

					//hash the chunk
					chunkHasher.Reset()
					chunkHasher.Write(buf[:n])
					chunkSum = chunkHasher.Sum(nil)

					if hex.EncodeToString(chunkSum) == sums[chunk] {
						if debug >= 2 {
							log.Printf("VALID: file %s chunk %d with checksum %x", fileName, chunk, chunkSum)
						}
						chunksValid++
						chunksValidFile++
					} else {
						if debug >= 2 {
							log.Printf("INVALID: file %s chunk %d with checksum %x, should be %s", fileName, chunk, chunkSum, sums[chunk])
						}
						chunksInvalid++
						chunksInvalidFile++
						chunksInvalidSlice = append(chunksInvalidSlice, chunk)
					}
				}
				if chunksInvalidFile == 0 {
					log.Printf("VALID: file %s checked completely (valid chunks: %d, invalid chunks: 0, time: %s)", fileName, chunksValidFile, time.Since(startFile))
					filesValid++
				} else {
					log.Printf("INVALID: file %s checked completely (valid chunks: %d, invalid chunks: %d, time: %s)", fileName, chunksValidFile, chunksInvalidFile, time.Since(startFile))
					filesInvalid++
					fileInvalid, err := os.OpenFile(filepath.FromSlash(directory+fileName+".invalid_chunks"), os.O_RDWR|os.O_CREATE, 0666)
					check(err, true)
					for _, c := range chunksInvalidSlice {
						fileInvalid.WriteString(strconv.Itoa(c) + "\n")
					}
				}
			} else {
				check(err, true)
			}
		}
	}
	if filesInvalid == 0 && filesValid != 0 {
		log.Printf("VALID: all files checked compeletely (valid files: %d, invalid files: 0, valid chunks: %d, invalid chunks: 0, total time: %s)", filesValid, chunksValid, time.Since(start))
	} else if filesValid == 0 {
		log.Printf("INVALID: no files found to verifiy")
	} else {
		log.Printf("INVALID: all files checked compeletely (valid files: %d, invalid files: %d, valid chunks: %d, invalid chunks: %d, total time: %s)", filesValid, filesInvalid, chunksValid, chunksInvalid, time.Since(start))
	}
}
