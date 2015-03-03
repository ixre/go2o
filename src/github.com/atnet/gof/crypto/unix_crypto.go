/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-25 21:23
 * description :
 * history :
 */

package crypto

//
//cyp := NewMd5Crypto("ops", "rdm")
//i :=2
//for {
//if i = i - 1; i ==0  {
//break
//}
//
//str := cyp.Encode()
//fmt.Println("str:", str)
//
//_,unix := cyp.Decode(str)
//fmt.Println(time.Now().Unix()-unix)
//}

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

var (
	d5 = md5.New()
)

const (
	unixLen = 10 //unix time长度为10
)

func getPos(token string) int {
	return len(token)/2 + 1
}
func getUnix() string {
	ux := time.Now().Unix()
	return strconv.FormatInt(ux, 10)
}

func getMd5(token, offset string) []byte {
	d5.Reset()
	d5.Write([]byte(token))
	src := d5.Sum([]byte(offset))
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

type UnixCrypto struct {
	pos      int
	md5Bytes []byte
	buf      *bytes.Buffer
	mux      sync.Mutex
}

func NewUnixCrypto(token, offset string) *UnixCrypto {
	return &UnixCrypto{
		pos:      len(token)/2 + 1,
		md5Bytes: getMd5(token, offset),
		buf:      bytes.NewBufferString(""),
	}
}

// return md5 bytes
func (this *UnixCrypto) GetBytes() []byte {
	return this.md5Bytes
}

func (this *UnixCrypto) Encode() []byte {
	unx := getUnix()
	l := this.pos

	this.mux.Lock()
	defer func() {
		this.buf.Reset()
		this.mux.Unlock()
	}()

	this.buf.Write(this.md5Bytes[:l])

	for i := 0; i < 10; i++ {
		this.buf.WriteString(unx[i : i+1])
		this.buf.Write(this.md5Bytes[l+i : l+i+1])
	}

	this.buf.Write(this.md5Bytes[10+l:])
	return this.buf.Bytes()
}

func (this *UnixCrypto) Decode(s string) ([]byte, int64) {
	smd := make([]byte, len(this.md5Bytes))
	unx := make([]byte, unixLen)

	if len(s) < len(smd) {
		return nil, 0
	}

	copy(smd, s[:this.pos])
	for i, v := range s[this.pos+unixLen*2:] {
		smd[this.pos+unixLen+i] = byte(v)
	}

	for i := 0; i < unixLen*2; i++ {
		v := s[this.pos+i]
		if i%2 == 0 {
			unx[i/2] = v
		} else {
			smd[this.pos+i/2] = v
		}
	}

	unix, _ := strconv.ParseInt(string(unx), 10, 32)
	return smd, unix
}

func (this *UnixCrypto) Compare(s string) (bool, []byte, int64) {
	b, u := this.Decode(s)
	return bytes.Compare(b, this.md5Bytes) == 0, b, u
}
