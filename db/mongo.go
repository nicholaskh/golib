package db

import (
	"fmt"
	"labix.org/v2/mgo"
	"time"
)

var (
	mgoSession *mgo.Session
)

type MongoInfo struct {
	SyncTimeout   time.Duration
	SocketTimeout time.Duration
}

func MgoSession(addr string, info MongoInfo) *mgo.Session {
	if mgoSession == nil {
		mgoSession = getMongoSession(addr, info)
	}

	return mgoSession
}

func getMongoSession(addr string, info MongoInfo) *mgo.Session {
	var err error
	mgoSession, err = mgo.Dial(addr)
	if err != nil {
		panic(fmt.Sprintf("Connect to mongo error: %s", err.Error()))
	}
	if info.SyncTimeout != 0 {
		mgoSession.SetSyncTimeout(info.SyncTimeout)
	}
	if info.SocketTimeout != 0 {
		mgoSession.SetSocketTimeout(info.SocketTimeout)
	} else if info.SyncTimeout != 0 {
		mgoSession.SetSocketTimeout(info.SyncTimeout)
	}
	return mgoSession
}

type MgoSessionPool struct {
	pool chan *mgo.Session
}

func NewMgoSessionPool(addr string, size int, info MongoInfo) *MgoSessionPool {
	this := new(MgoSessionPool)
	this.pool = make(chan *mgo.Session, size)
	for i := 0; i < size; i++ {
		mongoSession := getMongoSession(addr, info)
		this.pool <- mongoSession
	}

	return this
}

func (this *MgoSessionPool) Get() *mgo.Session {
	return <-this.pool
}

func (this *MgoSessionPool) Put(connection *mgo.Session) {
	this.pool <- connection
}
