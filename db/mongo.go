package db

import (
	"fmt"
	"labix.org/v2/mgo"
)

var (
	mgoSession *mgo.Session
)

func MgoSession(addr string) *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(addr)
		if err != nil {
			panic(fmt.Sprintf("Connect to mongo error: %s", err.Error()))
		}
	}
	return mgoSession
}

type MgoSessionPool struct {
	pool chan *mgo.Session
}

func NewMgoSessionPool(addr string, size int) *MgoSessionPool {
	this := new(MgoSessionPool)
	this.pool = make(chan *mgo.Session, size)
	for i := 0; i < size; i++ {
		mgoSession, err := mgo.Dial(addr)
		if err != nil {
			panic(fmt.Sprintf("Connect to mongo error: %s", err.Error()))
		}
		this.pool <- mgoSession
	}

	return this
}

func (this *MgoSessionPool) Get() *mgo.Session {
	return <-this.pool
}

func (this *MgoSessionPool) Put(connection *mgo.Session) {
	this.pool <- connection
}
