package main

import (
	"io"
	"math/big"
	"fmt"
	"encoding/binary"
	"bytes"
	"strconv"
	"compress/zlib"
)

type ErrUnknownTerm struct {
	termType byte
}

var (
	ErrFloatScan = fmt.Errorf("read: failed to sscanf float")
	be           = binary.BigEndian
	bTrue        = []byte("true")
	bFalse       = []byte("false")
)

func Read (r io.Reader) (term Term, err error) {
	var etype byte
	if etype, err = ruint8(r); err != nil {
		return nil, err
	}
	fmt.Println("etype", etype)
	if etype == EtCompress {
		if originLen, err := ruint32(r); err != nil {
			return nil, err
		} else {
			fmt.Println("origin length", originLen)
			if dr, err := zlib.NewReader(r); err != nil {
				dr.Close()
				return nil, err
			} else {
				if term, err = Read(dr); err != nil {
					dr.Close()
					return nil, err
				} else {
					dr.Close()
					return term, nil
				}
			}
		}
	}
	if etype == EtVersion {
		if term, err = Read(r); err != nil {
			return nil, err
		} else {
			return term, nil
		}
	}
	var b []byte
	switch etype {
	case ettSmallTuple:
		var arity uint8
		if arity, err = ruint8(r); err != nil {
			break
		}
		tupleData := make([]Term, arity)
		fmt.Println("tuple size", arity)
		for i := 0; i < cap(tupleData); i++ {
			if tupleData[i], err = Read(r); err != nil {
				break
			}
		}
		t := &Tuple{tupleData}
		term = t

	case ettList:
		var n uint32
		if n, err = ruint32(r); err != nil {
			return
		}
		list := make([]Term, n)
		fmt.Println("list size", n)
		for i:=0; i < cap(list); i++ {
			if list[i], err = Read(r); err != nil {
				return
			}
		}
		var tail Term
		if isNil(r) {
			//do nothing
			//TODO check logic
		} else {
			//TODO check logic
		}
		term = &ErlangList{list, tail}

	case ettBinary:
		if b, err = buint32(r); err == nil {
			_, err = io.ReadFull(r, b)
			term = b
		}

	case ettSmallBig:
		b = make([]byte, 2)
		if _, err = io.ReadFull(r, b); err != nil {
			break
		}
		sign := b[1]
		b = make([]byte, uint8(b[0]))
		term, err = readBigInt(r, b, sign)

	case ettNil:
		term = nil
	}
	return
}

func ToString(term Term, buffer *bytes.Buffer) {
	if term == nil {
		buffer.WriteString("null")
	} else {
		switch term.(type) {
		case *Tuple:
			var items []Term
			items = term.(*Tuple).data
			for i := 0; i < len(items); i++ {
				if i == 0 {
					buffer.WriteString("{")
					ToString(items[i], buffer)
					if len(items) == 1 {
						buffer.WriteString("}")
					} else {
						buffer.WriteString(",")
					}
				} else if i == len(items) - 1 {
					ToString(items[i], buffer)
					buffer.WriteString("}")
				} else {
					ToString(items[i], buffer)
					buffer.WriteString(",")
				}
			}
		case *ErlangList:
			var items []Term
			items = term.(*ErlangList).data
			for i := 0; i < len(items); i++ {
				if i == 0 {
					buffer.WriteString("[")
					ToString(items[i], buffer)
					if len(items) == 1 {
						buffer.WriteString("]")
					} else {
						buffer.WriteString(",")
					}
				} else if i == len(items) - 1 {
					ToString(items[i], buffer)
					buffer.WriteString("]")
				} else {
					ToString(items[i], buffer)
					buffer.WriteString(",")
				}
			}
		case int64:
			buffer.WriteString(strconv.FormatInt(term.(int64), 64))
		case int:
			buffer.WriteString(strconv.Itoa(term.(int)))
		case []byte:
			buffer.WriteString(string(term.([]byte)))
		}
	}
}

func isNil(r io.Reader) bool {
	if b, err := ruint8(r); err != nil {
		return false
	} else {
		if b == ettNil {
			return true
		}
	}
	return false
}

func readBigInt(r io.Reader, b []byte, sign byte) (interface{}, error) {
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}

	size := len(b)
	hsize := size >> 1
	for i := 0; i < hsize; i++ {
		b[i], b[size-i-1] = b[size-i-1], b[i]
	}

	v := new(big.Int).SetBytes(b)
	if sign != 0 {
		v = v.Neg(v)
	}

	// try int and int64
	v64 := v.Int64()
	if x := int(v64); v.Cmp(big.NewInt(int64(x))) == 0 {
		return x, nil
	} else if v.Cmp(big.NewInt(v64)) == 0 {
		return v64, nil
	}

	return v, nil
}

func ruint8(r io.Reader) (uint8, error) {
	b := []byte{0}
	_, err := io.ReadFull(r, b)
	return b[0], err
}

func ruint16(r io.Reader) (uint16, error) {
	b := []byte{0, 0}
	_, err := io.ReadFull(r, b)
	return be.Uint16(b), err
}

func ruint32(r io.Reader) (uint32, error) {
	b := []byte{0, 0, 0, 0}
	_, err := io.ReadFull(r, b)
	return be.Uint32(b), err
}

func buint8(r io.Reader) ([]byte, error) {
	size, err := ruint8(r)
	return make([]byte, size), err
}

func buint16(r io.Reader) ([]byte, error) {
	size, err := ruint16(r)
	return make([]byte, size), err
}

func buint32(r io.Reader) ([]byte, error) {
	size, err := ruint32(r)
	return make([]byte, size), err
}