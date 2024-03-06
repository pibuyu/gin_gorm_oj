package test

import (
	"gin_gorm_o/models"
	"testing"
	"time"
)
import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func TestGoRedis(t *testing.T) {
	rdb := models.Redis

	err := rdb.Set(ctx, "age", "24", time.Second*10).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "age").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("age:", val)

	val2, err := rdb.Get(ctx, "zhangchaoyu").Result()
	if err == redis.Nil {
		fmt.Println("zhangchaoyu does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("zhangchaoyu", val2)
	}
	// Output: key value
	// key2 does not exist
}
