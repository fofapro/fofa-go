// Copyright (c) 2016 baimaohui

// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.

// Package fofa implements some fofa-api utility functions.
package fofa

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"bytes"

	"github.com/DivineRapier/go-tools/log"

	"strings"

	"github.com/buger/jsonparser"
)

// Fofa a fofa client can be used to make queries
type Fofa struct {
	email []byte
	key   []byte
	*http.Client
}

// Result represents a record of the query results
// contain domain host  ip  port title country city
type result struct {
	Domain  string `json:"domain,omitempty"`
	Host    string `json:"host,omitempty"`
	IP      string `json:"ip,omitempty"`
	Port    string `json:"port,omitempty"`
	Title   string `json:"title,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
}

// User struct for fofa user
type User struct {
	Email  string `json:"email"`
	Fcoin  int    `json:"fcoin"`
	Vip    bool   `json:"bool"`
	Avatar string `json:"avatar"`
}

// Results fofa result set
type Results []result

const (
	defaultAPIUrl = "https://fofa.so/api/v1/search/all?"
)

var (
	errFofaReplyWrongFormat = errors.New("Fofa Reply With Wrong Format.")
	errFofaReplyNoData      = errors.New("No Data In Fofa Reply.")
)

// NewFofaClient create a fofa client
func NewFofaClient(email, key []byte) *Fofa {

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Fofa{
		email: email,
		key:   key,
		Client: &http.Client{
			Transport: transCfg, // disable tls verify
		},
	}
}

// QueryAsJSON make a fofa query and return json data as result
// echo 'domain="nosec.org"' | base64 - | xargs -I{}
// curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64={}"
// host title ip domain port country city
func (ff *Fofa) QueryAsJSON(page uint, args ...[]byte) ([]byte, error) {
	var (
		query  = []byte(nil)
		fields = []byte("domain,host,ip,port,title,country,city")
		q      = []byte(nil)
	)
	switch {
	case len(args) == 1 || (len(args) == 2 && args[1] == nil):
		query = args[0]
	case len(args) == 2:
		query = args[0]
		fields = args[1]
	}

	q = []byte(base64.StdEncoding.EncodeToString(query))
	q = bytes.Join([][]byte{[]byte(defaultAPIUrl),
		[]byte("email="), ff.email,
		[]byte("&key="), ff.key,
		[]byte("&qbase64="), q,
		[]byte("&fields="), fields,
		[]byte("&page="), []byte(strconv.Itoa(int(page))),
	}, []byte(""))
	fmt.Printf("%s\n", q)
	resp, err := ff.Get(string(q))
	if err != nil {
		log.Errorf("%v\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	errmsg, err := jsonparser.GetString(buf, "errmsg")
	if err == nil {
		err = errors.New(errmsg)
	} else {
		err = nil
	}
	return buf, err
}

// QueryAsArray make a fofa query and
// return array data as result
// echo 'domain="nosec.org"' | base64 - | xargs -I{}
// curl "https://fofa.so/api/v1/search/all?email=${FOFA_EMAIL}&key=${FOFA_KEY}&qbase64={}"
func (ff *Fofa) QueryAsArray(page uint, args ...[]byte) (Results, error) {

	var (
		mapFields   = make(map[string]int)
		resultArray = [][]byte(nil)
	)

	data, err := ff.QueryAsJSON(page, args...)
	if err != nil {
		log.Errorf("err: %v\n", err.Error())
		return nil, err
	}

	// map field to index
	if len(args) > 1 && args[1] != nil {
		fields := strings.Split(string(args[1]), ",")
		for i, field := range fields {
			mapFields[strings.Trim(field, " ")] = i
		}
	} else {
		mapFields["domain"] = 0
		mapFields["host"] = 1
		mapFields["ip"] = 2
		mapFields["port"] = 3
		mapFields["title"] = 4
		mapFields["country"] = 5
		mapFields["city"] = 6

	}

	errmsg, err := jsonparser.GetString(data, "errmsg")
	// err equals to nil on error
	if err == nil {
		err = errors.New(errmsg)
		log.Errorf("err: %v\n", errmsg)
		return nil, errors.New(errmsg)
	}

	results, dataType, _, err := jsonparser.Get(data, "results")

	switch {
	case dataType != jsonparser.Array:
		log.Errorf("err: %v\n", err.Error())
		return nil, err
	case err != nil:
		log.Errorf("err: %v\n", err.Error())
		return nil, err
	}
	size, err := jsonparser.GetInt(data, "size")
	if err != nil {
		log.Errorf("fofa reply with wrong format.\n%s\n", data)
		return nil, errFofaReplyWrongFormat
	}
	if size < 1 {
		log.Errorf("no data in fofa reply.\n%s\n", data)
		return nil, errFofaReplyNoData
	}
	if len(mapFields) > 1 {
		resultArray = bytes.Split(results[2:len(results)-2], []byte("],["))
	} else {
		resultArray = bytes.Split(results[1:len(results)-1], []byte{','})
	}
	queryArray := make(Results, len(resultArray), len(resultArray))
	for i, v := range resultArray {
		tmp := bytes.Split(v, []byte{','})

		if a, ok := mapFields["domain"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].Domain = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
		if a, ok := mapFields["host"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].Host = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
		if a, ok := mapFields["ip"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].IP = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
		if a, ok := mapFields["port"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].Port = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
		if a, ok := mapFields["title"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].Title = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
		if a, ok := mapFields["country"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].Country = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
		if a, ok := mapFields["city"]; ok {
			if len(tmp[a]) > 2 {
				queryArray[i].City = string(tmp[a][1 : len(tmp[a])-1])
			}
		}
	}
	return queryArray, nil
}

// UserInfo get user information
func (ff *Fofa) UserInfo() (user *User, err error) {
	user = new(User)
	queryStr := strings.Join([]string{"https://fofa.so/api/v1/info/my?email=", string(ff.email), "&key=", string(ff.key)}, "")
	resp, e := ff.Get(queryStr)
	if e != nil {
		user = nil
		err = e
	} else {
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			err = e
		} else {
			errmsg, e := jsonparser.GetString(body, "errmsg")
			if e == nil {
				err = errors.New(errmsg)
			} else {
				email, err1 := jsonparser.GetString(body, "email")
				fcoin, err2 := jsonparser.GetInt(body, "fcoin")
				vip, err3 := jsonparser.GetBoolean(body, "isvip")
				avatar, err4 := jsonparser.GetString(body, "avatar")
				if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
					err = errors.New("illegal format userinfo")
				} else {
					user.Avatar = avatar
					user.Email = email
					user.Fcoin = int(fcoin)
					user.Vip = vip
				}
			}
		}
	}
	return
}

func (u *User) String() string {
	data, err := json.Marshal(u)
	if err != nil {
		log.Error("json marshal failed. err:", err)
		return ""
	}
	return string(data)
}

func (r *result) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Error("json marshal failed. err:", err)
		return ""
	}
	return string(data)
}

func (r *Results) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Error("json marshal failed. err:", err)
		return ""
	}
	return string(data)
}
