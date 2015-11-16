package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileLogWriter struct {
	sync.Mutex
	fd       *os.File
	opendate int
	FileName string
	MaxDays  int
}

func (fw *FileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	var err error
	fw.fd, err = os.OpenFile(fw.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err == nil {
		fw.opendate = time.Now().Day()
	}
	return fw.fd, err
}

func (fw *FileLogWriter) Write(b []byte) (int, error) {
	fw.Lock()
	defer fw.Unlock()
	return fw.fd.Write(b)
}

func (fw *FileLogWriter) check() {
	fw.Lock()
	defer fw.Unlock()
	if time.Now().Day() != fw.opendate {
		if err := fw.rotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", fw.FileName, err)
			return
		}
	}
}

// write log in new file.
// new file name like xx.log.2013-01-01.2
func (fw *FileLogWriter) rotate() error {
	_, err := os.Lstat(fw.FileName)
	if err == nil { // file exists
		// Find the next available number
		fname := fw.FileName + fmt.Sprintf(".%s.%03d", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
		_, err = os.Lstat(fname)
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", fw.FileName)
		}

		// close file before rename
		fw.fd.Close()

		// Rename the file to its newfound home
		err = os.Rename(fw.FileName, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s\n", err)
		}
		// re-start logger
		_, err := fw.createLogFile()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s\n", err)
		}

		go fw.deleteOldLog()
	}

	return nil
}
func (fw *FileLogWriter) deleteOldLog() (err error) {
	if fw.MaxDays <= 0 {
		return
	}
	dir := filepath.Dir(fw.FileName)
	filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
				fmt.Println(returnErr)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*int64(fw.MaxDays)) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(fw.FileName)) {
				os.Remove(path)
			}
		}
		return
	})
	return
}
