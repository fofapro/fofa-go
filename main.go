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
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

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

    The commands are:

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

	var (
		json  []byte
		array fofa.Results
		err   error
	)
	switch *format {
	case "json":
		json, err = clt.QueryAsJSON(uint(*page), []byte(*query), []byte(*fields))
	case "array":
		array, err = clt.QueryAsArray(uint(*page), []byte(*query), []byte(*fields))
	default:
		_ = fmt.Errorf("Expect json or array as output format.")
		usage()
		return
	}

	if err != nil {
		_ = fmt.Errorf("%v\n", err.Error())
	}
	fd, err := os.Create(*out)
	defer fd.Close()
	if err != nil {
		_ = fmt.Errorf("%v\n", err.Error())
	}
	switch {
	case json != nil:
		ioutil.WriteFile(*out, json, 0666)
	case array != nil:
		for _, v := range array {
			if v.Domain != "" {
				io.WriteString(fd, "domain "+v.Domain+"\t")
			}
			if v.Host != "" {
				io.WriteString(fd, "host "+v.Host+"\t")
			}
			if v.IP != "" {
				io.WriteString(fd, "iP "+v.IP+"\t")
			}
			if v.Port != "" {
				io.WriteString(fd, "port "+v.Port+"\t")
			}
			if v.Country != "" {
				io.WriteString(fd, "country "+v.Country+"\t")
			}
			if v.City != "" {
				io.WriteString(fd, "city "+v.City)
			}
			io.WriteString(fd, "\n")
		}
	}
}
