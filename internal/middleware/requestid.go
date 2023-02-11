/*
   GoToSocial
   Copyright (C) 2021-2023 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package middleware

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"sync"

	"github.com/gin-gonic/gin"
)

type ridCtxType string

const (
	// RequestIDKey is a string to use as a map key, for example a logger field
	RequestIDKey            = "requestID"
	ridCtxKey    ridCtxType = RequestIDKey
)

var (
	ridLock sync.Mutex
	ridInit sync.Once
	ridSrc  *rand.Rand
)

func ridInitRandom() {
	ridInit.Do(func() {
		var rngSeed int64
		binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
		ridSrc = rand.New(rand.NewSource(rngSeed)) // nolint:gosec
	})
}

func ridGen() string {
	ridInitRandom()
	var id [16]byte

	ridLock.Lock()
	ridSrc.Read(id[:])
	ridLock.Unlock()

	// Use RawURLEncoding because we don't need the padding, but it's
	// possible the rid may be used in a URL at some point in other
	// systems
	return base64.RawURLEncoding.EncodeToString(id[:])
}

// RequestID returns a gin middleware which adds a unique ID for each request
// to the context. It currently directly wraps the upstream library but this
// makes it easier to set any custom options later on.
func RequestID(header string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get id from request
		rid := c.GetHeader(header)
		if rid == "" {
			rid = ridGen()
			c.Request.Header.Set(header, rid)
		}

		ctx := context.WithValue(c.Request.Context(), ridCtxKey, rid)
		c.Request = c.Request.WithContext(ctx)

		// Set the id to ensure that the requestid is in the response
		c.Header(header, rid)
		c.Next()
	}
}

func RequestIDFromCtx(c context.Context) string {
	v, ok := c.Value(ridCtxKey).(string)
	if !ok {
		return ""
	}
	return v
}
