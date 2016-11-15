// Copyright (c) 2016 baimaohui

// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.

// package main
// a tool developed by fofa sdk and
// an example for developer how to use fofa-sdk.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/fofapro/fofa-go/fofa"
)

var (
	fields = flag.String("fields", "host,ip,port", "fields which you want to select")
	query  = flag.String("query", "", "query string")
	email  = flag.String("email", os.Getenv("FOFA_EMAIL"), "an email which you login to fofa.so")
	key    = flag.String("key", os.Getenv("FOFA_KEY"), "md5 string which you can find on userinfo page")
	format = flag.String("format", "json", "output format")
	page   = flag.Int("page", 1, "page number you want to query")
	out    = flag.String("out", "fofa_"+strconv.FormatInt(time.Now().Unix(), 10), "output file path")
)

func usage() {
	fmt.Println(`
    Fofa is a tool for discovering assets.

    Usage:

            fofa option argument ...

    The options are:

            email           the email which you login to fofa.so
                            Use FOFA_EMAIL env by default.

            key             the md5 string which you can find on userinfo page
                            Use FOFA_KEY env by default.

            fields          fields which you want to select
                            Use host,ip,port as default.

            query           query statement which is similar to the statement used in the fofa.so

            format          output format
                            Json(json) as default, alternatively you can select array(array).

            page            page number you want to query, 1000 records per page
                            If page is not set or page is less than 1, page will be set to 1.

            out             output file path
                            Use fofa_${timestamp} as default.
    `)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if *email == "" {
		*email = os.Getenv("FOFA_EMAIL")
	}
	if *key == "" {
		*key = os.Getenv("FOFA_KEY")
	}
	clt := fofa.NewFofaClient([]byte(*email), []byte(*key))
	if clt == nil {
		_ = fmt.Errorf("Allocate Failed! Out Of Memery!\n")
		return
	}
	if *page <= 0 {
		*page = 1
	}

	if *query == "" {
		usage()
		return
	}

	*out = *out + ".json"

	var (
		jsonResult  []byte
		arrayResult fofa.Results
		err         error
	)
	switch *format {
	case "json":
		jsonResult, err = clt.QueryAsJSON(uint(*page), []byte(*query), []byte(*fields))
	case "array":
		arrayResult, err = clt.QueryAsArray(uint(*page), []byte(*query), []byte(*fields))
	default:
		_ = fmt.Errorf("Expect json or array as output format.")
		usage()
		return
	}

	if err != nil {
		_ = fmt.Errorf("%v\n", err.Error())
	}
	switch {
	case jsonResult != nil:
		ioutil.WriteFile(*out, jsonResult, 0666)
	case arrayResult != nil:
		marshalResult, err := json.Marshal(arrayResult)
		if err != nil {
			fmt.Printf("[Fatal] %s\n", err.Error())
			return
		}
		ioutil.WriteFile(*out, marshalResult, 0666)
	}
}
