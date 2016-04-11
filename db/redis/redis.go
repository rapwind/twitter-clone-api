package redis

import (
	"fmt"
	"strconv"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/techcampman/twitter-d-server/errors"
)

// Store is a wrapper redigo's *redigo.Pool structure
type Store struct {
	pool          *redigo.Pool
	defaultExpire time.Duration
}

// Config ... RedisConfig is a configuration structure
type Config struct {
	Host              string
	Password          string
	MaxActive         int
	MaxIdle           int
	IdleTimeout       time.Duration
	DefaultExpiration time.Duration
}

// NewRedisStore returns a pointer of Store structure
func NewRedisStore(conf Config) *Store {
	return &Store{pool: &redigo.Pool{
		MaxActive:   conf.MaxActive,
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: conf.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			// the redis protocol should probably be made sett-able
			c, err := redigo.Dial("tcp", conf.Host)
			if err != nil {
				return nil, err
			}
			if len(conf.Password) > 0 {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}, defaultExpire: conf.DefaultExpiration}
}

// Pool Returns the connection pool for redis
func (r *Store) Pool() *redigo.Pool {
	return r.pool
}

// Close closes the connection pool for redis
func (r *Store) Close() error {
	return r.pool.Close()
}

// Do is a wrapper around redigo's redigo.Do method that executes any redis command.
// Do does not support prefix support. Example usage: redigo.Do("INCR", "counter").
func (r *Store) Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	conn := r.pool.Get()
	// conn.Close() returns an error but we are already returning regarding error
	// while returning the Do(..) response
	defer conn.Close()
	reply, err = conn.Do(cmd, args...)
	return
}

// Send is a wrapper around redigo's redigo.Send method that writes the command to the client's output buffer.
func (r *Store) Send(cmd string, args ...interface{}) error {
	conn := r.pool.Get()
	// conn.Close() returns an error but we are already returning regarding error
	// while returning the Do(..) response
	defer conn.Close()
	return conn.Send(cmd, args...)
}

// Get is used to get the value of key. If the key does not exist an empty string is returned.
// Usage: redigo.Get("name")
func (r *Store) Get(key string) (reply interface{}, err error) {
	reply, err = r.Do("GET", key)
	if err == nil && reply == nil {
		err = errors.ErrDataNotFound
	}
	return
}

// GetString is used to get the value of key. If the key does not exist an empty string is returned.
// Usage: redigo.Get("John Do")
func (r *Store) GetString(key string) (string, error) {
	return redigo.String(r.Do("GET", key))
}

// GetInt is used the value of key as an integer.
// If the key does not exist or the stored value is a non-integer, zero is returned.
// Example usage: redigo.GetInt("counter")
func (r *Store) GetInt(key string) (int, error) {
	return redigo.Int(r.Do("GET", key))
}

// Delete is used to remove the specified a key.
// Key is ignored if it does not exist.
// It returns the number of keys that were removed.
// Example usage: redigo.Delete("name")
func (r *Store) Delete(key string) (int, error) {
	return redigo.Int(r.Do("DEL", key))
}

// BulkDelete is used to remove the specified keys.
// Key is ignored if it does not exist.
// It returns the number of keys that were removed.
// Example usage: redigo.Delete("name")

// Delete(key)を優先して使って下さい
//func (r *Store) BulkDelete(key ...string) (int, error) {
//	return redigo.Int(r.Do("DEL", key))
//}

// Increment increments the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or
// contains a string that can not be represented as integer
func (r *Store) Increment(key string) (int, error) {
	return redigo.Int(r.Do("INCR", key))
}

// IncrementBy increments the number stored at key by given number.
// If the key does not exist, it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or
// contains a string that can not be represented as integer
func (r *Store) IncrementBy(key string, by int64) (int64, error) {
	return redigo.Int64(r.Do("INCRBY", key, by))
}

// Decrement decrements the number stored at key by one. If the key does not exist,
// it is set to 0 before performing the operation.
// An error is returned if the key contains a value of the wrong type or
// contains a string that can not be represented as integer
func (r *Store) Decrement(key string) (int, error) {
	return redigo.Int(r.Do("DECR", key))
}

