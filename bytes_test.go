package govalid

import (
	"encoding/base64"
	"testing"
)

type customByte byte

func TestBytesAcceptsSlicesArraysAndAliases(t *testing.T) {
	values := []any{
		[]byte{1, 2, 3},
		[3]byte{1, 2, 3},
		[]customByte{1, 2, 3},
	}
	for _, value := range values {
		if err := runRule(value, Bytes()); err != nil {
			t.Fatalf("unexpected validation error for %T: %v", value, err)
		}
	}
	if err := runRule([]uint16{1}, Bytes()); err == nil {
		t.Fatal("expected a byte type error")
	}
}

func TestBytesNilEmptyAndLengthRules(t *testing.T) {
	var nilBytes []byte
	if err := runRule(nilBytes, BytesNil()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule([0]byte{}, BytesNotNil()); err != nil {
		t.Fatalf("arrays should be non-nil bytes: %v", err)
	}

	value := []byte{1, 2, 3}
	rules := []Rule{
		BytesNotEmpty(), BytesLength(3), BytesMinLength(2), BytesMaxLength(4),
		BytesLengthBetween(2, 4), BytesLengthNotBetween(4, 6),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
	if err := runRule(value, BytesLength(-1)); err == nil {
		t.Fatal("expected a negative length error")
	}
}

func TestBytesComparisonRules(t *testing.T) {
	value := []byte("govalid")
	rules := []Rule{
		BytesEqual([]byte("govalid")), BytesNotEqual([]byte("invalid")),
		BytesConstantTimeEqual([]byte("govalid")),
		BytesOneOf([]byte("other"), []byte("govalid")),
		BytesNoneOf([]byte("other"), []byte("invalid")),
		BytesContains([]byte("valid")), BytesNotContains([]byte("invalid")),
		BytesContainsAny([]byte("missing"), []byte("go")),
		BytesContainsAll([]byte("go"), []byte("valid")),
		BytesHasPrefix([]byte("go")), BytesNotHasPrefix([]byte("no")),
		BytesHasSuffix([]byte("valid")), BytesNotHasSuffix([]byte("go")),
		BytesGreaterThan([]byte("alpha")), BytesLessThan([]byte("zulu")),
		BytesBetween([]byte("alpha"), []byte("zulu")),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
}

func TestBytesRulesCloneExpectedArguments(t *testing.T) {
	expected := []byte("safe")
	rule := BytesEqual(expected)
	expected[0] = 'x'
	if err := runRule([]byte("safe"), rule); err != nil {
		t.Fatalf("rule should retain its original expected value: %v", err)
	}
}

func TestBytesContentRules(t *testing.T) {
	rules := []Rule{
		BytesNoZero(), BytesNotAllZero(), BytesUnique(), BytesASCII(), BytesPrintableASCII(),
	}
	for _, rule := range rules {
		if err := runRule([]byte("Go!"), rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
	if err := runRule([]byte{0, 0}, BytesAllZero()); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if err := runRule([]byte{1, 1}, BytesUnique()); err == nil {
		t.Fatal("expected duplicate byte error")
	}
}

func TestBytesFormatRules(t *testing.T) {
	encoded := []byte(base64.StdEncoding.EncodeToString([]byte("govalid")))
	tests := []struct {
		value []byte
		rule  Rule
	}{
		{[]byte("olá"), BytesUTF8()},
		{[]byte{0xff}, BytesNotUTF8()},
		{[]byte(`{"valid":true}`), BytesJSON()},
		{[]byte(`<root><valid>true</valid></root>`), BytesXML()},
		{[]byte("676f76616c6964"), BytesHex()},
		{encoded, BytesBase64()},
	}
	for _, test := range tests {
		if err := runRule(test.value, test.rule); err != nil {
			t.Fatalf("unexpected validation error for %q: %v", test.value, err)
		}
	}
	if err := runRule([]byte(`{"broken"`), BytesJSON()); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestBytesNestedRules(t *testing.T) {
	value := []byte{2, 4, 6}
	rules := []Rule{
		BytesEach(IntEven()),
		BytesAt(1, IntEqual(uint8(4))),
		BytesAtIfPresent(10, IntPositive()),
	}
	for _, rule := range rules {
		if err := runRule(value, rule); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	}
	if err := runRule(value, BytesAt(-1, IntPositive())); err == nil {
		t.Fatal("expected a negative index error")
	}
}
