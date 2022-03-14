package util

import (
	"bufio"
	"os"

	"github.com/suecodelabs/cnfuzz/src/log"
)

func WriteToFile(bytesToWrite *[]byte, filePath string) {
	// os.Mkdir("./reports", os.FileMode(0666))
	file, err := os.Create(filePath)
	if err != nil {
		log.L().Errorf("error while creating \"%s\" on filesystem: %+v", filePath, err)
		return
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(*bytesToWrite)
	if err != nil {
		log.L().Errorf("error while writing a report to filesystem: %+v", err)
		return
	}
	writer.Flush()
}
