package main


import (
	"github.com/juju/errors"
	"strconv"
)

// Erlang external term tags.
const (
	ettAtom          = 'd' //100
	ettAtomUTF8      = 'v' //118 // this is beyond retarded
	ettBinary        = 'm' //109
	ettBitBinary     = 'M' //77
	ettCachedAtom    = 'C' //67
	ettCacheRef      = 'R' //82
	ettExport        = 'q' //113
	ettFloat         = 'c' //99
	ettFun           = 'u' //117
	ettInteger       = 'b' //98
	ettLargeBig      = 'o' //111
	ettLargeTuple    = 'i' //105
	ettList          = 'l' //108
	ettNewCache      = 'N' //78
	ettNewFloat      = 'F' //70
	ettNewFun        = 'p' //112
	ettNewRef        = 'r' //114
	ettNil           = 'j' //106
	ettPid           = 'g' //103
	ettPort          = 'f' //102
	ettRef           = 'e' //101
	ettSmallAtom     = 's' //115
	ettSmallAtomUTF8 = 'w' //119 // this is beyond retarded
	ettSmallBig      = 'n' //110
	ettSmallInteger  = 'a' //97
	ettSmallTuple    = 'h' //104
	ettString        = 'k' //107
)

const (
	// Erlang external term format version
	EtVersion = byte(131)
	EtCompress = byte(80)
)

const (
	// Erlang distribution header
	EtDist = byte('D')
)

var tagNames = map[byte]string{
	ettAtom:          "ATOM_EXT",
	ettAtomUTF8:      "ATOM_UTF8_EXT",
	ettBinary:        "BINARY_EXT",
	ettBitBinary:     "BIT_BINARY_EXT",
	ettCachedAtom:    "ATOM_CACHE_REF",
	ettExport:        "EXPORT_EXT",
	ettFloat:         "FLOAT_EXT",
	ettFun:           "FUN_EXT",
	ettInteger:       "INTEGER_EXT",
	ettLargeBig:      "LARGE_BIG_EXT",
	ettLargeTuple:    "LARGE_TUPLE_EXT",
	ettList:          "LIST_EXT",
	ettNewCache:      "NEW_CACHE_EXT",
	ettNewFloat:      "NEW_FLOAT_EXT",
	ettNewFun:        "NEW_FUN_EXT",
	ettNewRef:        "NEW_REFERENCE_EXT",
	ettNil:           "NIL_EXT",
	ettPid:           "PID_EXT",
	ettPort:          "PORT_EXT",
	ettRef:           "REFERENCE_EXT",
	ettSmallAtom:     "SMALL_ATOM_EXT",
	ettSmallAtomUTF8: "SMALL_ATOM_UTF8_EXT",
	ettSmallBig:      "SMALL_BIG_EXT",
	ettSmallInteger:  "SMALL_INTEGER_EXT",
	ettSmallTuple:    "SMALL_TUPLE_EXT",
	ettString:        "STRING_EXT",
}

type Term interface {}

type Tuple struct {
	data []Term
}

type ErlangList struct {
	data []Term
	tail Term
}

//Override
func (t *Tuple) getType() byte {
	if t.isSmall() {
		return ettSmallTuple
	} else {
		return ettLargeTuple
	}
}

func (t *Tuple) isSmall() bool {
	return len(t.data) <= 255
}

func (t *Tuple) get(index int) Term {
	return t.data[index]
}

func (t *Tuple) size() int {
	return len(t.data)
}

//Override
func (l *ErlangList) getType() byte {
	return ettList
}

func (l *ErlangList) size() int {
	if l.isProper() {
		return len(l.data)
	} else {
		return len(l.data) + 1
	}
}

func (l *ErlangList) isProper() bool {
	if l.tail != nil {
		return false
	}
	return true
}

func (l *ErlangList) get(index int) Term {
	return l.data[index]
}

func StringToUuid(uuid string) (result []byte, err error) {
	if len(uuid) != 36 {
		return nil, errors.New("wrong format of uuid")
	}
	if string(uuid[8]) != "-" || string(uuid[13]) != "-" || string(uuid[18]) != "-" || string(uuid[23]) != "-" {
		return nil, errors.New("wrong format of uuid")
	}
	result = make([]byte, 16)
	var offset int = 0
	for i := 0; i < 16; i++ {
		if offset == 8 || offset == 13 || offset == 18 || offset == 23 {
			offset++
		}
		upper, _ := strconv.ParseInt(string(uuid[offset]), 16, 32)
		offset++
		lower, _ := strconv.ParseInt(string(uuid[offset]), 16, 32)
		offset++
		result[i] = byte(upper * 16 + lower)
	}
	return result, nil
}