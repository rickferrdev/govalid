package govalid

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"unicode/utf8"
)

// BytesUTF8 returns a byte-sequence validation rule for utf8.
func BytesUTF8() Rule {
	return bytesSlicePredicateRule(utf8.Valid, "bytes should contain valid UTF-8")
}

// BytesNotUTF8 returns a byte-sequence validation rule for not utf8.
func BytesNotUTF8() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !utf8.Valid(value) },
		"bytes should not contain valid UTF-8")
}

// BytesJSON returns a byte-sequence validation rule for json.
func BytesJSON() Rule {
	return bytesSlicePredicateRule(json.Valid, "bytes should contain valid JSON")
}

// BytesXML returns a byte-sequence validation rule for xml.
func BytesXML() Rule {
	return bytesSliceRule(func(value []byte) error {
		var destination struct{}
		if err := xml.Unmarshal(value, &destination); err != nil {
			return errors.New("bytes should contain valid XML")
		}
		return nil
	})
}

// BytesPEM returns a byte-sequence validation rule for pem.
func BytesPEM() Rule {
	return bytesSlicePredicateRule(func(value []byte) bool {
		block, rest := pem.Decode(value)
		return block != nil && len(bytes.TrimSpace(rest)) == 0
	}, "bytes should contain a valid PEM block")
}

// BytesHex returns a byte-sequence validation rule for hex.
func BytesHex() Rule {
	return bytesSliceRule(func(value []byte) error {
		if len(value)%2 != 0 {
			return errors.New("bytes should contain valid hexadecimal text")
		}
		decoded := make([]byte, hex.DecodedLen(len(value)))
		if _, err := hex.Decode(decoded, value); err != nil {
			return errors.New("bytes should contain valid hexadecimal text")
		}
		return nil
	})
}

// BytesBase64 returns a byte-sequence validation rule for base64.
func BytesBase64() Rule {
	return encodedBytesRule(base64.StdEncoding, "bytes should contain valid base64 text")
}

// BytesBase64URL returns a byte-sequence validation rule for base64url.
func BytesBase64URL() Rule {
	return encodedBytesRule(base64.URLEncoding, "bytes should contain valid base64url text")
}

func encodedBytesRule(encoding *base64.Encoding, message string) Rule {
	return bytesSliceRule(func(value []byte) error {
		decoded := make([]byte, encoding.DecodedLen(len(value)))
		if _, err := encoding.Decode(decoded, value); err != nil {
			return errors.New(message)
		}
		return nil
	})
}
