package redis

import "errors"

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
	// ErrEmptyKey is an error that indicates no redis key is set.
	ErrEmptyKey = errors.New("empty redis key")
)

type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Name string `json:",optional"`
		Host string
		Type string `json:",default=node,options=node|cluster"`
		Pass string `json:",optional"`
	}

	// A RedisKeyConf is a redis config with key.
	RedisKeyConf struct {
		RedisConf
		Key string `json:",optional"`
	}
)

// NewRedis returns a Redis.
func (rc RedisConf) NewRedis() *Redis {
	// 2021-03-09 by hujiachao
	// 由于redis会区分内网/外网地址，虽然是同一个实例，Host不同会导致一致性hash到不同的位置
	// 这样会影响本地debug线上问题，故增加Name字段用于替代Host计算hash，保证结果计算结果一致
	return NewRedisWithName(rc.Name, rc.Host, rc.Type, rc.Pass)
}

// Validate validates the RedisConf.
func (rc RedisConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

// Validate validates the RedisKeyConf.
func (rkc RedisKeyConf) Validate() error {
	if err := rkc.RedisConf.Validate(); err != nil {
		return err
	}

	if len(rkc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
