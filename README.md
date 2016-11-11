# FOFA Pro SDK 使用说明文档
## FOFA Pro API   
<a href="https://fofa.so/api"><font face="menlo">`FOFA Pro API`</font></a> 是资产搜索引擎 <a href="https://fofa.so/">`FOFA Pro`</a> 为开发者提供的 `RESTful API` 接口, 允许开发者在自己的项目中集成 `FOFA Pro` 的功能。    


## FOFA SDK
基于 `FOFA Pro API` 编写的 `golang` 版 `SDK`, 方便 `golang` 开发者快速将 `FOFA Pro` 集成到自己的项目中。


## 环境
### 开发环境
``` zsh
$ go version
go version go1.7.3 darwin/amd64
```
### 测试环境
``` zsh
$ go version
go version go1.7.3 darwin/amd64
```
### 使用环境
建议在 `go version > 1.7.2 && (Windows >= XP || OS X >= 10.8 || Linux >= 2.6.23) && 64-bit processor` 使用 `FOFA SDK`

## 获取
### `govendor`
安装 `FOFA SDK` 之前，请先确认已经安装 <a href="https://github.com/kardianos/govendor/"><font face="menlo">govendor</font></a>.   
``` zsh
$ go get -u -v github.com/kardianos/govendor
```
`点击`  <a href="https://github.com/kardianos/govendor/blob/master/doc/faq.md"><font face="menlo">GOVENDOR FAQ</font></a> `了解更详细的使用说明`      


### FOFA SDK
<strong>下载 `FOFA SDK` </strong>   
``` zsh
$ go get github.com/fofapro/fofa-go
```   
<strong>安装</strong>  
``` zsh
$ govendor install +local
```
<i>如果是 `*nix` 或者 `mac` 可以使用</i>   
``` zsh
$ make
$ make install
```
<strong>文档</strong> 
``` zsh
$ godoc -http=:6060 -index
```
`或者`
``` zsh
$ make doc
```
<strong>测试</strong>   
``` zsh
$ govendor test  ./fofa
```
<i>如果是 `*nix` 或者 `mac` 可以使用</i>   
``` zsh
$ make test
```
## 依赖
### Email & API Key   
| `Email` |用户登陆 `FOFA Pro` 使用的 `Email`|
|---------|:-----------------:|
|`Key`| 前往 <a href="https://fofa.so/my/users/info" style="color:#0000ff"><strong>`个人中心`</strong></a> 查看 `API Key` |
如果开发者经常使用固定的账号，建议将`email`与`key`添加到环境变量中   
``` zsh
echo "export FOFA_EMAIL=example@fofa.so" >> ~/.zshrc
echo "export FOFA_KEY=32charsMD5String" >> ~/.zshrc
```
`SDK` 提供的示例代码就是使用的这种形式。


## Example   
``` go
func FofaExample() {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")

	clt := fofa.NewFofaClient([]byte(email), []byte(key))
	if clt == nil {
		fmt.Printf("create fofa client\n")
		return
	}
	result, err := clt.QueryAsJSON(1, []byte(`body="小米"`))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	fmt.Printf("%s\n", result)
	arr, err := clt.QueryAsArray(1, []byte(`domain="haosec.cn"`), []byte("domain"))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	fmt.Printf("count: %d\n", len(arr))
	encodeArr, _ := json.Marshal(arr)
	fmt.Printf("\n%s\n", encodeArr)
}
```

## FOFA Tool   
使用`FOFA SDK` 开发的控制台程序,无序登陆即可进行`FOFA`查询.    
在同时满足安装`FOFA SDK`与`$GOPATH/bin`包含在`$PATH`环境变量时，可以直接运行.
### USAGE   
``` zsh
$ fofa-go

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
```

## 协议
`FOFA SDK` 遵循 `MIT` 协议 <a href="https://opensource.org/licenses/mit"><font face="menlo">https://opensource.org/licenses/mit</font></a>
