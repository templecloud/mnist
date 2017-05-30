//--------------------------------------------------------------------------------------------------
//==================================================================================================
// MNIST IDX ByteImageStore downloader
//--------------------------------------------------------------------------------------------------
package idxfile

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// http://yann.lecun.com/exdb/mnist/
const (
	MNIST_LOCAL        = "./data/"
	MNIST_REMOTE       = "http://yann.lecun.com/exdb/mnist/"
	MNIST_TRAIN_IMAGES = "train-images-idx3-ubyte"
	MNIST_TRAIN_LABELS = "train-labels-idx1-ubyte"
	MNIST_TEST_IMAGES  = "t10k-images-idx3-ubyte"
	MNIST_TEST_LABELS  = "t10k-labels-idx1-ubyte"
	GZIP               = ".gz"
	IDX                = ".idx"
)

type Download struct {
	Source, Destination string
}

func DownloadMNIST() (err error) {
	// ensure the specified local directory exist
	err = os.MkdirAll(MNIST_LOCAL, os.FileMode(0777))
	if err != nil {
		log.Fatal(err)
	}
	// download and unarchive each mnist file
	var mnistFiles = [...]string{
		MNIST_TRAIN_IMAGES, MNIST_TRAIN_LABELS, MNIST_TEST_IMAGES, MNIST_TEST_LABELS}
	for _, v := range mnistFiles {
		// download zipped file
		source := MNIST_REMOTE + v + GZIP
		destination := MNIST_LOCAL + v + GZIP
		unarchived := strings.TrimSuffix(destination, GZIP) + IDX
		if _, err := os.Stat(destination); os.IsNotExist(err) {
			toDownload := Download{source, destination}
			log.Printf("Downloading File: %+v\n", toDownload)
			err := download(source, destination)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Already downloaded: %s\n", destination)
		}
		// unarchive
		if _, err := os.Stat(unarchived); os.IsNotExist(err) {
			log.Printf("Local File     : %s\n", destination)
			log.Printf("Unarchived File: %s\n", unarchived)
			err := gunzip(destination, unarchived)
			if err != nil {
				return err
			}
		} else {
			log.Printf("Already unarchived: %s\n", unarchived)
		}
	}
	return nil
}

func download(url string, filepath string) (err error) {
	// create file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// download mnist
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// save mnsit
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func gunzip(source string, destination string) (err error) {
	// open an input stream on the source path
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	// open a gzip reader on the source stream
	archive, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer archive.Close()
	// open an ouput stream on the destination path
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	// copy the unarchived stream to the destination file
	_, err = io.Copy(out, archive)
	if err != nil {
		return err
	}
	return nil
}
