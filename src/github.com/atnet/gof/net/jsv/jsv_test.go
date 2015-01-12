/**
 * Copyright 2014 @ Yunduo Inc.
 * name :
 * author : newmin
 * date : 2014-12-26 16:29
 * description :
 * history :
 */

package jsv

import (
	"testing"
)

func Test_A(t *testing.T) {
	var buffer []byte = make([]byte, 20480)

	const server = "localhost:1002"
	conn, err := Dial("tcp", server)
	if err != nil {
		t.Fatalf("[TCP]: Connect Refused %s", server)
	}

	data := []byte(`{"partner_id":"666888","secret":"d435a520e50e960b"}>>Partner.GetPartner`)

	for i := 0; i < 10000; i++ {
		_, err := conn.Write(data)
		if err != nil {
			chkErr(t, err)
			break
		}
		_, err = conn.Read(buffer)

		if err != nil {
			chkErr(t, err)
			break
		}
		//n, _ := conn.Read(buffer)
		//w.Write(buffer[:n])
	}
	conn.Close()
}

func chkErr(t *testing.T, err error) {
	t.Fail()
	t.Error(err)
}
