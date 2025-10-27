package input

import (
	"bufio"
	"log"
	"os"
)

type FileInput struct {
	FilePath string

	file   *os.File
	reader *bufio.Reader
}

func (f *FileInput) Init() error {
	var err error
	f.file, err = os.Open(f.FilePath)
	if err != nil {
		log.Printf("Fail to open file: %v\n", f.FilePath)
		return err
	}
	f.reader = bufio.NewReader(f.file)
	return nil
}

func (f *FileInput) FetchEvent() (string, error) {
	line, err := f.reader.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			//log.Println("Read end of file")
			return "", err
		}
		log.Println("Fail to read line")
		return "", err
	}
	return line, nil
}

func (f *FileInput) Close() {
	if f.file == nil {
		return
	}
	f.file.Close()
	log.Println("File closed")
}
