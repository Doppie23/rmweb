///////////////////////////////////////////////////////////////////////////////
//
// rmweb/download_pdf.go
// John Simpson <jms1@jms1.net> 2023-12-17
//
// Download a PDF file from a reMarkable tablet

package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

///////////////////////////////////////////////////////////////////////////////
//
// Download a PDF file

func download_pdf(uuid string, localfile string) {

	////////////////////////////////////////////////////////////
	// If the output filename contains any directory names,
	// make sure any necessary directories exist.

	for n := 1; n < len(localfile); n++ {
		if localfile[n] == '/' {
			dir := localfile[:n]

			if flag_debug {
				fmt.Printf("checking dir='%s'\n", dir)
			}

			////////////////////////////////////////
			// Check the directory

			s, err := os.Stat(dir)
			if os.IsNotExist(err) {
				////////////////////////////////////////
				// doesn't exist yet - create it

				fmt.Printf("Creating    '%s' ...", dir)

				err := os.Mkdir(dir, 0755)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					os.Exit(1)
				}

				fmt.Println(" ok")

			} else if err != nil {
				////////////////////////////////////////
				// os.Stat() had some other error

				fmt.Printf("ERROR: os.Stat('%s'): %v\n", dir, err)
				os.Exit(1)

			} else if (s.Mode() & fs.ModeDir) == 0 {
				////////////////////////////////////////
				// exists and is not a directory

				fmt.Printf("ERROR: '%s' exists and is not a directory\n", dir)
				os.Exit(1)
			}

		} // if localfile[n] == '/'
	} // for n

	////////////////////////////////////////////////////////////
	// Download the file

	fmt.Printf("Downloading '%s'\n", localfile)

	////////////////////////////////////////
	// Request the file

	url := "http://" + tablet_addr + "/download/" + uuid + "/placeholder"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	////////////////////////////////////////
	// Create output file

	dest, err := os.Create(localfile)
	if err != nil {
		fmt.Printf("ERROR: os.Create('%s'): %v", localfile, err)
		os.Exit(1)
	}

	defer dest.Close()

	////////////////////////////////////////
	// Copy the output to the file

	var src io.Reader = &PassThru{Reader: resp.Body}

	_, err = io.Copy(dest, src)
	if err != nil {
		fmt.Printf("ERROR: os.Copy(): %v", err)
		os.Exit(1)
	}

	////////////////////////////////////////
	// done

	// fmt.Printf( "%d ... ok\n" , total )
}
