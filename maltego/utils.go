/*
 * NETCAP - Traffic Analysis Framework
 * Copyright (c) 2017-2020 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package maltego

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dreadl0ck/netcap/defaults"
	netio "github.com/dreadl0ck/netcap/io"
)

const (
	errUnexpectedFileType    = "unexpected file type"
	errUnexpectedReadFailure = "unexpected read failure"

	// ExecutablePath points to the netcap binary on disk
	ExecutablePath = "/usr/local/bin/net"

	// PropertyIpAddr is the name of maltego property that contains the IP address
	PropertyIpAddr      = "ipaddr"

	// PropertyIpAddrLabel is the label for the ip address property
	PropertyIpAddrLabel = "IPAddress"
)

// EscapeText ensures that the input text is safe to embed within XML.
func EscapeText(text string) string {
	var buf bytes.Buffer

	err := xml.EscapeText(&buf, []byte(text))
	if err != nil {
		fmt.Println(err)
	}

	return buf.String()
}

func die(err string, msg string) {
	trx := Transform{}
	// add error message for the user
	trx.AddUIMessage(msg+": "+err, UIMessageFatal)
	fmt.Println(trx.ReturnOutput())
	log.Println(msg, err)
	os.Exit(0) // don't signal an error for the transform invocation
}

func openFile(path string) (*os.File, string) {

	log.Println("open path:", path)
	f, err := os.Open(path)
	if err != nil {

		log.Println("failed to open path", err, "trying without .gz extension...")

		f, err = os.Open(strings.TrimSuffix(path, ".gz"))
		if err != nil {
			log.Println("failed to open path", err)
			trx := Transform{}
			trx.AddUIMessage("failed to open path: "+err.Error(), UIMessageInform)
			fmt.Println(trx.ReturnOutput())
			os.Exit(0) // don't signal an error for the transform invocation
		}
	}

	return f, path
}

func openNetcapArchive(path string) *netio.Reader {
	r, err := netio.Open(path, defaults.BufferSize)
	if err != nil {

		log.Println("failed to open path ", path, " trying without .gz extension...")
		path = strings.TrimSuffix(path, ".gz")

		r, err = netio.Open(path, defaults.BufferSize)
		if err != nil {
			die(err.Error(), "failed to open file")
		}
	}

	return r
}
