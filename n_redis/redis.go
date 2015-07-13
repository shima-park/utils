package n_redis

import (
	"git.wdwd.com/nova/n_utils"
	"github.com/garyburd/redigo/redis"
	"time"
)

var Pool *redis.Pool

func RedisInit(server string, pwd string, params map[string]string) {
	var max_idle, max_active int
	var idle_timeout time.Duration

	if params["maxIdle"] != "" {
		max_idle = int(n_utils.Be_int(params["maxIdle"]))
	} else {
		max_idle = 3
	}

	// 未设置则没限制
	if params["maxActive"] != "" {
		max_active = int(n_utils.Be_int(params["maxActive"]))
	} else {
		max_active = 0
	}

	if params["idleTimeout"] != "" {
		idle_timeout = time.Duration(n_utils.Be_int(params["idleTimeout"]))
	} else {
		idle_timeout = 240
	}

	Pool = &redis.Pool{
		MaxIdle:     max_idle,
		MaxActive:   max_active,
		IdleTimeout: idle_timeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
