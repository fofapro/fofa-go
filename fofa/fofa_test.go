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
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"runtime"

	"github.com/buger/jsonparser"
)

var (
	clt = NewFofaClient([]byte(os.Getenv("FOFA_EMAIL")), []byte(os.Getenv("FOFA_KEY")))
)

func EqualBytes(a, b []byte) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewFofaClient(t *testing.T) {
	rand.Seed(time.Now().Unix() ^ 0x1a2b3c4d)
	for i := 0; i < 100; i++ {
		email := strconv.Itoa(rand.Int())
		key := strconv.Itoa(rand.Int())
		clt := NewFofaClient([]byte(email), []byte(key))
		if !EqualBytes([]byte(email), clt.email) || !EqualBytes([]byte(key), clt.key) {
			t.Errorf("expect email = %s  key = %s , but email = %s  key = %s\n", email, key, clt.email, clt.key)
		}
	}
}

func TestQueryAsJSON(t *testing.T) {
	var (
		arr          = []byte(nil)
		err          = error(nil)
		modeNormal   = []byte("normal")
		modeExtended = []byte("extended")
		query        = []byte(nil)
		fields       = []byte(nil)
		page         = uint(0)
	)
	// -------------------------------------------
	clt := NewFofaClient([]byte(os.Getenv("FOFA_EMAIL")), []byte(os.Getenv("FOFA_KEY")))
	if clt == nil {
		t.Errorf("create fofa client failed!")
	}
	{
		{

			{ // ------extended-----domain-----------------------------

				query = []byte(`host="nosec.org"`)
				fields = []byte(`fields=domain`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeExtended, query, page, 6)
				}
			} // -------------------------------------------

			{ // -------extended-----host-----------------------------
				query = []byte(`host="nosec.org"`)
				fields = []byte(`fields=host`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeExtended, query, page, 6)
				}
			} // -------------------------------------------

			{ // -------normal------domain------------------------------
				query = []byte(`nosec.org`)
				fields = []byte(`fields=domain`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeNormal, query, page, 25)
				}
			} // -------------------------------------------

			{ // --------normal------host-----------------------------
				query = []byte(`nosec.org`)
				fields = []byte(`fields=host`)
				page = 1
				arr, err = clt.QueryAsJSON(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					jsonExpectEqual(t, arr, modeNormal, query, page, 25)
				}
			} // -------------------------------------------
		} // -------------------------------------------

		{ // -------------------------------------------
			query = []byte(`host="nosec.org"`)
			fields = nil
			page = 1
			arr, err = clt.QueryAsJSON(page, query, fields)
			if err != nil {
				t.Errorf("%v\n", err.Error())
			} else {
				jsonExpectEqual(t, arr, modeExtended, query, page, 6)
			}
		}

		{ // -------------------------------------------
			query = []byte(`host="nosec.org"`)
			fields = nil
			page = 1
			arr, err = clt.QueryAsJSON(page, query)
			if err != nil {
				t.Errorf("%v\n", err.Error())
			} else {
				jsonExpectEqual(t, arr, modeExtended, query, page, 6)
			}
		} // -------------------------------------------
	} // -------------------------------------------
}

