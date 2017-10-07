package main

import (
	"github.com/elastic/beats/libbeat/plugin"
	"github.com/elastic/beats/libbeat/processors"
)

var Bundle = plugin.Bundle(
	processors.Plugin("fingerprint", New),
)
