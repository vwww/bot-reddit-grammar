package storage

import (
	"context"

	"google.golang.org/appengine/datastore"
)

type ProcessingOffset struct {
	Last string `datastore:"last,noindex"`
}

func offsetKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "ProcessingOffset", "", 1337, nil)
}

func GetOffset(c context.Context) (*ProcessingOffset, error) {
	k := offsetKey(c)
	var p ProcessingOffset

	if err := datastore.Get(c, k, &p); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}

	return &p, nil
}

func SetOffset(c context.Context, offset *ProcessingOffset) error {
	k := offsetKey(c)
	_, err := datastore.Put(c, k, offset)
	return err
}

type StoredAuth struct {
	ModHash string `datastore:"uh,noindex"`
	Session string `datastore:"s,noindex"`
}

var cachedStoredAuth *StoredAuth

func authKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "StoredAuth", "", 1337, nil)
}

func GetAuth(c context.Context) (*StoredAuth, error) {
	k := authKey(c)
	var a StoredAuth

	if err := datastore.Get(c, k, &a); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	}

	return &a, nil
}

func SetAuth(c context.Context, auth *StoredAuth) error {
	k := authKey(c)
	_, err := datastore.Put(c, k, auth)
	return err
}
