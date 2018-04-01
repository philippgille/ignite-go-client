package ignite

import (
	"fmt"
	"math/rand"
)

// CacheReplaceIfEquals puts a value with a given key to cache only if
// the key already exists and value equals provided value.
func (c *client) CacheReplaceIfEquals(cache string, binary bool, key interface{}, valueCompare interface{}, valueNew interface{}, status *int32) (bool, error) {
	if status != nil {
		*status = StatusSuccess
	}

	uid := rand.Int63()

	o := c.Prepare(opCacheReplaceIfEquals, uid)
	// prepare data
	if err := o.WritePrimitives(hashCode(cache), binary); err != nil {
		return false, fmt.Errorf("failed to write cache id and binary flag: %s", err.Error())
	}
	if err := o.WriteObjects(key, valueCompare, valueNew); err != nil {
		return false, fmt.Errorf("failed to write cache key and values: %s", err.Error())
	}

	// execute
	r, err := c.Call(o)
	if err != nil {
		return false, fmt.Errorf("failed to execute operation: %s", err.Error())
	}
	if r.UID != uid {
		return false, fmt.Errorf("invalid response id (expected %d, but received %d)", uid, r.UID)
	}
	if status != nil {
		*status = r.Status
	}
	if r.Status != StatusSuccess {
		return false, fmt.Errorf("failed to execute operation: status=%d, message=%s", r.Status, r.Message)
	}

	// read response data
	var res bool
	if err := r.ReadPrimitives(&res); err != nil {
		return false, fmt.Errorf("failed to read result value: %s", err.Error())
	}

	return res, nil
}
