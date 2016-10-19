// Copyright (c) 2016- Daniel Oaks <daniel@danieloaks.net>
// released under the MIT license

package irc

import "fmt"

// ISupportList holds a list of ISUPPORT tokens
type ISupportList struct {
	Tokens      map[string]*string
	CachedReply [][]string
}

// NewISupportList returns a new ISupportList
func NewISupportList() *ISupportList {
	var il ISupportList
	il.Tokens = make(map[string]*string)
	il.CachedReply = make([][]string, 0)
	return &il
}

// Add adds an RPL_ISUPPORT token to our internal list
func (il *ISupportList) Add(name string, value string) {
	il.Tokens[name] = &value
}

// AddNoValue adds an RPL_ISUPPORT token that does not have a value
func (il *ISupportList) AddNoValue(name string) {
	il.Tokens[name] = nil
}

// RegenerateCachedReply regenerates the cached RPL_ISUPPORT reply
func (il *ISupportList) RegenerateCachedReply() {
	il.CachedReply = make([][]string, 0)
	maxlen := 400      // Max length of a single ISUPPORT token line
	var length int     // Length of the current cache
	var cache []string // Token list cache

	for name, value := range il.Tokens {
		var token string
		if value == nil {
			token = name
		} else {
			token = fmt.Sprintf("%s=%s", name, *value)
		}

		if len(token)+length <= maxlen {
			// account for the space separating tokens
			if len(cache) > 0 {
				length++
			}
			cache = append(cache, token)
			length += len(token)
		}

		if len(cache) == 13 || len(token)+length >= maxlen {
			cache = append(cache, "are supported by this server")
			il.CachedReply = append(il.CachedReply, cache)
			cache = make([]string, 0)
			length = 0
		}
	}

	if len(cache) > 0 {
		cache = append(cache, "are supported by this server")
		il.CachedReply = append(il.CachedReply, cache)
	}
}

// RplISupport outputs our ISUPPORT lines to the client. This is used on connection and in VERSION responses.
func (client *Client) RplISupport() {
	for _, tokenline := range client.server.isupport.CachedReply {
		// ugly trickery ahead
		client.Send(nil, client.server.name, RPL_ISUPPORT, append([]string{client.nick}, tokenline...)...)
	}
}