/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-12-26 23:50
 * description :
 * history :
 */

package crypto

import (
	"testing"
)

func Test_A(t *testing.T) {
	cyp := NewUnixCrypto("sonven", "3dsdgfdfgdfg")
	i := 2000000
	for {
		if i = i - 1; i == 0 {
			break
		}

		s := string(cyp.Encode())
		//fmt.Println("str:",s)
		cyp.Compare(s)

		//r,bytes,unix := cyp.Compare(s)
		//fmt.Println("dst:",string(bytes),time.Unix(unix,0).String())
		//fmt.Println("src:",string(cyp.GetBytes()))
		//fmt.Println("result:",r)
	}
}