func TestQueryAsArray(t *testing.T) {
	var (
		arr                                          = Results(nil)
		err                                          = error(nil)
		query, fields                                = []byte(nil), []byte(nil)
		page, size                                   = uint(0), uint(0)
		host, title, ip, domain, port, country, city bool
	)
	// -------------------------------------------
	if clt == nil {
		t.Fatalf("create fofa client failed!")
	}
	{
		{ // fields != nil
			{ // fields: domain   query:host=nosec.org
				host, title, ip, domain, port, country, city = false, false, false, true, false, false, false
				query = []byte(`host="nosec.org"`)
				fields = []byte(`domain`)
				page, size = 1, 6
				arr, err = clt.QueryAsArray(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {

					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)
				}
			} // fields: domain   query:host=nosec.org

			{ // fields: host   query:host=nosec.org
				host, title, ip, domain, port, country, city = true, false, false, false, false, false, false

				query = []byte(`host="nosec.org"`)
				fields = []byte(`host`)
				page, size = 1, 6
				arr, err = clt.QueryAsArray(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)
				}
			} // fields: host   query:host=nosec.org

			{ // fields: domain   query:nosec.org
				host, title, ip, domain, port, country, city = false, false, false, true, false, false, false

				query = []byte(`nosec.org`)
				fields = []byte(`domain`)
				page, size = 1, 25
				arr, err = clt.QueryAsArray(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)

				}
			} // fields: domain   query:nosec.org

			{ // fields: host   query:nosec.org
				host, title, ip, domain, port, country, city = true, false, false, false, false, false, false

				query = []byte(`nosec.org`)
				fields = []byte(`host`)
				page, size = 1, 25

				arr, err = clt.QueryAsArray(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)

				}
			} // fields: host   query:nosec.org
		} // fields != nil

		{ // fields == nil
			{ // query:host=nosec.org
				host, title, ip, domain, port, country, city = false, false, false, false, false, false, false
				query = []byte(`host="nosec.org"`)
				fields = nil
				page, size = 1, 6

				arr, err = clt.QueryAsArray(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)
				}
			} // query:host=nosec.org

			{ // query:nosec.org
				host, title, ip, domain, port, country, city = false, false, false, false, false, false, false
				query = []byte(`"nosec.org"`)
				fields = nil
				page, size = 1, 25

				arr, err = clt.QueryAsArray(page, query, fields)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)

				}
			} // query:nosec.org

		} // fields == nil

		{ // without fields
			{ // -------------------------------------------
				host, title, ip, domain, port, country, city = true, true, true, true, true, true, true
				query = []byte(`host="nosec.org"`)
				fields = nil
				page, size = 1, 6

				arr, err = clt.QueryAsArray(page, query)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)

				}
			}

			{ // -------------------------------------------
				host, title, ip, domain, port, country, city = true, true, true, true, true, true, true
				query = []byte(`host="nosec.org"`)
				fields = nil
				page = 1
				page, size = 1, 6

				arr, err = clt.QueryAsArray(page, query)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)

				}
			}
			{ // -------------------------------------------
				host, title, ip, domain, port, country, city = true, true, true, true, true, true, true
				query = []byte(`nosec.org`)
				fields = nil
				page, size = 1, 25

				arr, err = clt.QueryAsArray(page, query)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)

				}
			}
			{ // -------------------------------------------
				host, title, ip, domain, port, country, city = true, true, true, true, true, true, true
				query = []byte(`nosec.org`)
				fields = nil
				page, size = 1, 25

				arr, err = clt.QueryAsArray(page, query)
				if err != nil {
					t.Errorf("%v\n", err.Error())
				} else {
					arrayExpectEqual(t, arr, host, title, ip, domain, port, country, city, size)
				}
			}
		} // without fields
	}
}

func TestServices(t *testing.T) {

}

