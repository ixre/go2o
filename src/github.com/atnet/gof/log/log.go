/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-10 23:52
 * description :
 * history :
 */

package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	LOpen = 1 << iota

	LStdFlags

	// show file source
	LSource

	// only show error source
	LESource

	DEFAULT_DEPTH = 3 // 在此包中存在3层
)

var (
	std ILogger = NewLogger(os.Stderr, "", LOpen|LStdFlags|LESource|DEFAULT_DEPTH)
)

type ILogger interface {
	AddDepth(int)
	ResetDepth()
	Println(...interface{})
	Printf(string, ...interface{})
	PrintErr(error)
	SetFlag(int)
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
}

type simpleLogger struct {
	out         io.Writer
	buf         []byte
	mux         sync.Mutex
	flag        int
	opened      bool
	prefix      string
	callerDepth int
}

func NewLogger(writer io.Writer, prefix string, flag int) *simpleLogger {
	if writer == nil {
		writer = os.Stdout
	}
	l := &simpleLogger{
		out:         writer,
		flag:        flag,
		opened:      flag&LOpen != 0,
		callerDepth: DEFAULT_DEPTH,
	}

	if len(prefix) != 0 {
		l.prefix = "[" + prefix + "]"
	}

	return l
}

func (t *simpleLogger) SetFlag(flag int) {
	t.flag = flag
}

func (t *simpleLogger) formatHeader(b *[]byte) {
	if t.flag&LStdFlags != 0 {
		now := time.Now()
		*b = append(*b, now.Format("2006-01-02 15:04:05 ** ")...)
		if len(t.prefix) != 0 {
			*b = append(*b, t.prefix...)
		}
	}
}

func (t *simpleLogger) appendSource(b *[]byte) {
	_, f, l, ok := runtime.Caller(t.callerDepth)
	if ok {
		*b = append(*b, (" (Source:" + f)...)
		*b = append(*b, " - Line:"...)
		*b = append(*b, strconv.Itoa(l)...)
		*b = append(*b, ')')
	} else {
		*b = append(*b, "Source:???"...)
	}
}

func (t *simpleLogger) output(err bool, b *[]byte, newLine bool) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.flag&LESource != 0 {
		// if not error,but only show error source
		if err {
			t.appendSource(b)
		}
	} else {
		if t.flag&LSource != 0 {
			t.appendSource(b)
		}
	}

	if newLine {
		*b = append(*b, '\n')
	}

	t.out.Write(*b)
}

func (t *simpleLogger) AddDepth(depth int) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.callerDepth += depth
}
func (t *simpleLogger) ResetDepth() {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.callerDepth = DEFAULT_DEPTH
}

func (t *simpleLogger) Println(v ...interface{}) {
	if t.opened {
		t.buf = t.buf[:0]
		t.formatHeader(&t.buf)
		t.buf = append(t.buf, fmt.Sprint(v...)...)
		t.output(false, &t.buf, true)
	}
}

func (t *simpleLogger) Printf(s string, v ...interface{}) {
	if t.opened {
		t.buf = t.buf[:0]
		t.formatHeader(&t.buf)
		t.buf = append(t.buf, fmt.Sprintf(s, v...)...)
		t.output(false, &t.buf, false)
	}
}

func (t *simpleLogger) PrintErr(e error) {
	if t.opened && e != nil {
		t.buf = t.buf[:0]
		t.formatHeader(&t.buf)
		t.buf = append(t.buf, fmt.Sprintf("[ Error] - %s", e.Error())...)
		t.output(false, &t.buf, true)
	}
}

func (t *simpleLogger) Panicf(s string, v ...interface{}) {
	t.buf = t.buf[:0]
	str := fmt.Sprintf(s, v...)
	t.formatHeader(&t.buf)
	t.buf = append(t.buf, str...)
	t.output(false, &t.buf, false)
	panic(str)
}

func (t *simpleLogger) Panicln(v ...interface{}) {
	t.buf = t.buf[:0]
	str := fmt.Sprint(v...)
	t.formatHeader(&t.buf)
	t.buf = append(t.buf, str...)
	t.output(false, &t.buf, true)
	panic(str)
}

func (t *simpleLogger) Fatalf(s string, v ...interface{}) {
	t.buf = t.buf[:0]
	t.formatHeader(&t.buf)
	t.buf = append(t.buf, fmt.Sprintf(s, v...)...)
	t.output(false, &t.buf, false)
	os.Exit(1)
}

func (t *simpleLogger) Fatalln(v ...interface{}) {
	t.buf = t.buf[:0]
	t.formatHeader(&t.buf)
	t.buf = append(t.buf, fmt.Sprint(v...)...)
	t.output(false, &t.buf, true)
	os.Exit(1)
}

func AddDepth(i int) {
	std.AddDepth(i)
}

func ResetDepth() {
	std.ResetDepth()
}

func SetFlag(flag int) {
	std.SetFlag(flag)
}

func Println(v ...interface{}) {
	std.Println(v...)
}

func Printf(s string, v ...interface{}) {
	std.Printf(s, v...)
}

func PrintErr(e error) {
	std.PrintErr(e)
}

func Panicf(s string, v ...interface{}) {
	std.Panicf(s, v...)
}

func Panicln(v ...interface{}) {
	std.Panicln(v...)
}

func Fatalf(s string, v ...interface{}) {
	std.Fatalf(s, v...)
}

func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

//func test(){
//	var logger ILogger = NewLogger(nil,"WEB",LOpen | LStdFlags | LSource)
//	logger.Println("[ BOOTStrap] - is boot now!")
//	logger.Println("[ NOTICE] - fsdfsdfsdf")
//}
