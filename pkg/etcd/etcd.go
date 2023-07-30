package etcd

import (
	"context"
	"fmt"
	"time"

	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	client *clientv3.Client
}

func NewClient() (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		client: cli,
	}, nil
}

func (c *Client) Get(key string) (value *mvccpb.KeyValue, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := c.client.Get(ctx, key)
	defer cancel()
	if err != nil {
		return &mvccpb.KeyValue{}, err
	}
	if len(resp.Kvs) == 0 {
		return &mvccpb.KeyValue{}, fmt.Errorf("unable to find key '%v'", key)
	}
	return resp.Kvs[0], nil
}

func (c *Client) ListWithPrefix(prefix string) (kvs []*mvccpb.KeyValue, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := c.client.Get(ctx, prefix, clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		return nil, err
	}
	return resp.Kvs, nil
}

func (c *Client) Put(key string, value string) (resp *clientv3.PutResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	val, err := c.client.Put(ctx, key, value)
	defer cancel()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *Client) Delete(key string) (resp *clientv3.DeleteResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err = c.client.Delete(ctx, key)
	defer cancel()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteWithPrefix(prefix string) (resp *clientv3.DeleteResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err = c.client.Delete(ctx, prefix, clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
