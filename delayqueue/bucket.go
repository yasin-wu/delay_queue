package delayqueue

import (
	"strconv"
)

// BucketItem bucket中的元素
type BucketItem struct {
	timestamp int64
	jobId     string
}

func (this *BucketItem) pushToBucket(key string, timestamp int64, jobId string) error {
	_, err := execRedisCommand("ZADD", key, timestamp, jobId)
	return err
}

func (this *BucketItem) getFromBucket(key string) (*BucketItem, error) {
	value, err := execRedisCommand("ZRANGE", key, 0, Setting.ZRangeLimit, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, ErrRecordNotFound
	}

	var valueBytes []interface{}
	valueBytes = value.([]interface{})
	if len(valueBytes) == 0 {
		return nil, ErrRecordNotFound
	}
	timestampStr := string(valueBytes[1].([]byte))
	item := &BucketItem{}
	item.timestamp, _ = strconv.ParseInt(timestampStr, 10, 64)
	item.jobId = string(valueBytes[0].([]byte))
	return item, nil
}

func (this *BucketItem) removeFromBucket(bucket string, jobId string) error {
	_, err := execRedisCommand("ZREM", bucket, jobId)
	return err
}
