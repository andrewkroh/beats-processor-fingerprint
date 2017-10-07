package main

import (
	"bytes"
	"fmt"
	"hash"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

type Fingerprint struct {
	config FingerprintConfig
	hash   hash.Hash // Hash function.
}

func (f Fingerprint) String() string {
	b := new(bytes.Buffer)
	b.WriteString("fingerprint=[")

	b.WriteString("fields=")
	b.WriteString(strings.Join(f.config.Fields, ", "))

	b.WriteString(", hash=")
	b.WriteString(f.config.Hash.name)

	b.WriteString(", target=")
	b.WriteString(f.config.Target)

	b.WriteRune(']')

	return b.String()
}

func New(c *common.Config) (processors.Processor, error) {
	fc := defaultFingerprintConfig
	err := c.Unpack(&fc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack fingerprint config")
	}

	return &Fingerprint{
		config: fc,
		hash:   fc.Hash.producer(),
	}, nil
}

func (f *Fingerprint) Run(event *beat.Event) (*beat.Event, error) {
	f.hash.Reset()
	for _, field := range f.config.Fields {
		v, _ := event.GetValue(field)
		writeValue(f.hash, v)
	}

	// Compute the hash and encode the value.
	event.PutValue(f.config.Target, f.config.Encoder(f.hash.Sum(nil)))
	return event, nil
}

func writeValue(writer io.Writer, object interface{}) {
	val := reflect.ValueOf(object)

	// Follow the pointer.
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = reflect.ValueOf(val.Elem().Interface())
	}

	if val.IsValid() {
		// Ensure we consistently hash times in UTC.
		o := val.Interface()
		if ts, ok := o.(time.Time); ok {
			o = ts.UTC()
		}

		writer.Write([]byte(fmt.Sprintf("%v", val.Interface())))
	}
}
