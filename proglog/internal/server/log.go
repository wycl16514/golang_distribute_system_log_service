package server

import (
	"fmt"
	"sync"
)

type Log struct {
	ch_mux chan interface{} 
	records [] Record 
}

func NewLog() *Log {
	return &Log{ch : make(chan Record),} 
}

func(c *Log) Append(record Record) (uint64, error) {
	c.ch_mux <- struct{}{}
    defer <- c.ch_mux 
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	return record.Offset, nil 
}

func (c *Log) Read(offset uint64)(Record, error) {
	c.ch_mux <- struct{}{}
    defer <- c.ch_mux 

	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound 
	}

	return c.records[offset], nil 
}

type Record struct {
	Value []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")