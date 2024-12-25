package utils

import "testing"

func TestReplaceSensitive(t *testing.T) {
	mp := map[string]string{
		"text":        "共产党是中华人民共和国的执政党",
		"replacement": "*",
	}

	//testPost(t, "/fd/replace_sensitive", mp)
}
