package main

import (
	"os"
	"encoding/csv"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"github.com/pshiko/catle"
	"fmt"
	"io"
	"bufio"
	"strconv"
	"time"
)

const usage = `[Usage]
	csvhead <csvpath>
[Options]
	-nh:  without header
	-t: set '\t' as delimiter
	-s: set white space as delimiter
	-n <number>: skip <number> row
`

func ExtractOption(args []string, option string) ([]string, bool) {
	for i, arg := range args {
		if arg == option {
			return append(args[:i], args[i+1:]...), true
		}
	}
	return args, false
}

func main() {
	noheader := false
	delim := ','
	skip := 0
	args := os.Args
	// Read option
	time.Sleep(time.Second)
	for i, arg := range args {
		if arg == "-n" {
			if i == len(args) - 1 {
				log.Fatalln("-n option should set number.")
			}
			if n, err := strconv.Atoi(args[i+1]); err != nil {
				log.Fatalf("-n option error: %v", err)
			} else {
				skip = n
				if i < len(args) - 2 {
					args = append(args[:i], args[i+2:]...)
				}else{
					args = args[:i]
				}
			}
		}
	}

	if rests, exists := ExtractOption(args, "-nh"); exists {
		args = rests
		noheader = true
	}

	if rests, exists := ExtractOption(args, "-t"); exists {
		args = rests
		delim = '\t'
	}

	if rests, exists := ExtractOption(args, "-s"); exists {
		args = rests
		delim = ' '
	}

	// New CSV Reader
	var in io.Reader
	if terminal.IsTerminal(0) {
		if len(args) < 2 {
			fmt.Print(usage)
			os.Exit(1)
		}
		f, err := os.Open(args[1])
		if err != nil {
			log.Fatalf("os.Open error: %v", err)
		}
		if info, err := f.Stat(); err != nil {
			log.Fatalf("f.Stat error: %v", err)
		}else if info.IsDir() {
			log.Fatalln("input should be file.")
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}
	bufin := bufio.NewReaderSize(in, 4096)
	// skip row
	for i := 0; i < skip; i++ {
		bufin.ReadLine()
	}
	sc := csv.NewReader(bufin)
	sc.Comma = delim
	sc.TrimLeadingSpace = true

	c, err := catle.NewCatle(sc, noheader)
	defer c.Close()
	if err != nil {
		log.Fatalf("init error: %q", err)
	}
	c.PrintTable()
	c.PollEvent()
}
