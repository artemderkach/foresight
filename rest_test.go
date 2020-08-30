package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequiredParams(t *testing.T) {
	tt := []struct {
		Key1         string
		Value1       string
		Key2         string
		Value2       string
		ResponseBody string
	}{
		{
			"time",
			"123",
			"message",
			"hello",
			"done",
		},
		{
			"",
			"",
			"message",
			"hello",
			"time parameter required",
		},
		{
			"time",
			"123",
			"",
			"",
			"message parameter required",
		},
	}

	rest := &Rest{}
	ts := httptest.NewServer(rest.Router())

	for _, test := range tt {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/printMeAt", nil)
		require.Nil(t, err)

		q := req.URL.Query()
		q.Add(test.Key1, test.Value1)
		q.Add(test.Key2, test.Value2)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		require.Nil(t, err)

		body, err := ioutil.ReadAll(res.Body)
		require.Nil(t, err)

		require.Equal(t, string(body), test.ResponseBody)
	}
}
