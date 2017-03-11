package parser

import (
	"bytes"
	"testing"
)

func TestParseInt(t *testing.T) {

	input := []byte("i123e")
	reader := bytes.NewReader(input)

	item, err := parse(reader)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	res, ok := item.(int64)

	if !ok {
		t.Error("Expected integer")
		return
	}

	if res != 123 {
		t.Error("Expected", 123, " but found ", res)
		return
	}
}

func TestParseString(t *testing.T) {
	input := []byte("5:Hello")
	reader := bytes.NewReader(input)

	item, err := parse(reader)

	if err != nil {
		t.Error(err.Error())
		return
	}

	res, ok := item.(string)

	if !ok {
		t.Errorf("Expected string")
		return
	}

	if res != "Hello" {
		t.Error("Expected \"", "Hello", "\" but found \"", res, "\".")
		return
	}
}

func TestParseList(t *testing.T) {
	input := []byte("li12e5:Helloe")
	reader := bytes.NewReader(input)

	item, err := parse(reader)

	if err != nil {
		t.Error(err.Error())
		return
	}

	res, ok := item.([]Item)

	if !ok {
		t.Errorf("Expected list")
		return
	}

	item1, ok := res[0].(int64)
	if !ok {
		t.Errorf("Expected int as first item")
		return
	}

	item2, ok := res[1].(string)
	if !ok {
		t.Errorf("Expected int as first item")
		return
	}

	if item1 != 12 {
		t.Error("Expected ", 12, " but found ", item1)
		return
	}

	if item2 != "Hello" {
		t.Error("Expected \"", "Hello", "\" but found \"", item2, "\".")
		return
	}
}
