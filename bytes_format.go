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

func BytesUTF8() Rule {
	return bytesSlicePredicateRule(utf8.Valid, "bytes should contain valid UTF-8")
}

func BytesNotUTF8() Rule {
	return bytesSlicePredicateRule(
		func(value []byte) bool { return !utf8.Valid(value) },
		"bytes should not contain valid UTF-8")
}

func BytesJSON() Rule {
	return bytesSlicePredicateRule(json.Valid, "bytes should contain valid JSON")
}

func BytesXML() Rule {
	return bytesSliceRule(func(value []byte) error {
		var destination struct{}
		if err := xml.Unmarshal(value, &destination); err != nil {
			return errors.New("bytes should contain valid XML")
		}
		return nil
	})
}

func BytesPEM() Rule {
	return bytesSlicePredicateRule(func(value []byte) bool {
		block, rest := pem.Decode(value)
		return block != nil && len(bytes.TrimSpace(rest)) == 0
	}, "bytes should contain a valid PEM block")
}

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

func BytesBase64() Rule {
	return encodedBytesRule(base64.StdEncoding, "bytes should contain valid base64 text")
}

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
