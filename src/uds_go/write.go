package main

import (
	"io"
	"math"
	"fmt"
	"math/big"
	"compress/zlib"
	"bytes"
)

func WriteTerm(w io.Writer, term Term, compress bool) (err error) {
	_, err = w.Write([]byte{EtVersion})
	b := &bytes.Buffer{}
	err = Write(b, term)
	if compress {
		WriteCompress(w, b.Bytes())
	} else {
		w.Write(b.Bytes())
	}
	return
}

func WriteCompress(w io.Writer, o []byte) (err error) {
	n := len(o)
	_, err = w.Write([]byte{
		EtCompress,
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
	nw := zlib.NewWriter(w)
	_ ,err = nw.Write(o)
	nw.Close()
	return
}

func Write(b *bytes.Buffer, term Term) (err error) {
	switch v := term.(type) {
	case *Tuple:
		err = writeTuple(b, v)
	case []byte:
		err = writeBinary(b, v)
	case *ErlangList:
		err = writeList(b, v)
	case int:
		err = writeInt(b, int64(v))
	case nil:
		err = writeNil(b)
	}
	return
}

func writeTuple(b *bytes.Buffer, tuple *Tuple) (err error) {
	n := len(tuple.data)
	if n <= math.MaxUint8 {
		_, err = b.Write([]byte{ettSmallTuple, byte(n)})
	} else {
		_, err = b.Write([]byte{
			ettLargeTuple,
			byte(n >> 24),
			byte(n >> 16),
			byte(n >> 8),
			byte(n),
		})
	}
	if err != nil {
		return
	}
	for _, v := range tuple.data {
		if err = Write(b, v); err != nil {
			return
		}
	}

	return

}

func writeBinary(b *bytes.Buffer, bytes []byte) (err error) {
	switch size := int64(len(bytes)); {
	case size <= math.MaxUint32:
		data := []byte{
			ettBinary,
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		}
		if _, err = b.Write(data); err == nil {
			_, err = b.Write(bytes)
		}

	default:
		err = fmt.Errorf("bad binary size (%d)", size)
	}

	return
}

func writeList(b *bytes.Buffer, l *ErlangList) (err error) {
	var items []Term = l.data
	n := len(items)
	_, err = b.Write([]byte{
		ettList,
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
	if err != nil {
		return
	}
	for i := 0; i < n; i++ {
		v := items[i]
		if err = Write(b, v); err != nil {
			return
		}
	}
	_, err = b.Write([]byte{ettNil})
	return
}


func writeInt(b *bytes.Buffer, x int64) (err error) {
	switch {
	case x >= 0 && x <= math.MaxUint8:
		_, err = b.Write([]byte{ettSmallInteger, byte(x)})

	case x >= math.MinInt32 && x <= math.MaxInt32:
		x := int32(x)
		_, err = b.Write([]byte{
			ettInteger,
			byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x),
		})

	default:
		err = writeBigInt(b, big.NewInt(x))
	}

	return
}

func writeBigInt(b *bytes.Buffer, x *big.Int) (err error) {
	sign := 0
	if x.Sign() < 0 {
		sign = 1
	}

	bytes := reverse(new(big.Int).Abs(x).Bytes())

	switch size := int64(len(bytes)); {
	case size <= math.MaxUint8:
		_, err = b.Write([]byte{ettSmallBig, byte(size), byte(sign)})

	case size <= math.MaxUint32:
		_, err = b.Write([]byte{
			ettLargeBig,
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
			byte(sign),
		})

	default:
		err = fmt.Errorf("bad big int size (%d)", size)
	}

	if err == nil {
		_, err = b.Write(bytes)
	}

	return
}

func writeNil(b *bytes.Buffer) (err error) {
	_, err = b.Write([]byte{ettNil})
	return
}

func reverse(b []byte) []byte {
	size := len(b)
	hsize := size >> 1

	for i := 0; i < hsize; i++ {
		b[i], b[size-i-1] = b[size-i-1], b[i]
	}

	return b
}
