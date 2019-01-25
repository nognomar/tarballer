package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	source := flag.String("source", "", "Source directory")
	target := flag.String("target", "", "Target tarball file")
	flag.Parse()

	log.Println("Tarballer started.")
	createTarball(*source, *target)
	log.Println("Done.")
}

func createTarball(source string, target string) {
	source = strings.Replace(source, "\\", "/", -1)
	if !strings.HasSuffix(source, "/") {
		source += "/"
	}

	if !strings.HasSuffix(target, ".tar.gz") {
		target += ".tar.gz"
	}

	fw, err := os.Create(target)
	checkErr(err)
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	compressDir(source, tw, source)
}

func compressDir(dirPath string, tw *tar.Writer, stripPath string) {
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	dir, err := os.Open(dirPath)
	checkErr(err)
	defer dir.Close()

	files, err := dir.Readdir(0)
	checkErr(err)
	for _, file := range files {
		filePath := dirPath + file.Name()
		if file.IsDir() {
			compressDir(filePath, tw, stripPath)
		} else {
			log.Println("add " + filePath)
			compressFile(filePath, tw, file, stripPath)
		}
	}
}

func compressFile(filePath string, tw *tar.Writer, file os.FileInfo, stripPath string) {
	fr, err := os.Open(filePath)
	checkErr(err)
	defer fr.Close()

	header := new(tar.Header)
	header.Name = strings.Replace(filePath, stripPath, "", -1)
	header.Size = file.Size()
	header.Mode = int64(file.Mode())
	header.ModTime = file.ModTime()

	err = tw.WriteHeader(header)
	checkErr(err)

	_, err = io.Copy(tw, fr)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-666)
	}
}
