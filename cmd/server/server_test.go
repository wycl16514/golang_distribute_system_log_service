package main

import(
	"encoding/json"
	"net/http"
	"internal/server"
	"bytes"
	"testing"
	"io/ioutil"
)

func TestServerLogWrite(t *testing.T) {
	var tests = []struct {
		request server.WriteRequest
	    want_response server.WriteResponse 
	} {
		{request: server.WriteRequest{server.Record{[]byte(`this is log request 1`), 0}}, 
		 want_response:  server.WriteResponse{Offset: 0, },},
		 {request: server.WriteRequest{server.Record{[]byte(`this is log request 2`), 0}}, 
		 want_response:  server.WriteResponse{Offset: 1, },},
		 {request: server.WriteRequest{server.Record{[]byte(`this is log request 3`), 0}}, 
		 want_response:  server.WriteResponse{Offset: 2, },},
	}

	for _, test := range tests {
		//将请求转换成json格式并post给日志服务
		request := &test.request 
		request_json, err := json.Marshal(request)
		if err != nil {
			t.Errorf("convert request to json fail")
			return 
		}

		resp, err := http.Post("http://localhost:8080", "application/json",bytes.NewBuffer(request_json))
		defer resp.Body.Close()
		if err != nil {
			t.Errorf("http post request fail: %v", err)
			return
		}

		//解析日志服务返回结果
		body, err := ioutil.ReadAll(resp.Body)
		var response server.WriteResponse 
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Errorf("Unmarshal write response fail: %v", err)
		}

		//检测结果是否与预期一致
		if response.Offset != test.want_response.Offset {
			t.Errorf("got offset: %d, but want offset: %d", response.Offset, test.want_response.Offset)
		}
		
	}

	var read_tests = []struct {
		request server.ReadRequest 
		want server.ReadResponse 
	} {
		{request: server.ReadRequest{Offset : 0,}, 
		want: server.ReadResponse{server.Record{[]byte(`this is log request 1`), 0}} },
		{request: server.ReadRequest{Offset : 1,}, 
		want: server.ReadResponse{server.Record{[]byte(`this is log request 2`), 0}} },
		{request: server.ReadRequest{Offset : 2,}, 
		want: server.ReadResponse{server.Record{[]byte(`this is log request 3`), 0}} },
	}
 
	for _, test := range read_tests {
		request := test.request 
		request_json , err := json.Marshal(request)
		if err != nil {
			t.Errorf("convert read request to json fail")
			return 
		}

		//将请求转换为json并放入GET请求体
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", bytes.NewBuffer(request_json))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("read request fail: %v", err)
			return 
		}

		//解析读请求返回的结果
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var response server.ReadResponse
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Errorf("Unmarshal read response fail: %v", err)
			return 
		}

		res := bytes.Compare(response.Record.Value, test.want.Record.Value)
		if res != 0 {
			t.Errorf("got value: %q, but want value : %q", response.Record.Value, test.want.Record.Value)
		}
	}


}

