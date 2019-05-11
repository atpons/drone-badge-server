package db

import (
	"testing"
)

func TestSetValue(t *testing.T) {
	db, _ := NewDb()
	stg := Stage{
		string(9223372036854775807),
		100,
		1,
		true,
	}
	err := SetValue(db, &stg)
	if err != nil {
		t.Fatal("Error SetValue")
	}
}

func TestGetValue(t *testing.T) {
	db, _ := NewDb()
	stg := Stage{
		string(9223372036854775807),
		100,
		1,
		true,
	}
	err := SetValue(db, &stg)
	stg2 := Stage{
		string(9223372036854775807),
		100,
		2,
		true,
	}
	err = SetValue(db, &stg2)
	val, err := GetValue(db, 114, 2)
	t.Log(val)
	if err != nil {
		t.Fatal("Record Not Found")
	}
}
