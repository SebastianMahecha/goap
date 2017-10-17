# Goap

It's another golang soap client.

## Usage

First define request body and response body structure :

```go
type addRequest struct {
    XMLName xml.Name `xml:"http://tempuri.org/ Add"`

    IntA int `xml:"intA"`
    IntB int `xml:"intB"`
}

type addResponse struct {
    XMLName xml.Name `xml:"http://tempuri.org/ AddResponse"`

    AddResult int `xml:"AddResult"`
}
```

Then you can use goap to send your request and process response :

```go
response := addResponse{}
err := goap.Call(
    "http://www.dneonline.com/calculator.asmx",  // service url
    "http://tempuri.org/Add", // soap action
    nil, // soap request headers (optional)
    addRequest{IntA: 10, IntB: 15}, // request
    nil, // soap response headers(optional)
    &response) // soap response pointer

if err != nil {
    panic(err)
}

log.Println("add result is", response.AddResult)
```

## Thanks
This project inspired by [Zarinpal go example](https://github.com/sijad/zarinpal-go) by [sijad](https://github.com/sijad)