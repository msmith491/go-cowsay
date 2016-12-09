package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

var (
	whichCow = flag.String("f", "default", "Which cow should say")
	listCow  = flag.Bool("l", false, "List available cowfiles")
	wrapCow  = flag.Bool("n", false, "Disable word wrap")
	widthCow = flag.Int("W", 40, "Width of cow bubble in characters")

	borgCow     = flag.Bool("b", false, "Borg Cow Mode")
	greedyCow   = flag.Bool("g", false, "Greedy Cow Mode")
	paranoidCow = flag.Bool("p", false, "Paranoid Cow Mode")
	stonedCow   = flag.Bool("s", false, "Stoned Cow Mode")
	tiredCow    = flag.Bool("t", false, "Tired Cow Mode")
	wiredCow    = flag.Bool("w", false, "Wired Cow Mode")
	youthfulCow = flag.Bool("y", false, "Youthful Cow Mode")
	deadCow     = flag.Bool("d", false, "Dead Cow Mode")
	customCow   = flag.String("e", "", "Custom Cow Eye String")
	tongueCow   = flag.String("T", "", "Custom Cow Tongue String")

	randomCow = flag.Bool("random", false, "Choose Random Cow")

	smallSideL string
	smallSideR string
	bigSideL   string
	bigSideR   string
	startL     string
	startR     string
	endL       string
	endR       string
	voice      string
)

func getAsset(choice string) []byte {
	if strings.HasSuffix(choice, ".cow") {
		data, err := Asset(choice)
		if err != nil {
			fmt.Println("Couldn't access asset")
		}
		return data
	} else {
		data, err := Asset(fmt.Sprintf("src/go-cowsay/cows/%s.cow", choice))
		if err != nil {
			fmt.Println("Couldn't access asset")
		}
		return data
	}
}

func getEyes() string {
	var result string
	if *borgCow {
		result = "=="
	} else if *greedyCow {
		result = "$$"
	} else if *paranoidCow {
		result = "@@"
	} else if *stonedCow {
		result = "**"
	} else if *tiredCow {
		result = "--"
	} else if *wiredCow {
		result = "OO"
	} else if *youthfulCow {
		result = ".."
	} else if *deadCow {
		result = "xx"
	} else if *customCow != "" {
		result = (*customCow)[:2]
	} else {
		result = "oo"
	}
	return result
}

func makeBubble(s string) string {
	// Buffer to write and clear for each line
	var b bytes.Buffer
	// String array to join and return
	var result []string
	longest := 0

	pad := " "
	if len(s) < *widthCow || *wrapCow {
		b.WriteString(smallSideL + pad + s + pad + smallSideR)
		result = append(result, b.String())
		b.Reset()
	} else {
		index := 0
		text := ""
		for true {
			if len(s) <= index+*widthCow {
				text = s[index:]
			} else {
				text = s[index : index+*widthCow]
			}
			if len(text) < *widthCow {
				// Last Text Lines
				b.WriteString(endL + pad + strings.TrimSpace(text) + pad + endR)
				result = append(result, b.String())
				b.Reset()
				break
			}
			// Middle Lines
			if index == 0 {
				b.WriteString(startL + pad)
			} else {
				b.WriteString(bigSideL + pad)
			}
			if text[len(text)-1] != ' ' || s[len(text)+index] != ' ' {
				split := strings.Split(text, " ")
				text = strings.Join(split[:len(split)-1], " ")
				b.WriteString(text)
				lenLast := len(split[len(split)-1])
				index -= lenLast
			} else {
				b.WriteString(text)
			}
			if index <= 0 {
				b.WriteString(pad + startR)
			} else {
				b.WriteString(pad + bigSideR)
			}
			result = append(result, b.String())
			b.Reset()
			index += *widthCow
		}
	}

	for _, line := range result {
		if len(line) > longest {
			longest = len(line)
		}
	}

	for index, line := range result {
		tmp := line[len(line)-1]
		line = line[:len(line)-1]
		diff := longest - len(line)
		for i := 0; i <= diff; i++ {
			line += pad
		}
		line += string(tmp)
		result[index] = line
	}

	b.WriteString(pad)
	for i := 0; i < longest; i++ {
		b.WriteString("_")
	}
	result = append([]string{b.String()}, result...)
	b.Reset()

	b.WriteString(pad)
	for i := 0; i < longest; i++ {
		b.WriteString("-")
	}
	result = append(result, b.String())
	b.Reset()

	return strings.Join(result, "\n")
}

func formatAnimal(s string) string {
	var animal string
	var tongue string
	if *tongueCow != "" {
		tongue = (*tongueCow)[:2]
	} else {
		tongue = "  "
	}
	animal = strings.Replace(s, "$eyes", getEyes(), -1)
	animal = strings.Replace(animal, "${eyes}", getEyes(), -1)
	animal = strings.Replace(animal, "$tongue", tongue, -1)
	animal = strings.Replace(animal, "$thoughts", voice, -1)
	animal = strings.Replace(animal, "\\\\", "\\", -1)
	animal = strings.Replace(animal, "\\@", "@", -1)
	return animal
}

func main() {
	flag.Parse()

	think := false
	if strings.Contains(os.Args[0], "cowthink") {
		think = true
	}

	var message string
	fi, _ := os.Stdin.Stat()

	if *listCow {
		for key := range _bindata {
			f := strings.Split(strings.Replace(key, ".cow", "", -1), "/")
			fmt.Printf(f[len(f)-1] + " ")
		}
		fmt.Println()
		os.Exit(0)
	}

	if think {
		smallSideL = "("
		smallSideR = ")"
		bigSideL = "("
		bigSideR = ")"
		startL = "("
		startR = ")"
		endL = "("
		endR = ")"
		voice = "o"
	} else {
		smallSideL = "<"
		smallSideR = ">"
		bigSideL = "|"
		bigSideR = "|"
		startL = "/"
		startR = "\\"
		endL = "\\"
		endR = "/"
		voice = "\\"
	}

	if fi.Size() > 0 {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		message = makeBubble(strings.TrimSpace(text))
	} else {
		args := strings.Join(flag.Args(), " ")
		if args == "" {
			flag.Usage()
			os.Exit(1)
		}
		message = makeBubble(args)
	}

	var data []byte
	var err error
	if *randomCow {
		var assets []string
		for key, _ := range _bindata {
			assets = append(assets, key)
		}
		choice := assets[rand.Intn(len(assets))]
		data = getAsset(choice)
	} else if strings.HasSuffix(*whichCow, ".cow") {
		data, err = ioutil.ReadFile(*whichCow)
		if err != nil {
			fmt.Println("Couldn't access file")
		}
	} else {
		data = getAsset(*whichCow)
	}
	sdata := fmt.Sprintf("%s", data)
	cow := formatAnimal(sdata)
	for _, line := range strings.Split(message, "\n") {
		fmt.Println(line)
	}
	for _, line := range strings.Split(cow, "\n")[1:] {
		if !strings.Contains(line, "EOC") && !strings.HasPrefix(line, "##") {
			fmt.Println(line)
		}
	}
}
