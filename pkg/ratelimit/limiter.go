package ratelimit

import (
	"log"
	"net/http"
	"sync"
	"time"	
	"net"

	errs "github.com/DoubleChuang/Drl/pkg/errors"
	
)

//ConnectLimiter contains token bucket map and max connect and bucket fill during config
type ConnectLimiter struct {
	maxConnection    int64
	tokenBuckets     sync.Map
	bucketFillDuring time.Duration
}

const (
	MaxConnection    = 60
	BucketFillDuring = 60 * time.Second
)


var (
	connectLimiter *ConnectLimiter
	once           sync.Once
)
//NewConnectLimiter is used to create ConnectLimiter
func GetConnectLimiter() *ConnectLimiter {
	once.Do(func() {
		connectLimiter = &ConnectLimiter{
			maxConnection:    MaxConnection,
			bucketFillDuring: BucketFillDuring,
		}
	})
	return connectLimiter
}
func remoteIp(req *http.Request) string {
    remoteAddr := req.RemoteAddr
    if ip := req.Header.Get("Remote_addr"); ip != "" {
        remoteAddr = ip
    } else {
        remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
    }

    if remoteAddr == "::1" {
        remoteAddr = "127.0.0.1"
    }
    return remoteAddr
}
//Take is used to available
func (cl *ConnectLimiter) Take(r *http.Request) (int64, error) {
	var (
		ok          bool
		v           interface{}
		tokenBucket *Bucket
		available   int64
	)

	v, _ = cl.tokenBuckets.LoadOrStore(remoteIp(r), NewBucket(cl.bucketFillDuring, cl.maxConnection))

	if tokenBucket, ok = v.(*Bucket); !ok {
		return 0, errs.ERR_GET_TOKEN_BUCKET_FAIL
	}
	if available = tokenBucket.TakeOnce(); available <= 0 {
		return 0, errs.ERR_NO_ENOUGH_TOKEN_BUCKET
	}
	return available, nil
}


//Get is used to available
func (cl *ConnectLimiter) Get(r *http.Request) (int64, error) {
	var ok bool
	var v interface{}
	var tokenBucket *Bucket

	var available int64 = cl.maxConnection - 1
	
	v, ok = cl.tokenBuckets.Load(remoteIp(r))

	if !ok {
		return 0, errs.ERR_GET_TOKEN_BUCKET_FAIL
	}

	if tokenBucket, ok = v.(*Bucket); !ok {
		return 0, errs.ERR_GET_TOKEN_BUCKET_FAIL
	}
	if available = tokenBucket.Get(); available <= 0 {
		return 0, errs.ERR_NO_ENOUGH_TOKEN_BUCKET
	}
	log.Println("Get available:", available)
	return available, nil
}
