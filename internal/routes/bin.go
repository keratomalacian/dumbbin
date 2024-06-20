package routes

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// handler for getting bins
func GetBin(binsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get binID from url parameter
		var binID = chi.URLParam(r, "binID")

		// get bin path and then read its content
		var binPath = filepath.Join(binsPath, binID)
		binContent, err := os.ReadFile(binPath)

		if err != nil {
			// check if bin doesnt exist
			if errors.Is(err, os.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "bin does not exist")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "internal server error")
			}
			return
		}

		// send the bin content to the user
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(binContent))
	}
}

// handler for creating bins
func CreateBin(binsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// reads the request body into 'body'
		var body, err = io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "error reading request body")
			return
		}

		// converts the bin content from bytes to a string and checks if it is empty
		var binContent = string(body)
		if len(strings.TrimSpace(binContent)) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "missing request body")
			return
		}

		// hashes the bin content using SHA1
		var sha = sha1.New()
		sha.Write(body)

		// converts the hash to hexadecimal
		var hash = hex.EncodeToString(sha.Sum(nil))

		// gets the path where the bin will be saved
		var binPath = filepath.Join(binsPath, hash)

		// writes the bin content as bytes to 'binPath'
		err = os.WriteFile(binPath, body, 0777)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "internal server error")
			log.Println("could not write file to", binPath, ":", err)
			return
		}

		// checks if the request was done using http or https
		var protocol string
		if r.TLS == nil {
			protocol = "http://"
		} else {
			protocol = "https://"
		}

		// returns the path of the newly created bin to the user
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "created bin at", fmt.Sprint(protocol, r.Host, "/", hash))
	}
}
