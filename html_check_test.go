package main

import (
	"UrlParser-1/html_check"
	"fmt"
	"reflect"
	"testing"
)

func TestHtml_check(t *testing.T) {

	TestTable := []struct {
		filename string
		result   []html_check.Info
	}{
		{
			"static/test1.txt",
			[]html_check.Info{{5, 39, "Invalid protocol!"},
				{6, 8, "Invalid protocol!"}},
		},
		{
			"static/test2.txt",
			[]html_check.Info{{1, 39, "Too many double slashes!"}},
		},
		{
			"static/test3.txt",
			[]html_check.Info{{6, 8, "Invalid protocol!"}},
		},
	}
	for _, test := range TestTable {
		if _, res := html_check.HtmlCheck(test.filename); !reflect.DeepEqual(res, test.result) {
			fmt.Println(res, test.result)
			t.Errorf("Expected: %v\nFound: %v", test.result, res)
		}
	}
}
