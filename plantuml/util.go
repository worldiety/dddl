package plantuml

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var idNum int
var mutex sync.Mutex

type WithId interface {
	Id() string
}

func nextId() string {
	mutex.Lock()
	defer mutex.Unlock()
	idNum++

	return strconv.Itoa(idNum)
}

type strWriter struct {
	Writer io.Writer
	Err    error
}

func (w strWriter) Print(str string) {
	if w.Err != nil {
		return
	}

	_, err := w.Writer.Write([]byte(str))
	if err != nil {
		w.Err = err
	}
}

func (w strWriter) Printf(format string, args ...interface{}) {
	if w.Err != nil {
		return
	}

	_, err := w.Writer.Write([]byte(fmt.Sprintf(format, args...)))
	if err != nil {
		w.Err = err
	}
}

func (w strWriter) Write(p []byte) (n int, err error) {
	if w.Err != nil {
		return 0, err
	}

	n, err = w.Writer.Write(p)
	if err != nil {
		w.Err = err
	}

	return
}

func escapeP(str string) string {
	return strings.ReplaceAll(str, `"`, "<U+0022>")
}

type Renderable interface {
	Render(wr io.Writer) error
}

type Raw string

func (r Raw) Render(wr io.Writer) error {
	_, err := wr.Write([]byte(r))
	return err
}

func String(r Renderable) string {
	sb := &bytes.Buffer{}
	if err := r.Render(sb); err != nil {
		return err.Error()
	}

	return sb.String()
}

var globalCache = map[string][]byte{}
var cacheMutex sync.Mutex

// Render invokes a local plantuml command. fileType can be svg, or png etc.
func RenderLocal(fileType string, renderable Renderable) ([]byte, error) {
	now := time.Now()
	defer func() {
		slog.Info("plantuml render took", slog.Any("dur", time.Now().Sub(now)))
	}()
	tmp := String(renderable)
	cacheMutex.Lock()
	buf, ok := globalCache[tmp]
	if ok {
		cacheMutex.Unlock()
		return buf, nil
	}
	cacheMutex.Unlock()

	if buf := readFileCache(fileType, tmp); buf != nil {
		return buf, nil
	}

	cmd := exec.Command("plantuml", "-t"+fileType, "-p")
	cmd.Env = os.Environ()

	w, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if _, err := w.Write([]byte(tmp)); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	buf, err = cmd.Output()
	if err != nil {
		log.Println(tmp)
		return buf, err
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	globalCache[tmp] = buf

	writeFileCache(fileType, tmp, buf)

	return buf, nil
}

func fileCacheDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = os.TempDir()
	}

	myTmpDir := filepath.Join(dir, ".plantuml-cache")
	if err := os.MkdirAll(myTmpDir, os.ModePerm); err != nil {
		log.Println("cannot access plantuml-cache tmp dir", err, myTmpDir)
	}

	return myTmpDir
}

func cacheFilename(ftype, str string) string {
	myTmpDir := fileCacheDir()
	sum := sha256.Sum224([]byte(ftype + str))
	return filepath.Join(myTmpDir, hex.EncodeToString(sum[:]))

}

func writeFileCache(ftype, str string, buf []byte) {
	fname := cacheFilename(ftype, str)
	if err := os.WriteFile(fname, buf, os.ModePerm); err != nil {
		log.Println("cannot write into plantuml-cache tmp dir", err, fname)
	}
}

func readFileCache(ftype, str string) []byte {
	fname := cacheFilename(ftype, str)
	buf, err := os.ReadFile(fname)
	if err == nil {
		return buf
	}

	return nil
}