// Expire sets a timeout on a key. After the timeout has expired, the key will automatically be deleted.
// Calling Expire on a key that has already expire set will update the expire value.
func (r *Store) Expire(key string, timeout time.Duration) (err error) {
	seconds := strconv.Itoa(int(timeout.Seconds()))
	reply, err := redigo.Int(r.Do("EXPIRE", key, seconds))
	if err != nil {
		return
	}
	if reply != 1 {
		err = fmt.Errorf("key does not exist or the timeout could not be set")
	}

	return
}

// Set key to hold the string value and set key to timeout after a given number of seconds.
// This command is equivalent to executing the following commands:
// SET mykey value
// EXPIRE mykey seconds
// SETEX is atomic, and can be reproduced by using the previous two commands inside an MULTI / EXEC block.
// It is provided as a faster alternative to the given sequence of operations,
// because this operation is very common when Redis is used as a cache.
// An error is returned when seconds is invalid.
func (r *Store) Set(key string, value interface{}, expire time.Duration) (err error) {
	reply, err := redigo.String(r.Do("SETEX", key, strconv.Itoa(int(expire.Seconds())), value))
	if err != nil {
		return
	}
	if reply != "OK" {
		err = fmt.Errorf("reply string is wrong!: %s", reply)
	}

	return
}

// CreatePubSubConn wraps a Conn with convenience methods for subscribers.
func (r *Store) CreatePubSubConn() *redigo.PubSubConn {
	return &redigo.PubSubConn{Conn: r.pool.Get()}
}

// Exists returns true if key exists or false if not.
func (r *Store) Exists(key string) bool {
	// does not have any err message to be checked, it return either 1 or 0
	reply, _ := redigo.Int(r.Do("EXISTS", key))
	return reply == 1
}

// Scard gets the member count of a Set with given key
func (r *Store) Scard(key string) (int, error) {
	return redigo.Int(r.Do("SCARD", key))
}

// SortedSetIncrementBy increments the value of a member in a sorted set.
// This function tries to return last floating value of the item,
// if it fails to parse reply to float64, returns parsing error along with Reply it self
func (r *Store) SortedSetIncrementBy(key string, incrementBy, item interface{}) (float64, error) {
	prefixed := []interface{}{}
	// add key
	prefixed = append(prefixed, key)

	// add incrementBy
	prefixed = append(prefixed, incrementBy)

	// add item
	prefixed = append(prefixed, item)

	return redigo.Float64(r.Do("ZINCRBY", prefixed...))
}

// SortedSetReverseRange returns the specified range of elements in the sorted set stored at key.
// ZREVRANGE key start stop [WITHSCORES]
// The elements are considered to be ordered from the highest to the lowest score.
// Descending lexicographical order is used for elements with equal score.
// Apart from the reversed ordering, ZREVRANGE is similar to ZRANGE.
func (r *Store) SortedSetReverseRange(key string, rest ...interface{}) ([]interface{}, error) {
	// create a slice with rest length +1
	// because we are gonna prepend key to it
	prefixedReq := make([]interface{}, len(rest)+1)

	// prepend prefixed key
	prefixedReq[0] = key

	for i, el := range rest {
		prefixedReq[i+1] = el
	}

	return redigo.Values(r.Do("ZREVRANGE", prefixedReq...))
}

// AddSetMembers adds given elements to the set stored at key.
// Given elements that are already included in set are ignored.
// Returns successfully added key count and error state
func (r *Store) AddSetMembers(key string, rest ...interface{}) (int, error) {
	return redigo.Int(r.Do("SADD", (r.prepareArgsWithKey(key, rest...))...))
}

// RemoveSetMembers removes given elements from the set stored at key
// Returns successfully removed key count and error state
func (r *Store) RemoveSetMembers(key string, rest ...interface{}) (int, error) {
	return redigo.Int(r.Do("SREM", (r.prepareArgsWithKey(key, rest...))...))
}

// GetSetMembers returns all members included in the set at key
// Returns members array and error state
func (r *Store) GetSetMembers(key string) ([]interface{}, error) {
	return redigo.Values(r.Do("SMEMBERS", key))
}

