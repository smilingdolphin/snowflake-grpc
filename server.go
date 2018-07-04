package snowflake

import (
	"context"
	"log"

	pb "github.com/fpay/snowflake-go/pb"

	"github.com/fpay/foundation-go/cache"
)

const (
	SnowflakeSequenceKey = "snowflake:sequence:key"
	SequenceStep         = 1
)

type SnowflakeService struct {
	redis *cache.RedisCache
	node  *Snowflake
}

func NewSnowflakeServer(r *cache.RedisCache) (*SnowflakeService, error) {
	id, err := r.IncrBy(SnowflakeSequenceKey, SequenceStep)
	if err != nil {
		log.Fatalf("Fetch snowflake sequence key err: %s", err)
		return nil, err
	}
	n, err := NewSnowflake(id % workeridMax)
	if err != nil {
		log.Fatalf("Initialize snowflake node err: %s", err)
		return nil, err
	}

	return &SnowflakeService{
		redis: r,
		node:  n,
	}, nil
}

func (sf *SnowflakeService) Generate(context.Context, *pb.Request) (*pb.Response, error) {
	r := new(pb.Response)
	id := sf.node.Generate()
	r.Uniqid = id
	return r, nil
}
