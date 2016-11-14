// Copyright 2016 The FOFA SDK Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package examples shows how to use fofa sdk
package examples

import (
	"encoding/json"
	"fmt"

	"github.com/fofapro/fofa-go/fofa"

	"os"
)

// FofaExample fofa sdk functons included
func FofaExample() {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")

	clt := fofa.NewFofaClient([]byte(email), []byte(key))
	if clt == nil {
		fmt.Printf("create fofa client\n")
		return
	}
	ret, err := clt.QueryAsJSON(1, []byte(`body="小米"`))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	fmt.Printf("%s\n", ret)
	arr, err := clt.QueryAsArray(1, []byte(`domain="163.com"`), []byte("ip,host,title"))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	fmt.Printf("count: %d\n", len(arr))
	encodeArr, _ := json.Marshal(arr)
	fmt.Printf("\n%s\n", encodeArr)
}
