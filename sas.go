package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// NewSharedAccessKey creates new shared access key for subsequent token generation.
func NewSharedAccessKey(hostname, policy, key string) *SharedAccessKey {
	return &SharedAccessKey{
		HostName:            hostname,
		SharedAccessKeyName: policy,
		SharedAccessKey:     key,
	}
}

// SharedAccessKey is SAS token generator.
type SharedAccessKey struct {
	HostName            string
	SharedAccessKeyName string
	SharedAccessKey     string
}

// Token generates a shared access signature for the named resource and lifetime.
func (c *SharedAccessKey) Token(resource string, lifetime time.Duration) (*SharedAccessSignature, error) {
	return NewSharedAccessSignature(
		resource, c.SharedAccessKeyName, c.SharedAccessKey, time.Now().Add(lifetime),
	)
}

// NewSharedAccessSignature initialized a new shared access signature
// and generates signature fields based on the given input.
func NewSharedAccessSignature(resource, policy, key string, expiry time.Time) (*SharedAccessSignature, error) {
	sig, err := mksig(resource, key, expiry)
	if err != nil {
		return nil, err
	}
	return &SharedAccessSignature{
		Sr:  resource,
		Sig: sig,
		Se:  expiry,
		Skn: policy,
	}, nil
}

func mksig(sr, key string, se time.Time) (string, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, b)
	if _, err := fmt.Fprintf(h, "%s\n%d", url.QueryEscape(sr), se.Unix()); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// SharedAccessSignature is a shared access signature instance.
type SharedAccessSignature struct {
	Sr  string
	Sig string
	Se  time.Time
	Skn string
}

// String converts the signature to a token string.
func (sas *SharedAccessSignature) String() string {
	s := "SharedAccessSignature " +
		"sr=" + url.QueryEscape(sas.Sr) +
		"&sig=" + url.QueryEscape(sas.Sig) +
		"&se=" + url.QueryEscape(strconv.FormatInt(sas.Se.Unix(), 10))
	if sas.Skn != "" {
		s += "&skn=" + url.QueryEscape(sas.Skn)
	}
	return s
}
