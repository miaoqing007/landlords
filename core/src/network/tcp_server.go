////=============================================================================
////     FileName: tcp_server.go
////         Desc: 简单的tcpserver实现
////       Author: chenshungen
////        Email: 953524751@qq.com
////     HomePage: http://www.chenshungen.cn
////      Version: 0.0.1
////   LastChange: 2018-11-29 16:25:42
////      History:
////=============================================================================
package network
//
//import (
//	"net"
//	"os"
//	"reflect"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//type ConnSet map[net.Conn]interface{}
//
//type TcpServer struct {
//	pid          int64
//	Addr         string
//	MaxConnNum   int
//	ln           *net.TCPListener
//	conns        ConnSet
//	counter      int64
//	idCounter    int64
//	mutexConns   sync.Mutex
//	wgLn         sync.WaitGroup
//	wgConns      sync.WaitGroup
//	connBuffSize int
//
//	NewClient func(conn net.Conn) IConn
//}
//
//// Init ...
//func (srv *TcpServer) Init(addr string, maxConnNum int, buffSize int64) {
//
//	srv.MaxConnNum = maxConnNum
//	srv.Addr = addr
//
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", srv.Addr)
//
//	if err != nil {
//		//logger.Fatal("[net] addr resolve error", tcpAddr, err)
//	}
//
//	ln, err := net.ListenTCP("tcp", tcpAddr)
//
//	if err != nil {
//		//logger.Fatalf("%v", err)
//	}
//
//	if srv.MaxConnNum <= 0 {
//		srv.MaxConnNum = 100
//		//logger.Warningf("invalid MaxConnNum, reset to %v", srv.MaxConnNum)
//	}
//
//	srv.ln = ln
//	srv.conns = make(ConnSet)
//	srv.counter = 1
//	srv.idCounter = 1
//	srv.pid = int64(os.Getpid())
//	srv.connBuffSize = int(buffSize)
//	//logger.Infof("TcpServer Listen %s", srv.ln.Addr().String())
//}
//
//// Run ...
//func (srv *TcpServer) Run() {
//	// 捕获异常
//	defer func() {
//		if err := recover(); err != nil {
//			//logger.Error("[net] panic", err, "\n", string(debug.Stack()))
//		}
//	}()
//
//	srv.wgLn.Add(1)
//	defer srv.wgLn.Done()
//
//	var tempDelay time.Duration
//	for {
//		conn, err := srv.ln.AcceptTCP()
//
//		if err != nil {
//			if ne, ok := err.(net.Error); ok && ne.Temporary() {
//				if tempDelay == 0 {
//					tempDelay = 5 * time.Millisecond
//				} else {
//					tempDelay *= 2
//				}
//				if max := 1 * time.Second; tempDelay > max {
//					tempDelay = max
//				}
//				//logger.Warningf("accept error: %v; retrying in %v", err, tempDelay)
//				time.Sleep(tempDelay)
//				continue
//			}
//			return
//		}
//		tempDelay = 0
//
//		if atomic.LoadInt64(&srv.counter) >= int64(srv.MaxConnNum) {
//			conn.Close()
//			//logger.Warning("too many connections %v", atomic.LoadInt64(&srv.counter))
//			continue
//		}
//
//		// Try to open keepalive for tcp.
//		conn.SetKeepAlive(true)
//		conn.SetKeepAlivePeriod(1 * time.Minute)
//		// disable Nagle's algorithm.
//		conn.SetNoDelay(true)
//		conn.SetWriteBuffer(srv.connBuffSize)
//		conn.SetReadBuffer(srv.connBuffSize)
//		client := srv.NewClient(conn)
//		srv.createNewConn(conn, client)
//
//		srv.wgConns.Add(1)
//		go func() {
//			client.Run()
//			// cleanup Run执行结束 连接断开 清理连接管理
//			srv.removeConn(conn, client)
//			srv.wgConns.Done()
//		}()
//	}
//}
//
//func (srv *TcpServer) Close() {
//	srv.ln.Close()
//	srv.wgLn.Wait()
//
//	srv.mutexConns.Lock()
//	for conn := range srv.conns {
//		conn.Close()
//	}
//	srv.conns = nil
//	srv.mutexConns.Unlock()
//	srv.wgConns.Wait()
//}
//
//func (srv *TcpServer) createNewConn(conn net.Conn, client IConn) {
//	srv.mutexConns.Lock()
//	atomic.AddInt64(&srv.counter, 1)
//	srv.conns[conn] = conn
//	nowTime := time.Now().Unix()
//	idCounter := atomic.AddInt64(&srv.idCounter, 1)
//	ccid := (nowTime << 32) | (srv.pid << 24) | idCounter
//	v := reflect.ValueOf(client).Elem()
//	v.FieldByName("ConnID").Set(reflect.ValueOf(ccid))
//	srv.mutexConns.Unlock()
//	client.OnConnect() // 去找实现的功能
//}
//
//func (srv *TcpServer) removeConn(conn net.Conn, client IConn) {
//	client.Close()
//	srv.mutexConns.Lock()
//	atomic.AddInt64(&srv.counter, -1)
//	delete(srv.conns, conn)
//	srv.mutexConns.Unlock()
//}
