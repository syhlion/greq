# restful web service client

restful web service reqeust tool

## Requirements

* [requestwork](https://github.com/syhlion/requestwork.v2)


## Install

`go get github.com/syhlion/restclient`


## Usage

```
func main(){

    //need import https://github.com/syhlion/requestwork.v2
    worker:=requestwork.New(50)

    //GET
    client:=restclient.New(worker,15*time.Second)
    data,httpstatus,err:=client.Get("https://tw.yahoo.com",nil)

    //POST
    v := url.Values{}
    v.Add("data", string(data))
    data,httpstatus,err:=client.Post("https://tw.yahoo.com",bytes.NewBufferString(v.Encode()))

}
```
