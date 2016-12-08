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

// boolCow  = flag.Bool("boolcow", false, "Bool cow opt")
)

func main() {
	flag.Parse()
	var message string
	fi, _ := os.Stdin.Stat()
	if fi.Size() > 0 {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		message = makeBubble(strings.TrimSpace(text))
	} else {
		message = makeBubble(strings.Join(flag.Args(), " "))
	}
	// fmt.Println(flag.Args())
	// splitArgs := strings.Split(strings.Join(flag.Args(), " "), "-")
	// firstArgs := splitArgs[:len(splitArgs)]

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
			if *whichCow == "default" && strings.Contains(line, "----w ") {
				fmt.Println("  " + line)
			} else {
				fmt.Println(line)
			}
		}
	}
}

func makeBubble(s string) string {
	var b bytes.Buffer
	var b2 bytes.Buffer
	var writelen int

	pad := " "
	if len(s) < 40 {
		writelen = len(s)
		b.WriteString("<" + pad)
		b.WriteString(s)
		b.WriteString(pad + ">\n")
	} else {
		writelen = 41
		index := 0
		b.WriteString("/" + pad)
		b.WriteString(s[:index+40])
		b.WriteString(pad + "\\\n")
		for true {
			index += 40
			if len(s) < index+40 {
				break
			}
			b.WriteString("|" + pad)
			if s[index+40] == ' ' {
				b.WriteString(s[index : index+40])
			} else {
				subindex := index + 40
				for true {
					if s[subindex] == ' ' {
						break
					}
					subindex -= 1
				}
				b.WriteString(s[index:subindex])
				// Padding for wordwrap
				for i := 0; i < index+40-subindex; i++ {
					b.WriteString(pad)
				}
				index = subindex + 1
			}
			b.WriteString(pad + "|\n")
		}
		b.WriteString("\\" + pad)
		b.WriteString(strings.TrimSpace(s[index:]))
		for i := 0; i < 40-(len(s)-index); i++ {
			b.WriteString(pad)
		}
		b.WriteString(pad + "/\n")
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
	// b2.WriteString("\n")
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
