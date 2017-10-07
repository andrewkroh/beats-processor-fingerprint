# beats-processor-fingerprint plugin

[![Build Status](http://img.shields.io/travis/andrewkroh/beats-processor-fingerprint.svg?style=flat-square)][travis]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[travis]: http://travis-ci.org/andrewkroh/beats-processor-fingerprint
[godocs]: http://godoc.org/github.com/andrewkroh/beats-processor-fingerprint
[releases]: https://github.com/andrewkroh/beats-processor-fingerprint/releases

beats-processor-fingerprint is a processor plugin for Elastic Beats that "fingerprints"
events by computing a hash of fields in the event. This can be useful as a checksum of the
data or it can be useful for deduplication purposes if the hash is used as a primary
key in your data store.

## Installation and Usage

Build the plugin or download a [release][releases]. Go plugins are only supported on
Linux at the current time. They must be compiled with the same Go version as the Beat.

```
go build -buildmode=plugin
```

Start a Beat with the plugin.

```
filebeat -e --plugin ./processor-fingerprint-linux-amd64.so
```

Add the processor to your configuration file.

```
processors:
- fingerprint:
    hash: sha256
    encoding: hex
    target: fingerprint
    fields: [source, timestamp]
```

## Configuration Options

- **`hash`**: The hashing algoritm to use. The supported algoritms are `md5`, `sha1`,
  `sha224`, `sha256` (default), `sha384`, `sha512`, `sha512_224`, `sha512_256`,
  `sha3_224`, `sha3_256`, `sha3_384`, and `sha3_512`.
- **`encoding`**: Encoding type to use on the hash value. The supported encoding types
  are `hex` (default), `base32`, and `base64`.
- **`target`**: The name of the field to which the encoded hash value will be written.
  The default value is `fingerprint`.
- **`fields`**: A list of field names whose values will be concatenated and hashed.
  Missing fields will be ignored. The default value is `[message]`.

## Example Output

```json
{
  "@timestamp": "2017-10-07T03:09:50.201Z",
  "@metadata": {
    "beat": "filebeat",
    "type": "doc",
    "version": "7.0.0-alpha1"
  },
  "source": "/home/andrew_kroh/go/src/github.com/elastic/beats/filebeat/messages",
  "offset": 68379,
  "message": "Oct  6 23:17:59 localhost systemd: Starting Session 6 of user andrew_kroh.",
  "beat": {
    "name": "host.example.com",
    "hostname": "host.example.com",
    "version": "7.0.0-alpha1"
  },
  "fingerprint": "42c4f7ffd2c28adbba04abf6c3bf28b3de9a8afea2227fcba8ac73d595a4209e"
}
```

