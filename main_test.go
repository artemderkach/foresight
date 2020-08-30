package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	// kill current process after 1 second
	go func() {
		time.Sleep(3 * time.Second)
		e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.Nil(t, e)
	}()

	go func() {
		main()
	}()

	// give service .5 seconds to start
	time.Sleep(500 * time.Millisecond)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/printMeAt", nil)
	require.Nil(t, err)

	// now plus few seconds
	now := time.Now()
	tm := now.Add(time.Duration(2) * time.Second)

	q := req.URL.Query()
	q.Add("time", strconv.FormatInt(tm.Unix(), 10))
	q.Add("message", "this message should be printed")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	require.Nil(t, err)

	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, string(body), "done")

	// give time to print msg

	// uncomment to view print message
	// time.Sleep(5 * time.Second)
	// require.Nil(t, "")
}