func TestOnlyOneResult(t *testing.T) {
	// query about will return just one record
	// it maybe occurs an error.
	if clt == nil {
		t.Fatalf("create fofa client failed!")
	}
	var (
		query, fields = []byte(nil), []byte(nil)
		arr           = Results(nil)
		err           = error(nil)
	)
	{ // fields == nil
		query = []byte(`domain="haosec.cn"`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "Beijing" {
			t.Errorf("Expect City Beijing  But %s\n", haosec.City)
		}
		if haosec.Title != "首页 - 安徒生 - 企业威胁情报感知平台" {
			t.Errorf("Expect Title 首页 - 安徒生 - 企业威胁情报感知平台  But %s\n", haosec.Title)
		}
		if haosec.Country != "CN" {
			t.Errorf("Expect Country CN  But %s\n", haosec.Country)
		}
		if haosec.IP != "123.59.94.182" {
			t.Errorf("Expect IP 123.59.94.182  But %s\n", haosec.IP)
		}
		if haosec.Domain != "haosec.cn" {
			t.Errorf("Expect Domain haosec.cn  But %s\n", haosec.Domain)
		}

		if haosec.Host != "https://haosec.cn" {
			t.Errorf("Expect Host https://haosec.cn  But %s\n", haosec.City)
		}
	}

	{ // fields == host
		query = []byte(`domain="haosec.cn"`)
		fields = []byte(`host`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "" {
			t.Errorf("Expect City EMPTY  But %s\n", haosec.City)
		}
		if haosec.Title != "" {
			t.Errorf("Expect Title EMPTY  But %s\n", haosec.Title)
		}
		if haosec.Country != "" {
			t.Errorf("Expect Country EMPTY  But %s\n", haosec.Country)
		}
		if haosec.IP != "" {
			t.Errorf("Expect IP EMPTY  But %s\n", haosec.IP)
		}
		if haosec.Domain != "" {
			t.Errorf("Expect Domain EMPTY  But %s\n", haosec.Domain)
		}

		if haosec.Host != "https://haosec.cn" {
			t.Errorf("Expect Host https://haosec.cn  But %s\n", haosec.City)
		}
	}

	{ // fields == domain
		query = []byte(`domain="haosec.cn"`)
		fields = []byte(`domain`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "" {
			t.Errorf("Expect City EMPTY  But %s\n", haosec.City)
		}
		if haosec.Title != "" {
			t.Errorf("Expect Title EMPTY  But %s\n", haosec.Title)
		}
		if haosec.Country != "" {
			t.Errorf("Expect Country EMPTY  But %s\n", haosec.Country)
		}
		if haosec.IP != "" {
			t.Errorf("Expect IP EMPTY  But %s\n", haosec.IP)
		}
		if haosec.Domain != "haosec.cn" {
			t.Errorf("Expect Domain haosec.cn  But %s\n", haosec.Domain)
		}

		if haosec.Host != "" {
			t.Errorf("Expect Host EMPTY  But %s\n", haosec.City)
		}
	}
	{ // fields == ip
		query = []byte(`domain="haosec.cn"`)
		fields = []byte(`ip`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "" {
			t.Errorf("Expect City EMPTY  But %s\n", haosec.City)
		}
		if haosec.Title != "" {
			t.Errorf("Expect Title EMPTY  But %s\n", haosec.Title)
		}
		if haosec.Country != "" {
			t.Errorf("Expect Country EMPTY  But %s\n", haosec.Country)
		}
		if haosec.IP != "123.59.94.182" {
			t.Errorf("Expect IP 123.59.94.182  But %s\n", haosec.IP)
		}
		if haosec.Domain != "" {
			t.Errorf("Expect Domain EMPTY  But %s\n", haosec.Domain)
		}

		if haosec.Host != "" {
			t.Errorf("Expect Host EMPTY  But %s\n", haosec.City)
		}
	}
	{ // fields == city
		query = []byte(`domain="haosec.cn"`)
		fields = []byte(`city`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "Beijing" {
			t.Errorf("Expect City Beijing  But %s\n", haosec.City)
		}
		if haosec.Title != "" {
			t.Errorf("Expect Title EMPTY  But %s\n", haosec.Title)
		}
		if haosec.Country != "" {
			t.Errorf("Expect Country EMPTY  But %s\n", haosec.Country)
		}
		if haosec.IP != "" {
			t.Errorf("Expect IP EMPTY  But %s\n", haosec.IP)
		}
		if haosec.Domain != "" {
			t.Errorf("Expect Domain EMPTY  But %s\n", haosec.Domain)
		}

		if haosec.Host != "" {
			t.Errorf("Expect Host EMPTY  But %s\n", haosec.City)
		}
	}
	{ // fields == country
		query = []byte(`domain="haosec.cn"`)
		fields = []byte(`country`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "" {
			t.Errorf("Expect City EMPTY  But %s\n", haosec.City)
		}
		if haosec.Title != "" {
			t.Errorf("Expect Title EMPTY  But %s\n", haosec.Title)
		}
		if haosec.Country != "CN" {
			t.Errorf("Expect Country CN  But %s\n", haosec.Country)
		}
		if haosec.IP != "" {
			t.Errorf("Expect IP EMPTY  But %s\n", haosec.IP)
		}
		if haosec.Domain != "" {
			t.Errorf("Expect Domain EMPTY  But %s\n", haosec.Domain)
		}

		if haosec.Host != "" {
			t.Errorf("Expect Host EMPTY  But %s\n", haosec.City)
		}
	}

	{ // fields == domain,city
		query = []byte(`domain="haosec.cn"`)
		fields = []byte(`domain,city`)
		arr, err = clt.QueryAsArray(1, query, fields)
		if err != nil {
			t.Errorf("%v\n", err.Error())
		}
		haosec := arr[0]
		if haosec.City != "Beijing" {
			t.Errorf("Expect City Beijing  But %s\n", haosec.City)
		}
		if haosec.Title != "" {
			t.Errorf("Expect Title EMPTY  But %s\n", haosec.Title)
		}
		if haosec.Country != "" {
			t.Errorf("Expect Country EMPTY  But %s\n", haosec.Country)
		}
		if haosec.IP != "" {
			t.Errorf("Expect IP EMPTY  But %s\n", haosec.IP)
		}
		if haosec.Domain != "haosec.cn" {
			t.Errorf("Expect Domain haosec.cn  But %s\n", haosec.Domain)
		}

		if haosec.Host != "" {
			t.Errorf("Expect Host EMPTY  But %s\n", haosec.City)
		}
	}
}

func TestError(t *testing.T) {
	if clt == nil {
		t.Fatalf("create fofa client failed!")
	}
	var (
		query, json = []byte(nil), []byte(nil)
		//arr         = Results(nil)
		err = error(nil)
	)
	query = []byte("domain=/")
	json, err = clt.QueryAsJSON(1, query)
	if err == nil {
		t.Fatalf("fofa query make an error: %v\n", err.Error())
	}
	if a, b := jsonparser.GetBoolean(json, "error"); b != nil {
		t.Errorf("parse error's json should not make error, but %v\n", b.Error())
	} else if !a {
		t.Errorf("error must be true, but %s\n", json)
	} else {
		_, d := jsonparser.GetString(json, "errmsg")
		if d != nil {
			t.Errorf("parse error's json should not make error, but %v\n", d.Error())
		}
	}
}

func TestIP(t *testing.T) {
	if clt == nil {
		t.Fatalf("create fofa client failed!")
	}

	var (
		query, fields = []byte(nil), []byte(nil)
		arr           = Results(nil)
		err           = error(nil)
	)

	query = []byte(`ip="106.75.75.204"`)
	fields = []byte(`host,domain,ip,port,title,city,country`)
	arr, err = clt.QueryAsArray(1, query, fields)
	switch {
	case err != nil:
		t.Errorf("%v\n", err.Error())
	case len(arr) != 4:
		t.Errorf("expect 4 records, but get %d\n", len(arr))
	}
	for _, v := range arr {
		switch {
		case v.Port == "6379":
			if v.Domain != "" || v.Title != "" {
				t.Errorf("%s\n", v)
			}
		case v.Port == "443":
			switch v.Title {
			case "FOFA Pro - 网络空间安全搜索引擎":
				if v.Domain != "fofa.so" && v.Domain != "106.75.75.204" {
					t.Errorf("%s\n", v)
				}
			}
		}
	}
}

func jsonExpectEqual(t *testing.T, data, mode, query []byte, page, size uint) {
	m := getMode(data)
	q := getQuery(data)
	p := getPage(data)
	s := getSize(data)
	if !EqualBytes(m, mode) || !EqualBytes(q, query) || page != p || size != s {
		_, f, r, _ := runtime.Caller(1)
		t.Errorf("%s %d: Expect\tmode=%s  query=%s  page=%d  size=%d\nBut\tmode=%s  query=%s  page=%d  size=%d\n", f, r, mode, query, page, size, m, q, p, s)
	}
}

func arrayExpectEqual(t *testing.T, data Results, host, title, ip, domain, port, country, city bool, size uint) {
	if !host && !title && !ip && !domain && !port && !country && !city {
		host, title, ip, domain, port, country, city = true, true, true, true, true, true, true
	}
	_, f, r, _ := runtime.Caller(1)
	if len(data) != int(size) {
		t.Errorf("%s %d: Expect count=%d. But count=%d.\n", f, r, size, len(data))
	}
	for _, v := range data {

		if (host && v.Host == "") || (!host && v.Host != "") {
			t.Errorf("%s %d: Expect host=%v\nBut %s\n", f, r, host, v)
		}
		if (title && v.Title == "") || (!title && v.Title != "") {
			t.Errorf("%s %d: Expect title=%v\nBut %s\n", f, r, title, v)
		}
		if (ip && v.IP == "") || (!ip && v.IP != "") {
			t.Errorf("%s %d: Expect ip=%v\nBut %s\n", f, r, ip, v)
		}
		if (domain && v.Domain == "") || (!domain && v.Domain != "") {
			t.Errorf("%s %d: Expect domain=%v\nBut %s\n", f, r, domain, v)
		}
		// TODO ...
		// 由于 http协议的返回结果很有可能没有 port 字段的内容，所以无法对这个字段的返回值做测试
		// if (port && v.Port == "") || (!port && v.Port != "") {
		// 	t.Errorf("%s %d: Expect port=%v\nBut port=%s\n", f, r, port, v.Port)
		// }
		if (country && v.Country == "") || (!country && v.Country != "") {
			t.Errorf("%s %d: Expect country=%v\nBut %s\n", f, r, country, v)
		}
		if (city && v.City == "") || (!city && v.City != "") {
			t.Errorf("%s %d: Expect city=%v\nBut %s\n", f, r, city, v)
		}
	}
}

func getDomain(data []byte) []byte {
	domain, _ := jsonparser.GetString(data, "domain")
	return []byte(domain)
}
func getHost(data []byte) []byte {
	host, _ := jsonparser.GetString(data, "host")
	return []byte(host)
}
func getIP(data []byte) []byte {
	ip, _ := jsonparser.GetString(data, "ip")
	return []byte(ip)
}
func getPort(data []byte) []byte {
	port, _ := jsonparser.GetString(data, "port")
	return []byte(port)
}
func getCountry(data []byte) []byte {
	country, _ := jsonparser.GetString(data, "country")
	return []byte(country)
}
func getCity(data []byte) []byte {
	city, _ := jsonparser.GetString(data, "city")
	return []byte(city)
}
func getPage(data []byte) uint {
	page, _ := jsonparser.GetInt(data, "page")
	return uint(page)
}
func getMode(data []byte) []byte {
	mode, _ := jsonparser.GetString(data, "mode")
	return []byte(mode)
}
func getQuery(data []byte) []byte {
	query, _ := jsonparser.GetString(data, "query")
	return []byte(query)
}
func getSize(data []byte) uint {
	size, _ := jsonparser.GetInt(data, "size")
	return uint(size)
}