// PopSetMember removes and returns a random element from the set stored at key
func (r *Store) PopSetMember(key string) (string, error) {
	return redigo.String(r.Do("SPOP", key))
}

// IsSetMember checks existence of a member set
func (r *Store) IsSetMember(key string, value string) (int, error) {
	return redigo.Int(r.Do("SISMEMBER", (r.prepareArgsWithKey(key, value))...))
}

// RandomSetMember returns random from set, but not removes unline PopSetMember
func (r *Store) RandomSetMember(key string) (string, error) {
	return redigo.String(r.Do("SRANDMEMBER", key))
}

// SortBy sorts elements stored at key with given weight and order(ASC|DESC)
//
// i.e. Suppose we have elements stored at key as object_1, object_2 and object_3
// and their weight is relatively stored at object_1:weight, object_2:weight, object_3:weight
// When we give sortBy parameter as *:weight,
// it gets all weight values and sorts the objects at given key with specified order.
func (r *Store) SortBy(key, sortBy, order string) ([]interface{}, error) {
	return redigo.Values(r.Do("SORT", key, "by", sortBy, order))
}

// Keys returns all keys with given pattern
// WARNING: Redis Doc says: "Don't use KEYS in your regular application code."
func (r *Store) Keys(key string) ([]interface{}, error) {
	return redigo.Values(r.Do("KEYS", key))
}

// Bool converts the given value to boolean
func (r *Store) Bool(reply interface{}) (bool, error) {
	return redigo.Bool(reply, nil)
}

// Int converts the given value to integer
func (r *Store) Int(reply interface{}) (int, error) {
	return redigo.Int(reply, nil)
}

// String converts the given value to string
func (r *Store) String(reply interface{}) (string, error) {
	return redigo.String(reply, nil)
}

// Int64 converts the given value to 64 bit integer
func (r *Store) Int64(reply interface{}) (int64, error) {
	return redigo.Int64(reply, nil)
}

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err.
// Otherwise, Values converts the reply as follows:
// Reply type      Result
// array           reply, nil
// nil             nil, ErrNil
// other           nil, error
func (r *Store) Values(reply interface{}) ([]interface{}, error) {
	return redigo.Values(reply, nil)
}

// prepareArgsWithKey helper method prepends key to given variadic parameter
func (r *Store) prepareArgsWithKey(key string, rest ...interface{}) []interface{} {
	prefixedReq := make([]interface{}, len(rest)+1)

	// prepend prefixed key
	prefixedReq[0] = key

	for i, el := range rest {
		prefixedReq[i+1] = el
	}

	return prefixedReq
}

// SortedSetsUnion creates a combined set from given list of sorted set keys.
//
// See: http://redigo.io/commands/zunionstore
func (r *Store) SortedSetsUnion(destination string, keys []string, weights []interface{}, aggregate string) (reply int64, err error) {
	if destination == "" {
		err = fmt.Errorf("no destination to store")
		return
	}

	if len(keys) == 0 {
		err = fmt.Errorf("no keys")
		return
	}

	prefixed := []interface{}{destination, len(keys)}

	for _, key := range keys {
		prefixed = append(prefixed, key)
	}

	if len(weights) != 0 {
		prefixed = append(prefixed, "WEIGHTS")
		prefixed = append(prefixed, weights...)
	}

	if aggregate != "" {
		prefixed = append(prefixed, "AGGREGATE", aggregate)
	}

	return redigo.Int64(r.Do("ZUNIONSTORE", prefixed...))
}

// SortedSetScore returns score of a member in a sorted set. If no member,
// an error is returned.
//
// See: http://redigo.io/commands/zscore
func (r *Store) SortedSetScore(key string, member interface{}) (float64, error) {
	return redigo.Float64(r.Do("ZSCORE", key, member))
}

// SortedSetRem removes a member from a sorted set. If no member, an error
// is returned.
//
// See: http://redigo.io/commands/zrem
func (r *Store) SortedSetRem(key string, members ...interface{}) (int64, error) {
	prefixed := []interface{}{key}
	prefixed = append(prefixed, members...)

	return redigo.Int64(r.Do("ZREM", prefixed...))
}
