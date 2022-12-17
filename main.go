package main

import (
	"fmt"
	"io"
	"os"
)

/**
IO with readers and writers
The io.Reader interface
The io.Writer interface

*/
type Reader interface {
	Read(p []byte) (n int, err error)
}

//type alphaReader string // an alphabet reader to read only letters A-Z||a-z
//
//func (a alphaReader) Read(p []byte) (n int, err error) {
//	count := 0
//
//	for i := 0; i < len(a); i++ {
//		if (a[i] >= 'A' && a[i] <= 'Z') || (a[i] >= 'a' && a[i] <= 'z') {
//			p[i] = a[i]
//		}
//		count++
//	}
//
//	return count, io.EOF
//}

type alphaReader struct {
	src io.Reader
}

func NewAlphaReader(source io.Reader) *alphaReader {
	return &alphaReader{source}
}

func (a *alphaReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	count, err := a.src.Read(p)
	if err != nil {
		return count, err
	}

	for i := 0; i < len(p); i++ {
		if (p[i] >= 'A' && p[i] <= 'Z') || (p[i] >= 'a' && p[i] <= 'z') {
			continue
		} else {
			p[i] = 0
		}
	}

	return count, io.EOF
}

func mainReader() {
	file, _ := os.Open("./go.mod")
	alpha := NewAlphaReader(file)
	io.Copy(os.Stdout, alpha)
	fmt.Println()
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type channelWriter struct {
	Channel chan byte
}

func NewChannelWriter() *channelWriter {
	return &channelWriter{
		Channel: make(chan byte, 1024),
	}
}

func (c *channelWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	go func() {
		defer close(c.Channel)
		for _, b := range p {
			c.Channel <- b
		}
		//defer is going to run here
	}()

	fmt.Println(len(p))

	return len(p), nil
}

//func main() {
//	cw := NewChannelWriter()
//
//	go func() {
//		fmt.Fprintf(cw, "We built a channel writer!")
//	}()
//
//	for c := range cw.Channel {
//		fmt.Printf("%c \n", c)
//	}
//}

func main() {
	cw := NewChannelWriter()

	file, err := os.Open("./writer.go")

	if err != nil {
		fmt.Println("Error reading file: ", err)
		os.Exit(1)
	}

	_, err = io.Copy(cw, file)

	if err != nil {
		fmt.Println("Error copying: ", err)
		os.Exit(1)
	}

	for c := range cw.Channel {
		fmt.Printf("%c \n", c)
	}

}
