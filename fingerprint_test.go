package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)

const message = "[Mon Mar 8 05:31:47 2004] [info] [client 64.242.88.10] " +
	"(104)Connection reset by peer: client stopped connection before send " +
	"body completed"

func newTestFingerprint(t testing.TB, encoding, hash string, fields ...string) *Fingerprint {
	if len(fields) == 0 {
		fields = defaultFingerprintConfig.Fields
	}
	c, err := common.NewConfigFrom(map[string]interface{}{
		"encoding": encoding,
		"hash":     hash,
		"fields":   fields,
	})
	if err != nil {
		t.Fatal(err)
	}

	f, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	return f.(*Fingerprint)
}

func TestFingerprintHashes(t *testing.T) {
	var tests = []struct {
		encoding    string
		hash        string
		fingerprint string
	}{
		{"hex", "md5", "ffc2c1d636ac17a38860df03350be0e4"},
		{"hex", "sha1", "fe7b2aede2119f5508f466209a26a863d405c1ee"},
		{"hex", "sha256", "5c2736f2a1b8ec165ffe3e904b4171ad7581db825f492bca7bdfa0cca4e5630f"},
		{"hex", "sha512", "3142c7ff002141dba401d5cae154f4c534d5e3b8403699398c3e1dd20bcad764accccfab108188f14b01e11072c3a705d8d355473ebbea576f008920d6953b4e"},
	}

	for _, testcase := range tests {
		f := newTestFingerprint(t, testcase.encoding, testcase.hash)
		event := &beat.Event{Fields: common.MapStr{"message": message}}
		event, err := f.Run(event)
		if assert.NoError(t, err) {
			assert.Equal(t, testcase.fingerprint, event.Fields["fingerprint"], "hash: %v", f.config.Hash.name)
		}
	}
}

func TestFingerprintFieldConcat(t *testing.T) {
	f := newTestFingerprint(t, "hex", "sha1",
		"@timestamp", "record_number", "beat.host", "message")
	event := &beat.Event{
		Timestamp: time.Unix(1091067890, 0).UTC(),
		Fields: common.MapStr{
			"record_number": 1888399992,
			"beat": common.MapStr{
				"host": "example",
			},
			"message": message,
		},
	}

	event, err := f.Run(event)
	if assert.NoError(t, err) {
		assert.Equal(t,
			"ee89e405f814a308440c20f11adcc36e9e51c393",
			event.Fields["fingerprint"])
	}
}

func TestFingerprintMissingField(t *testing.T) {
	f := newTestFingerprint(t, "hex", "md5", "other")
	event := &beat.Event{Fields: common.MapStr{"message": message}}

	event, err := f.Run(event)
	if assert.NoError(t, err) {
		assert.Equal(t,
			"d41d8cd98f00b204e9800998ecf8427e",
			event.Fields["fingerprint"])
	}
}

func TestFingerprintString(t *testing.T) {
	f := newTestFingerprint(t, "hex", "sha1")
	assert.Equal(t, "fingerprint=[fields=message, hash=sha1, target=fingerprint]", f.String())
}

func TestWriteValue(t *testing.T) {
	var tests = []struct {
		in  interface{}
		out string
	}{
		{nil, ""},
		{true, "true"},
		{8, "8"},
		{uint(10), "10"},
		{18.123, "18.123"},
		{"hello", "hello"},
	}

	b := new(bytes.Buffer)
	for _, testcase := range tests {
		b.Reset()
		writeValue(b, testcase.in)
		assert.Equal(t, testcase.out, b.String())

		b.Reset()
		writeValue(b, &testcase.in)
		assert.Equal(t, testcase.out, b.String())
	}
}

func BenchmarkFingerprintFilterSHA1(b *testing.B) {
	f := newTestFingerprint(b, "hex", "sha1")
	event := &beat.Event{Fields: common.MapStr{"message": message}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Run(event)
	}
}

func BenchmarkFingerprintFilterSHA256(b *testing.B) {
	f := newTestFingerprint(b, "hex", "sha256")
	event := &beat.Event{Fields: common.MapStr{"message": message}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Run(event)
	}
}

func BenchmarkFingerprintFilterSHA512(b *testing.B) {
	f := newTestFingerprint(b, "hex", "sha512")
	event := &beat.Event{Fields: common.MapStr{"message": message}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Run(event)
	}
}

func BenchmarkFingerprintFilterMD5(b *testing.B) {
	f := newTestFingerprint(b, "hex", "md5")
	event := &beat.Event{Fields: common.MapStr{"message": message}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Run(event)
	}
}
