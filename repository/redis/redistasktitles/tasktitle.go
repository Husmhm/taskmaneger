package redistasktitles

import (
	"context"
	"encoding/json"
	"time"
)

func (d DB) Set(key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return d.adapter.Client().Set(context.Background(), key, p, expiration).Err()

}

func (d DB) Get(key string, dest interface{}) error {
	p, err := d.adapter.Client().Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(p, dest)
}
