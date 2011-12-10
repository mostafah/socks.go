/*
Copyright 2011 Mostafa Hajizdeh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	socksAddr string
)

// socksDial performs as a replacement of the default net.Dial function. It
// is capable of transfering the connection through a SOCKS5 proxy, defined
// by socksAddr.
func socksDial(network, addr string) (net.Conn, error) {
	sAddr, _ := net.ResolveTCPAddr(network, socksAddr)
	domain := addr[:len(addr)-3]
	var port uint16 = 80
	conn, err := net.Dial(network, fmt.Sprintf("%s:%d", sAddr.IP, sAddr.Port))
	conn.Write([]byte{5, 1, 0})
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("<", n, buf[:n])
	send := make([]byte, 5+len(domain)+2)
	send[0] = 5
	send[1] = 1
	send[2] = 0
	send[3] = 3
	send[4] = byte(len(domain))
	for i, c := range []byte(domain) {
		send[5+i] = c
	}
	send[len(send)-2] = byte(port >> 8)
	send[len(send)-1] = byte(port & 0xff)
	fmt.Println(">", len(send), send)
	n, err = conn.Write(send)
	fmt.Println(n, err)
	n, _ = conn.Read(buf)
	fmt.Println("<", n, buf[:n])
	return conn, err
}

// getHTMLTitle parses r as HTML and returns its title. An empty string is
// returned if no <title> tag is found or error happens. Please note that this
// function doesn't check for the <html> -> <head> -> <title> tag hierarchy,
// but just picks the first <title> tag and returns its text.
func getHTMLTitle(r io.Reader) string {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			// probably the end of buffer, not an actual error
			break
		}

		if tt == html.StartTagToken {
			tn, _ := z.TagName()
			if string(tn) == "title" {
				// found the <title> tag, now return the next
				// token, which is actually the text just
				// after <title>
				z.Next()
				return strings.Trim(z.Token().String(), " \t\n\r")
			}
		}
	}

	// return empty in case of error or when no <title> tag was found
	return ""
}

// testUrl downloads the webpage at the url u and prints the HTTP status code
// and its HTML title. It returns HTTP status code.
func testUrl(u string) int {
	res, err := http.Get(u)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0
	}
	defer res.Body.Close()
	fmt.Println(res.Status)
	fmt.Println(getHTMLTitle(res.Body))
	return res.StatusCode
}

func main() {
	flag.Parse()

	socksAddr = "127.0.0.1:1920"
	http.DefaultTransport = &http.Transport{Dial: socksDial}

	if flag.NArg() < 1 {
		testUrl("http://www.google.com/")
	} else {
		u := flag.Arg(0)
		if strings.ToLower(u)[:7] != "http://" {
			u = "http://" + flag.Arg(0)
		}
		testUrl(u)
	}
}
