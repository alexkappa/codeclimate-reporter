package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alexkappa/errors"
)

type SourceFile struct {
	Name     string        `json:"name"`
	Coverage []interface{} `json:"coverage"`
	BlobID   string        `json:"blob_id"`
}

func (s *SourceFile) String() string {
	c := make(map[int]int64)
	for line, hits := range s.Coverage {
		if hits != nil {
			c[line] = hits.(int64)
		}
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Name: %s\n", s.Name)
	fmt.Fprintf(&buf, "Coverage: %v\n", c)
	fmt.Fprintf(&buf, "BlobID: %v\n\n", s.BlobID)
	return buf.String()
}

func collectSource(c coverage) ([]SourceFile, error) {
	src := make([]SourceFile, 0)

	for file, info := range c {
		cwd, _ := os.Getwd()
		rel, _ := filepath.Rel(cwd, file)

		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		len := bytes.Count(b, []byte{'\n'})

		cov := make([]interface{}, len)

		for _, line := range info.lines {
			cov[line.line-1] = line.hits
		}

		src = append(src, SourceFile{
			Name:     rel,
			Coverage: cov,
			BlobID:   calculateBlobID(b),
		})
	}

	return src, nil
}

type coverage map[string]coverageFile

type coverageFile struct {
	lines []coverageLine
}

type coverageLine struct {
	line int64
	hits int64
}

func collectCoverage(r io.Reader) (coverage, error) {
	b := bufio.NewReader(r)

	f, err := b.ReadString('\n')
	if err != nil {
		return nil, errors.Wrap(err, "Read failed")
	}

	if !strings.HasPrefix(f, "mode:") {
		return nil, errors.New("Unknown coverage format")
	}

	cov := make(coverage)

	for {
		l, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "Read failed")
		}

		regex := regexp.MustCompile(`(.*?):(\d+).\d+,(\d+).\d+ \d+ (\d+)`)

		var match []string
		if match = regex.FindStringSubmatch(l); len(match) != 5 {
			return nil, errors.Wrap(err, "Regexp match failed")
		}

		file := os.Getenv("GOPATH") + "/src/" + match[1]
		start, _ := strconv.ParseInt(match[2], 10, 64)
		end, _ := strconv.ParseInt(match[3], 10, 64)
		hits, _ := strconv.ParseInt(match[4], 10, 64)

		covFile := coverageFile{}

		if _, ok := cov[file]; !ok {
			covFile.lines = make([]coverageLine, 0)
		}

		for line := start; line <= end; line++ {
			covFile.lines = append(covFile.lines, coverageLine{line, hits})
		}

		cov[file] = covFile
	}

	return cov, nil
}

func calculateBlobID(b []byte) string {
	h := sha1.New()
	io.WriteString(h, "blob ")
	io.WriteString(h, strconv.Itoa(len(b)))
	io.WriteString(h, ` \0`)
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}
