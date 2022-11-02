package cache

import (
	"encoding/base64"
	"encoding/json"
	"gotag/pkg/l"
	"gotag/pkg/v"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Redis struct {
	cli  *redis.Client
	host string
}

func InitRedisFromViper() (Cache, error) {
	host := v.GetViper().GetString("cache.host")
	password := v.GetViper().GetString("cache.password")
	cli := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
	})

	_, err := cli.Ping().Result()
	if err != nil {
		l.GetLogger().Error("Connect to redis host"+host+"failed ", zap.Error(err))
		return nil, err
	}

	l.GetLogger().Info("connect redis success", zap.String("host", host))
	return &Redis{
		cli:  cli,
		host: host,
	}, nil
}

func (c *Redis) Set(key string, value interface{}, timeout int) error {
	if c.cli == nil {
		err := c.reconnect()
		if err != nil {
			return err
		}
	}

	timeouts := time.Duration(timeout) * time.Second
	cmd := c.cli.Set(key, value, timeouts)
	if cmd.Err() != nil {
		l.GetLogger().Error("Set failed, key:"+key, zap.Error(cmd.Err()))
		return cmd.Err()
	} else {
		return nil
	}
}

func (c *Redis) Get(key string, to interface{}) error {
	if c.cli == nil {
		err := c.reconnect()
		if err != nil {
			return err
		}
	}

	exist, err := c.Exist(key)
	if err != nil {
		l.GetLogger().Error("Get Exist failed", zap.String("key", key))
		return errors.New("Cache get failed:" + key)
	}

	if !exist {
		l.GetLogger().Info("Get key do not existed", zap.String("key", key))
		return errors.New("Cache key do not existed")
	}

	cmd := c.cli.Get(key)
	if cmd == nil || cmd.Err() != nil {
		l.GetLogger().Error("Get failed", zap.String("key", key))
		return errors.New("Cache get failed:" + key)
	}

	err = json.Unmarshal([]byte(cmd.Val()), to)
	if err != nil {
		return err
	}

	return nil
}

func (c *Redis) Del(key string) error {
	if c.cli == nil {
		err := c.reconnect()
		if err != nil {
			return err
		}
	}

	cmd := c.cli.Del(key)
	if cmd == nil || cmd.Err() != nil {
		l.GetLogger().Error("Del failed", zap.String("key", key), zap.Error(cmd.Err()))
	}

	return cmd.Err()
}

func (c *Redis) Exist(key string) (bool, error) {
	if c.cli == nil {
		err := c.reconnect()
		if err != nil {
			return false, err
		}
	}

	cmd := c.cli.Exists(key)
	if cmd.Err() != nil {
		l.GetLogger().Error("Exist failed", zap.String("key", key), zap.Error(cmd.Err()))
		return false, cmd.Err()
	}

	return cmd.Val() > 0, nil
}

func (c *Redis) SetExpireTime(key string, timeout int) error {
	if c.cli == nil {
		err := c.reconnect()
		if err != nil {
			return err
		}
	}

	cmd := c.cli.Expire(key, time.Duration(timeout)*time.Second)
	if cmd == nil || cmd.Err() != nil {
		l.GetLogger().Error("Exist failed", zap.String("key", key), zap.Error(cmd.Err()))
	}

	return nil
}

func (c *Redis) Close() {
	c.cli.Close()
}

// SetNx Set if not exists
func (c *Redis) SetNx(key string, value interface{}, timeout int) (bool, error) {
	data, err := c.encode(value)
	if err != nil {
		return false, err
	}

	if c.cli == nil {
		err = c.reconnect()
		if err != nil {
			return false, err
		}
	}

	cmd := c.cli.SetNX(key, data, time.Duration(timeout)*time.Second)
	if cmd.Err() != nil {
		l.GetLogger().Error("SetNX failed, key:"+key, zap.Error(cmd.Err()))
		return false, cmd.Err()
	}
	return cmd.Val(), nil
}

func (c *Redis) reconnect() error {
	c.cli = redis.NewClient(&redis.Options{
		Addr: c.host,
	})

	_, err := c.cli.Ping().Result()

	return err
}

//gob encoding
func (c *Redis) encode(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)

	return encoded, nil
}
