package notrelevant

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type WordCounter int
type LineCounter int
type countingWriter struct {
	w     io.Writer
	count *int64
}

func (cw countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)

	if err != nil {
		*cw.count += int64(n)
	}

	return n, nil
}

func (c *WordCounter) Write(p []byte) (int, error) {
	count := counter(p, bufio.ScanWords)
	*c += WordCounter(count)

	return count, nil
}

func (l *LineCounter) Write(p []byte) (int, error) {
	count := counter(p, bufio.ScanLines)
	*l += LineCounter(count)

	return count, nil
}

func counter(p []byte, fn bufio.SplitFunc) (count int) {
	sc := bufio.NewScanner(strings.NewReader(string(p)))

	sc.Split(fn)

	for sc.Scan() {
		count++
	}

	if err := sc.Err(); err != nil {
		fmt.Errorf("%s", err)
	}

	return
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var cn int64

	val := countingWriter{w, &cn}

	return val, val.count
}
