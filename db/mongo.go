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

