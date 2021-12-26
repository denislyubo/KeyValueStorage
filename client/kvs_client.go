package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type KeyValueStorage interface {
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key, val string) error
	//Delete(ctx context.Context, key string)
}

type clientKVS struct {
	conStr string
}

func NewClientKVS(conStr string) KeyValueStorage {
	return &clientKVS{conStr: conStr}
}

func (c *clientKVS) Get(ctx context.Context, key string) (string, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", c.conStr+"/"+key, nil)

	if err != nil {
		return "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	val, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (c *clientKVS) Put(ctx context.Context, key, value string) error {

	bodyRequest := []byte(value)
	bufBody := bytes.NewBuffer(bodyRequest)
	req, err := http.NewRequestWithContext(ctx, "PUT", c.conStr+"/"+key, bufBody)

	if err != nil {
		return err
	}

	client := &http.Client{}

	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	client := NewClientKVS("http://localhost:8080/v1")

	var action, key, value string
	if len(os.Args) > 0 {
		action, key = os.Args[1], os.Args[2]
		value = strings.Join(os.Args[3:], " ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch action {
	case "get":
		val, err := client.Get(ctx, key)

		if err != nil {
			log.Fatalf("could not get value for key %s: %v\n", key, err)
		}

		log.Printf("Get %s returns: %s", key, val)
	case "put":
		err := client.Put(ctx, key, value)
		if err != nil {
			log.Fatalf("could not get put key %s: %v\n", key, err)
		}
		log.Printf("Put %s", key)
	default:
		log.Fatalf("Syntax: go run [get|put] KEY VALUE...")
	}
}
