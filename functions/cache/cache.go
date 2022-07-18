package cache

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/akrylysov/pogreb"
	"strconv"
)

func State(id int64, state ...uint16) (uint16, error) {
	db, err := pogreb.Open("cache", nil)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	if len(state) == 0 {
		state, err := db.Get([]byte(strconv.FormatInt(id, 10) + "state"))
		if err != nil {
			return 0, err
		}
		if bytes.Compare(state, []byte{}) == 0 {
			return 0, errors.New("stateless user")
		}
		return binary.BigEndian.Uint16(state), err
	} else if len(state) == 1 {
		b := make([]byte, 4)
		binary.BigEndian.PutUint16(b, state[0])
		err := db.Put([]byte(strconv.FormatInt(id, 10)+"state"), b)
		if err != nil {
			return 0, err
		}
		return state[0], err
	} else {
		return 0, errors.New("too many arguments")
	}
}
