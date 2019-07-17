package main

import (
	"github.com/tsuna/gohbase"
	riak "github.com/basho/riak-go-client"
	"strings"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"github.com/juju/errors"
	"fmt"
	"encoding/json"
	"bytes"
)

type Client interface {
	getName() string
	build() error
	setEndpoint(string)
	open() error
	close() error
}

type HBaseClient struct {
	zkQuorum string
	timeout int32
	h_client gohbase.Client
}

type RiakClient struct {
	remoteAddresses string
	timeout int32
	r_cluster *riak.Cluster
}

func (client *HBaseClient) getName() string {
	return "HBase"
}

func (client *HBaseClient) build() error {
	client.h_client = gohbase.NewClient(client.zkQuorum)
	return nil
}

func (client *HBaseClient) setEndpoint(zk string) {
	client.zkQuorum = zk
}

func (client *HBaseClient) open() error {
	return nil
}

func (client *HBaseClient) close() error {
	client.h_client.Close()
	return nil
}

func (client *RiakClient) getName() string {
	return "Riak"
}

func (client *RiakClient) build() error {
	addresses := strings.Split(client.remoteAddresses, ",")
	var nodes []*riak.Node
	for _, addr := range addresses {
		nodeOpt := &riak.NodeOptions {
			RemoteAddress: addr,
		}
		var node *riak.Node
		var err error
		if node, err = riak.NewNode(nodeOpt); err != nil {
			return err
		}
		nodes = append(nodes, node)
	}
	opts := &riak.ClusterOptions {
		Nodes: nodes,
	}
	cluster, err := riak.NewCluster(opts)
	client.r_cluster = cluster
	return err
}

func (client *RiakClient) setEndpoint(addresses string) {
	client.remoteAddresses = addresses
}

func (client *RiakClient) open() error {
	if err := client.r_cluster.Start(); err != nil {
		return err
	}
	return nil
}

func (client *RiakClient) close() error {
	if err := client.r_cluster.Stop(); err != nil {
		return err
	}
	return nil
}

type DbRequest interface {
	get(client Client, key string, table string) (interface{}, error) // Try change string to DbRequest: get(client Client, key string, table string) (DbRequest, error)
	put(client Client, table string) error
}

type HBaseRow struct {
	Key string `json:"row"`
	ColFamily string `json:"familiy"`
	Column string `json:"column"`
	Value string `json:"value"`
}

type HBaseRows struct {
	rows []*HBaseRow
}

type RiakObject struct {
	Bucket string `json:"bucket"`
	Key string `json:"key"`
	Value string `json:"value"`
}

func (records *HBaseRows) get(client Client, key string, table string) (string, error) {
	if cl, ok := client.(*HBaseClient); ok {
		getRequest, _ := hrpc.NewGetStr(context.Background(), table, key)
		getRsp, err := cl.h_client.Get(getRequest)
		if err != nil {
			return "", err
		}
		var results []*HBaseRow
		for i, c := range getRsp.Cells {
			var record *HBaseRow
			fmt.Println("cell", i, "value", string(c.Value), "row", string(c.Row), "column", string(c.Qualifier), "familiy", string(c.Family))
			record = &HBaseRow{Key: string(c.Row), Column: string(c.Qualifier), ColFamily: string(c.Family), Value: string(c.Value)}
			results = append(results, record)
		}
		fmt.Println("row len:", len(results))
		recordJson, err := json.Marshal(results)
		fmt.Println("record json:", string(recordJson))
		if err != nil {
			return "", err
		}
		return string(recordJson), nil
	}
	return "", errors.New("HBaseClient assertion failed")
}

func (records *HBaseRows) put(client Client, table string) error {
	for _, rec := range records.rows {
		byteArray := []byte(rec.Value)
		values := map[string]map[string][]byte{rec.ColFamily: map[string][]byte{rec.Column: byteArray}}
		putRequest, _ := hrpc.NewPutStr(context.Background(), table, rec.Key, values)
		if cl, ok := client.(*HBaseClient); ok {
			_, err := cl.h_client.Put(putRequest)
			if err != nil {
				return err
			}
		} else {
			return errors.New("HBaseClient assertion failed")
		}
	}
	return nil
}


func (record *RiakObject) get(client Client, key string, table string) (interface{}, error) {
	//cmd, err := riak.NewFetchValueCommandBuilder().WithBucket(table).WithKey(key).Build()
	binaryKey, _ := StringToUuid(key)
	cmd, err := riak.NewFetchValueCommandBuilder().WithBucket(table).WithBinaryKey(binaryKey).Build()
	if err != nil {
		return "", err
	}
	if cl, ok := client.(*RiakClient); ok {
		if err := cl.r_cluster.Execute(cmd); err != nil {
			return "", err
		}
		if svc, ok := cmd.(*riak.FetchValueCommand); ok {
			valLen := len(svc.Response.Values)
			if valLen != 0 {
				/*
				obj := svc.Response.Values[0]
				fmt.Println("result value:", string(obj.Value))
				var newObj *RiakObject
				newObj = &RiakObject{Value: string(obj.Value), Bucket: table, Key: key}
				recordJson, err := json.Marshal(newObj)
				*/
				rb := svc.Response.Values[0].Value
				r := bytes.NewReader(rb)
				term, err := Read(r)
				if err != nil {
					return nil, err
				}
				return term, nil
				//return string(recordJson), nil
			} else {
				return nil, errors.New("Empty Result")
			}
		} else {
			return nil, errors.New("FetchValueCommand assertion failed")
		}
	} else {
		return nil, errors.New("RiakClient assertion failed")
	}
	return nil, nil
}

func (record *RiakObject) put(client Client, table string) error {
	obj := &riak.Object{Bucket: record.Bucket, Key: record.Key, ContentType: "application/octet-stream", Value: []byte(record.Value),}
	cmd, err := riak.NewStoreValueCommandBuilder().WithContent(obj).Build()
	if err != nil {
		return err
	}
	if cl, ok := client.(*RiakClient); ok {
		if err := cl.r_cluster.Execute(cmd); err != nil {
			return err
		}
	} else {
		return errors.New("RiakClient aseertion failed")
	}
	return nil
}


func getRequest(request DbRequest, client Client, key string, table string) (interface{}, error) {
	return request.get(client, key, table)
}

func putRequest(request DbRequest, client Client, table string) error {
	return request.put(client, table)
}
