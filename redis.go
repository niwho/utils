package utils

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/niwho/logs"
)

type RedisClient struct {
	pool   *redis.Pool
	server string
	db     int
	redis.Conn
}

func NewRedisClient(address string, db int) *RedisClient {
	rc := &RedisClient{
		server: address,
		db:     db,
	}
	rc.pool = &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", rc.server, 280*time.Millisecond, 200*time.Millisecond, 200*time.Millisecond)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: time.Second * 30,
	}
	return rc
}

func (rc *RedisClient) GetString(key string) (string, error) {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()

	return redis.String(conn.Do("GET", key))

}

func (rc *RedisClient) SetString(key, val string, ex int) error {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()

	_, err := conn.Do("SET", key, val, "EX", ex)
	return err

}

func (rc *RedisClient) ZScore(key string, member string) (string, error) {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()
	return redis.String(conn.Do("ZSCORE", key, member))

}

func (rc *RedisClient) Remove(key string) {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()
	conn.Do("DEL", key)

}

func (rc *RedisClient) ZAdd(key string, members map[string]int64, ttl int) {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()

	batch_num := 100
	var k int = 0
	for mem, val := range members {
		conn.Send("ZADD", key, val, mem)
		k += 1
		if k > batch_num {
			conn.Flush()
			conn.Receive() // reply from SET
			// _, err := conn.Do("EXEC")
			_, err := conn.Receive()
			logs.Log(logs.F{"err": err}).Error("redis zadd")
			k = 0
		}
	}
	if k > 0 {
		conn.Flush()
		conn.Receive() // reply from SET
		// _, err := conn.Do("EXEC")
		_, err := conn.Receive()
		logs.Log(logs.F{"err": err}).Error("redis zadd")

	}
	conn.Do("EXPIRE", key, ttl)
}

func (rc *RedisClient) getConn() redis.Conn {
	return rc.pool.Get()
}

func (rc *RedisClient) GetInt(key string) (int, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return redis.Int(conn.Do("GET", key))
}

func (rc *RedisClient) SetInt(key string, val int, ex time.Duration) error {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()

	_, err := conn.Do("SET", key, val, "EX", int64(ex.Seconds()))
	return err
}

func (rc *RedisClient) GetInt64(key string) (int64, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return redis.Int64(conn.Do("GET", key))
}

func (rc *RedisClient) SetInt64(key string, val int64, ex time.Duration) error {
	conn := rc.pool.Get()
	defer func() {
		conn.Close()
	}()

	_, err := conn.Do("SET", key, val, "EX", int64(ex.Seconds()))
	return err
}

func (rc *RedisClient) ZAddOne(key string, score int64, val string, limit int, ex time.Duration) error {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	_, err := redis.Int(conn.Do("ZADD", key, score, val))
	if err != nil {
		return err
	}
	_, err = redis.Int(conn.Do("EXPIRE", key, int64(ex.Seconds())))

	count, _ := redis.Int(conn.Do("ZCARD", key))
	if count > limit {
		_, err = redis.Int(conn.Do("ZREMRANGEBYRANK", key, 0, 0))
	}
	if err != nil {
		logs.Log(logs.F{"key": key, "count": count, "limit": limit, "err": err}).Error("ZAddOne")
	}
	return err
}

func (rc *RedisClient) ZRange(key string, start int, stop int) ([]string, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	reply, err := conn.Do("ZRANGE", key, start, stop)
	fmt.Println("zzzz", reply, err)
	return redis.Strings(reply, err)
}

func (rc *RedisClient) ZRangeWithScore(key string, start int, stop int) (map[string]int64, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return redis.Int64Map(conn.Do("ZRANGE", key, start, stop, "WITHSCORES"))
}

func (rc *RedisClient) ZRevRange(key string, start int, stop int) ([]string, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return redis.Strings(conn.Do("ZREVRANGE", key, start, stop))
}

func (rc *RedisClient) ZRemRangeByRank(key string, start int, stop int) (int, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return redis.Int(conn.Do("ZREMRANGEBYRANK", key, start, stop))
}

func (rc *RedisClient) ZRem(key, member string) (int, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return redis.Int(conn.Do("ZREM", key, member))
}

func (rc *RedisClient) ZAddAndTrim(key string, members map[string]int64, start int, stop int, ex time.Duration) error {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	args := []interface{}{key}
	for val, score := range members {
		args = append(args, val, score)
	}

	_, err := redis.Int(conn.Do("ZAdd", args...))
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("expire", key, int64(ex.Seconds())))
	if err != nil {
		return err
	}

	_, err = redis.Int(conn.Do("ZRemRangeByRank", start, stop))
	return err
}

func (rc *RedisClient) HGet(key, member string) (interface{}, error) {
	conn := rc.getConn()
	defer func() {
		conn.Close()
	}()

	return conn.Do("HGET", key, member)
}
func (rc *RedisClient) HGetInt64(key, member string) (int64, error) {
	return redis.Int64(rc.HGet(key, member))
}

func (rc *RedisClient) HSet(key, member string, val interface{}) error {
	conn := rc.getConn()
	defer conn.Close()

	_, err := redis.Int(conn.Do("HSET", key, member, val))
	return err
}

func (rc *RedisClient) Del(key string) error {
	conn := rc.getConn()
	defer conn.Close()

	_, err := conn.Do("del", key)
	return err
}
