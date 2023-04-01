package commands

import (
	"context"
	"sync"
	"time"

	"github.com/Coderx44/gg/apperrors"
	"github.com/Coderx44/gg/domain"
)

type Datastore struct {
	sync.RWMutex
	Cond  *sync.Cond
	data  map[string]*DataItem
	queue map[string][]string
}

type DataItem struct {
	value      string
	expiryTime time.Time
}

func NewDatastore() *Datastore {
	return &Datastore{
		data:  make(map[string]*DataItem),
		queue: make(map[string][]string),
		Cond:  sync.NewCond(&Datastore{}),
	}
}

func (ds *Datastore) Get(ctx context.Context, key string) (string, error) {
	ds.RLock()
	defer ds.RUnlock()

	item, ok := ds.data[key]
	if !ok {
		return "", apperrors.ErrKeyNotFound
	}

	if !item.expiryTime.IsZero() && item.expiryTime.Before(time.Now()) {
		// If the item has expired, delete it and return false
		delete(ds.data, key)
		return "", apperrors.ErrKeyNotFound
	}

	return item.value, nil
}

func (ds *Datastore) Set(ctx context.Context, info domain.SetCommand) error {
	ds.Lock()
	defer ds.Unlock()

	item, ok := ds.data[info.Key]

	if ok && !item.expiryTime.IsZero() && item.expiryTime.Before(time.Now()) {
		// If the item has expired, delete it and return false
		delete(ds.data, info.Key)
	}
	if info.Condition == "NX" && ok {
		return apperrors.ErrKeyExists
	}
	if info.Condition == "XX" && !ok {
		return apperrors.ErrKeyNotFound
	}
	if info.Expiry_time != -1 {
		ds.data[info.Key] = &DataItem{
			value:      info.Value,
			expiryTime: time.Now().Add(time.Duration(info.Expiry_time) * time.Second),
		}
	} else {
		ds.data[info.Key] = &DataItem{
			value:      info.Value,
			expiryTime: time.Time{},
		}
	}
	return nil
}
func (ds *Datastore) Bqpop(ctx context.Context, key string, timeout float64) (string, error) {

	ds.Lock()
	defer ds.Unlock()

	values, ok := ds.queue[key]
	if !ok {
		return "", apperrors.ErrKeyNotFound
	}

	// Wait for a value to be added to the queue or the timeout to expire
	expiryTime := time.Now().Add(time.Duration(timeout) * time.Second)
	go func() {
		for len(values) == 0 {

			if timeout == 0 {
				ds.Cond.Broadcast()
				return
			}

			// Check if the timeout has expired
			if time.Now().After(expiryTime) {
				ds.Cond.Broadcast()
				return
			}

		}
	}()
	for len(values) == 0 {
		// fmt.Println("timeout")
		if timeout == 0 {
			return "", apperrors.ErrQueueEmpty
		}

		// Check if the timeout has expired
		if time.Now().After(expiryTime) {
			return "", apperrors.ErrQueueTimeout
		}

		// Wait for a signal to wake up or the timeout to expire
		ds.Cond.Wait()
		values = ds.queue[key]
	}

	value := values[len(values)-1]

	ds.queue[key] = values[:len(values)-1]

	return value, nil
}

func (ds *Datastore) Qpush(ctx context.Context, key string, list []string) {

	// Process the request
	ds.Lock()
	defer ds.Unlock()

	if ds.queue[key] == nil {
		ds.queue[key] = make([]string, 0)
	}

	ds.Cond.Broadcast()
	ds.queue[key] = append(ds.queue[key], list...)

}

func (ds *Datastore) Qpop(ctx context.Context, key string) (value string) {
	ds.Lock()
	defer ds.Unlock()

	values := ds.queue[key]
	if len(values) == 0 {
		return
	}

	value = values[len(values)-1]

	ds.queue[key] = values[:len(values)-1]
	return
}
