package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// implement 2 functions to help us in parsing process
// readline and read integer

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}

		n += 1

		line = append(line, b)

		// check if the line finished or not
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// readArray function
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// Read the length of the array
	_len, _, err := r.readInteger()

	if err != nil {
		return Value{}, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, 0)
	for i := 0; i < _len; i++ {
		val, err := r.Read()

		if err != nil {
			return v, err
		}

		// append parsed value to the array
		v.array = append(v.array, val)
	}
	return v, nil
}

// Read bulk function

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	_len, _, err := r.readInteger()

	if err != nil {
		return v, err
	}

	bulk := make([]byte, _len) // array of char (string)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// to read \r\n
	r.readLine()

	return v, nil
}

// marshal function 
func (v Value) Marshal () []byte {
	switch v.typ{
		case "string":
			return v.marshalString()
		case "bulk":
			return v.marshalBulk()
		case "array":
			return v.marshalArray()
		case "error":
			return v.marshalError()
		case "null":
			return v.marshalNull()
		default:
			return []byte{}	
	}
}

// marshal string function 
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

//marshal bulk function
func(v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// Marshal Array function 
func (v Value) marshalArray() []byte {
	lenArray := len(v.array)
	var bytes[]byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(lenArray)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i< lenArray; i ++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

// Marshal Error function 
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// Marshal null function 
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// writer 
type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)

	if err != nil {
		return err
	}
	return nil
}