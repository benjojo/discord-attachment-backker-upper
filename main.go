package main

import (
	"archive/zip"
	"crypto/md5"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	outputDir := flag.String("output-dir", "./dump/", "Where you want all the attachments to go")
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Fatalf("Please provide a path to a dump zip file")
	}

	zipF, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("Cannot open zip file (OS Level) %v", err)
	}

	fi, err := zipF.Stat()
	if err != nil {
		log.Fatalf("Cannot open zip file (OS stat Level) %v", err)
	}

	zipReader, err := zip.NewReader(zipF, fi.Size())
	if err != nil {
		log.Fatalf("Cannot open zip file (Zip Level) %v", err)
	}

	for _, f := range zipReader.File {
		if !(strings.HasPrefix(f.Name, "messages/") && strings.HasSuffix(f.Name, ".csv")) {
			continue
		}

		csvFile, err := f.Open()
		if err != nil {
			log.Printf("Cannot open %v inside zip file %v", f.Name, err)
			continue
		}

		csvReader := csv.NewReader(csvFile)

		headers, err := csvReader.Read()
		if err != nil {
			log.Printf("Failed to read CSV readers %v", err)
			continue
		}

		if len(headers) < 3 {
			log.Printf("Failed to read CSV readers %v", err)
			continue
		}

		if headers[3] != "Attachments" {
			log.Printf("Unexpected file format in %v", f.Name)
			continue
		}

		log.Printf("Processing %v", f.Name)

		for {
			rows, err := csvReader.Read()
			if err != nil {
				break
			}

			if rows[3] == "" {
				continue
			}

			url, err := url.Parse(rows[3])
			if err != nil {
				log.Printf("Invalid URL %v", rows[3])
				continue
			}

			// 2022-08-03 22:56:28.245000+00:00
			timestamp, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", rows[1])
			if err != nil {
				log.Printf("bad ts %v", err)
				continue
			}

			resp, err := http.Get(url.String())
			if err != nil {
				log.Printf("Failed to fetch %v", err)
			}

			urlHash := md5.New()
			urlHash.Write([]byte(url.String()))
			filename := fmt.Sprintf("%x", urlHash.Sum(nil))
			urlBits := strings.Split(url.String(), ".")
			ext := urlBits[len(urlBits)-1]

			os.MkdirAll(fmt.Sprintf("%s/%s", *outputDir, timestamp.Format("2006-01-02")), 0777)
			cf, _ := os.Create(fmt.Sprintf("%s/%s/%s.%s", *outputDir, timestamp.Format("2006-01-02"), filename, ext))
			io.Copy(cf, resp.Body)
			cf.Close()

		}

	}
}
