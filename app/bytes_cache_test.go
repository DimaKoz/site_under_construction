package app

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestAddKeyAndPath(t *testing.T) {

	err := AddKeyAndPath(HtmlUnderConstruction, "./../html/under_construction.html")
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile("./../html/under_construction.html")
	if err != nil {
		t.Fatal(err)
	}
	data, err := GetBytes(HtmlUnderConstruction)
	if err != nil {
		t.Errorf("bytes_cache returned unexpected body: got \n %v \nwant\n %v",
			err.Error(), "[]byte")
		return
	}

	if bytes.Compare(expected, *data) != 0 {
		t.Errorf("bytes_cache returned unexpected body: got \n %v \nwant\n %v",
			string(*data), string(expected))
	}
}

func TestAddKeyAndPathError(t *testing.T) {
	RemoveKey(HtmlUnderConstruction)
	err := AddKeyAndPath(HtmlUnderConstruction, "./html/under_construction.html")
	if err == nil {
		t.Errorf("bytes_cache didn't return error: got \n %v \nwant\n %v",
			"error", "nil")
	}
}

func TestGetBytes(t *testing.T) {
	expected, err := ioutil.ReadFile("./../html/under_construction.html")
	if err != nil {
		t.Fatal(err)
	}
	result, err := GetBytes("./../html/under_construction.html")
	if err != nil {
		t.Errorf("bytes_cache returned unexpected body: got \n %v \nwant\n %v",
			err.Error(), "[]byte")
		return
	}

	if bytes.Compare(expected, *result) != 0 {
		t.Errorf("bytes_cache returned unexpected body: got \n %v \nwant\n %v",
			string(*result), string(expected))
	}

}

func TestGetBytesError(t *testing.T) {
	_, err := GetBytes("")
	if err == nil {
		t.Errorf("bytes_cache didn't return error: got \n %v \nwant\n %v",
			"error", "nil")
	}

}
