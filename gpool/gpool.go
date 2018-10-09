package gpool

import (
    "context"
    "sync"
    "time"
    "errors"

    "google.golang.org/grpc"
)

var (
    // ErrInvalidConfig  config invalid
    ErrInvalidConfig = errors.New("gpool : invalid pool config")
    // ErrTimeout get conn timeout
    ErrTimeout       = errors.New("gpool : get conn timeout")
    // ErrPoolClosed pool has been closed
    ErrPoolClosed    = errors.New("gpool : pool is closed")
    // ErrPoolFulled pool is fulled
    ErrPoolFulled    = errors.New("gpool : pool is full")
)

// Factory function type to create grpc.ClientConn object
type Factory func() (*grpc.ClientConn, error)

// Conn a wrapper of grpc client
type Conn struct {
    C         *grpc.ClientConn  // connection object
    lastUsed  time.Time         // timestamp used last time
}

// GPool simillar to nginx connection with idle timeout
type GPool struct {
    factory   Factory        // factory method to create connection
    conns     chan Conn      // connections' chan

    init      uint32         // init pool size
    capacity  uint32         // max pool size

    maxIdle   time.Duration  // how long an idle connection remains open

    rwl       sync.RWMutex   // read-write mutex
}

// NewPool create pool
func NewPool(factory Factory, init, capacity uint32, maxIdle time.Duration) (*GPool, error) {
    // check input args
    if init < 0 || capacity <= 0 || capacity < init {
        return nil, ErrInvalidConfig
    }

    // create pool
    p := &GPool{
        conns :    make(chan Conn, capacity),
        factory :  factory,
        maxIdle :  maxIdle,

        init :     init,
        capacity : capacity,
    }

    // create conns with size : init
    var idx uint32
    for idx = 0; idx < init; idx++ {
        c, err := factory()
        if err != nil {
            return nil, err
        }

        p.conns <- Conn{
            C :        c,
            lastUsed : time.Now(),
        }
    }

    // fill the rest with empty conns
    rest := capacity - init
    for idx = 0; idx < rest; idx++ {
        p.conns <- Conn{}
    }

    return p, nil
}

// Capacity  capacity of pool
func (p* GPool) Capacity() uint32 {
    conns := p.getConnsRLock()
    if conns == nil {
        return 0
    }
    return uint32(cap(conns))
}

// Available  available connection
func (p* GPool) Available() uint32 {
    conns := p.getConnsRLock()
    if conns == nil {
        return 0
    }
    return uint32(len(conns))
}

// get conns
func (p* GPool) getConnsRLock() chan Conn {
    p.rwl.RLock()
    defer p.rwl.RUnlock()

    return p.conns
}

// Get get conn with context
func (p* GPool) Get(ctx context.Context) (*Conn, error) {
    conns := p.getConnsRLock()
    if conns == nil {
        return nil, ErrPoolClosed
    }

    // fetch conn until timeout
    var conn Conn
    select {
        case conn = <-conns:
            // do nothing
        case <-ctx.Done():
            return nil, ErrTimeout
    }

    // check its idle time and reset connection if needed
    if conn.C != nil && p.maxIdle > 0 &&
                        conn.lastUsed.Add(p.maxIdle).Before(time.Now()) {
        conn.C.Close()
        conn.C = nil
    }

    // create one if needed
    var err error
    if conn.C == nil {
        conn.C, err = p.factory()
        if err != nil { // if failed, should return this conn to pool
            conns <- Conn{}
        }
    }
    return &conn, err
}

// Put return conn to pool
func (p *GPool) Put(conn *Conn) error {
    nconn := Conn{
        C :        conn.C,
        lastUsed : time.Now(), // update lastUsed
    }

    select {
        case p.conns <- nconn:
            // do nothing
        default:
            return ErrPoolFulled
    }
    return nil
}

// Close close conns of pool
func (p *GPool) Close() {
    // fetch connes : rwlock
    p.rwl.Lock()
    conns := p.conns
    p.conns = nil
    p.rwl.Unlock()

    // close them
    if conns == nil {
        return
    }

    var idx uint32
    for idx = 0; idx < p.capacity; idx++ {
        conn := <-conns
        if conn.C != nil {
            conn.C.Close()
        }
    }
    close(conns)
}

