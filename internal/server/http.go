package server 

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func NewHttpServer(addr string) *http.Server {
	httpsrv := newHttpServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleLogWrite).Methods("POST")
	r.HandleFunc("/", httpsrv.hadnleLogRead).Methods("GET")

	return &http.Server{
		Addr : addr,
		Handler: r,
	}
}

type httpServer struct{
	Log *Log 
}

func newHttpServer() *httpServer {
	return &httpServer {
		Log: NewLog(),
	}
}

type WriteRequest struct {
	Record Record `json:"record"`
}

type WriteResponse struct {
	Offset uint64 `json:"offset"`
}

type ReadRequest struct {
	Offset uint64 `json:"offset"`
}

type ReadResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleLogWrite(w http.ResponseWriter, r * http.Request) {
	var req WriteRequest 
	//服务以json格式接收请求
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	res := WriteResponse{Offset: off}
	//服务以json格式返回结果
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
}

func (s *httpServer) hadnleLogRead(w http.ResponseWriter, r *http.Request) {
	var req ReadRequest 
	//服务以json格式接收请求
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	res := ReadResponse{Record: record}
	//服务以json格式返回结果
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}