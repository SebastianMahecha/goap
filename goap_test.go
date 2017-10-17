package goap

import (
	"encoding/xml"
	"testing"
)

func TestCall(t *testing.T) {
	type args struct {
		url             string
		action          string
		requestHeaders  interface{}
		request         interface{}
		responseHeaders interface{}
		response        interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "add",
			args: args{
				url:      "http://www.dneonline.com/calculator.asmx",
				action:   "http://tempuri.org/Add",
				request:  addRequest{IntA: 10, IntB: 15},
				response: &addResponse{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Call(tt.args.url, tt.args.action, tt.args.requestHeaders, tt.args.request, tt.args.responseHeaders, tt.args.response); err != nil {
				t.Fatalf("Call() error = %v", err)
			}

			if tt.args.response == nil {
				t.Fatalf("response is null")
			}

			result := tt.args.response.(*addResponse).AddResult
			mustBe := tt.args.request.(addRequest).IntA + tt.args.request.(addRequest).IntB
			if result != mustBe {
				t.Errorf("result must be %v but is %v", mustBe, result)
			}
		})
	}
}

type addRequest struct {
	XMLName xml.Name `xml:"http://tempuri.org/ Add"`

	IntA int `xml:"intA"`
	IntB int `xml:"intB"`
}

type addResponse struct {
	XMLName xml.Name `xml:"http://tempuri.org/ AddResponse"`

	AddResult int `xml:"AddResult"`
}
