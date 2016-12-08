package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	whichCow = flag.String("f", "default", "Which cow should say")
	listCow  = flag.Bool("l", false, "List available cowfiles")
	wrapCow  = flag.Bool("n", false, "Disable word wrap")
	widthCow = flag.Int("W", 40, "Width of cow bubble in characters")
)

func main() {
	flag.Parse()
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
	if fi.Size() > 0 {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		message = makeBubble(strings.TrimSpace(text), *wrapCow)
	} else {
		message = makeBubble(strings.Join(flag.Args(), " "), *wrapCow)
	}

	data, err := Asset(fmt.Sprintf("src/go-cowsay/cows/%s.cow", *whichCow))
	if err != nil {
		fmt.Println("Couldn't access asset")
	}
	sdata := fmt.Sprintf("%s", data)
	cow := formatAnimal(sdata)
	for _, line := range strings.Split(message, "\n") {
		fmt.Println(line)
	}
	for _, line := range strings.Split(cow, "\n")[1:] {
		if !strings.Contains(line, "EOC") && !strings.HasPrefix(line, "##") {
			// Hack to fix the cow belly/legs
			if strings.Contains(line, "||----w ") {
				fmt.Println("  " + line)
			} else {
				fmt.Println(line)
			}
		}
	}
}

func makeBubble(s string, wordWrap bool) string {
	var b bytes.Buffer
	var b2 bytes.Buffer
	var writelen int

	pad := " "
	if len(s) < *widthCow || wordWrap {
		writelen = len(s)
		b.WriteString("<" + pad)
		b.WriteString(s)
		b.WriteString(pad + ">\n")
	} else {
		writelen = *widthCow + 1
		index := 0
		// Top text line
		text := s[:index+*widthCow]
		b.WriteString("/" + pad)
		b.WriteString(text)
		b.WriteString(pad + "\\\n")
		index += *widthCow
		for true {
			if len(s) <= index+*widthCow {
				text = s[index:]
			} else {
				text = s[index : index+*widthCow]
			}
			if len(text) < *widthCow {
				// Last Text Lines
				b.WriteString("\\" + pad)
				b.WriteString(strings.TrimSpace(text))
				for i := 0; i < *widthCow-len(text); i++ {
					b.WriteString(pad)
				}
				b.WriteString(pad + "/\n")
				break
			}
			// Middle Lines
			b.WriteString("|" + pad)
			if text[len(text)-1] != ' ' || s[len(text)+index] != ' ' {
				split := strings.Split(text, " ")
				b.WriteString(strings.Join(split[:len(split)-1], " "))
				lenLast := len(split[len(split)-1])
				index -= lenLast
				for i := 0; i <= lenLast; i++ {
					b.WriteString(pad)
				}
			} else {
				b.WriteString(text)
			}
			b.WriteString(pad + "|\n")
			index += *widthCow
		}
	}
	b2.WriteString(pad)
	for i := 0; i <= writelen; i++ {
		b2.WriteString("_")
	}
	b2.WriteString("\n")
	b2.Write(b.Bytes())
	b2.WriteString(pad)
	for i := 0; i <= writelen; i++ {
		b2.WriteString("-")
	}
	return b2.String()
}

func formatAnimal(s string) string {
	var animal string
	animal = strings.Replace(s, "$eyes", "oo", -1)
	animal = strings.Replace(animal, "$tongue", "", -1)
	animal = strings.Replace(animal, "$thoughts", "\\", -1)
	animal = strings.Replace(animal, "\\\\", "\\", -1)
	animal = strings.Replace(animal, "\\@", "@", -1)
	return animal
}
