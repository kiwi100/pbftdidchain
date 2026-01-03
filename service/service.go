package service

import (
	"encoding/json"
	"fmt"
	"net"
	"pbftdidchain/message"
)

/*
Our algorithm can be used to implement any deterministic replicated service with a state and some operations. The
operations are not restricted to simple reads or writes of portions of the service state; they can perform arbitrary
deterministic computations using the state and operation arguments. Clients issue requests to the replicated service to
invoke operations and block waiting for a reply. The replicated service is implemented by   replicas. Clients and
replicas are non-faulty if they follow the algorithm in Section 4 and if no attacker can forge their signature.
*/

type Service struct {
	SrvHub      *net.UDPConn
	nodeChan    chan interface{}
	clientAddrs map[string]*net.UDPAddr // 存储客户端地址映射
}

func InitService(port int, msgChan chan interface{}) *Service {
	locAddr := net.UDPAddr{
		Port: port,
	}
	srv, err := net.ListenUDP("udp4", &locAddr)
	if err != nil {
		return nil
	}
	fmt.Printf("\n===>Service Listening at[%d]", port)
	s := &Service{
		SrvHub:      srv,
		nodeChan:    msgChan,
		clientAddrs: make(map[string]*net.UDPAddr),
	}
	return s
}

func (s *Service) WaitRequest(sig chan interface{}) {

	defer func() {
		if r := recover(); r != nil {
			sig <- r
		}
	}()

	buf := make([]byte, 2048)
	for {
		n, rAddr, err := s.SrvHub.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Service received err:%s\n", err)
			continue
		}
		fmt.Printf("\nService message[%d] from[%s]\n", n, rAddr.String())
		bo := &message.Request{}
		if err := json.Unmarshal(buf[:n], bo); err != nil {
			fmt.Printf("\nService message parse err:%s", err)
			continue
		}
		// 保存客户端地址（客户端监听在8088端口）
		clientAddr := &net.UDPAddr{
			IP:   rAddr.IP,
			Port: 8088, // 客户端监听端口
		}
		s.clientAddrs[bo.ClientID] = clientAddr
		go s.process(bo)
	}
}

func (s *Service) process(op *message.Request) {

	/*
		TODO:: Check operation
		1. if clientID is authorized
		2. if operation is valid
	*/
	s.nodeChan <- op
}

/*
	Each replica i executes the operation requested by m  after committed-local(m, v, n, i)is true and i’s state

reflects the sequential execution of all requests with lower sequence numbers. This ensures that all non- faulty replicas
execute requests in the same order as required to provide the safety property. After executing the requested operation,
replicas send a reply to the client. Replicas discard requests whose timestamp is lower than the timestamp in the last
reply they sent to the client to guarantee exactly-once semantics.

	We do not rely on ordered message delivery, and therefore it is possible for a replica to commit requests out

of order. This does not matter since it keeps the pre- prepare, prepare, and commit messages logged until the
corresponding request can be executed.
*/
func (s *Service) Execute(v, n, seq int64, o *message.Request) (reply *message.Reply, err error) {

	fmt.Printf("Service is executing opertion[%s]......\n", o.Operation)
	r := &message.Reply{
		SeqID:     seq,
		ViewID:    v,
		Timestamp: o.TimeStamp,
		ClientID:  o.ClientID,
		NodeID:    n,
		Result:    "success",
	}

	bs, _ := json.Marshal(r)

	// 从保存的客户端地址映射中获取地址
	cAddr, ok := s.clientAddrs[o.ClientID]
	if !ok {
		// 如果找不到客户端地址，使用localhost作为默认值
		cAddr = &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 8088,
		}
		fmt.Printf("Warning: Client address not found for %s, using localhost\n", o.ClientID)
	}

	no, err := s.SrvHub.WriteToUDP(bs, cAddr)
	if err != nil {
		fmt.Printf("Reply client failed:%s\n", err)
		return nil, err
	}
	fmt.Printf("Reply Success!:%d seq=%d\n", no, seq)
	return r, nil
}

func (s *Service) DirectReply(r *message.Reply) error {
	bs, _ := json.Marshal(r)

	// 从保存的客户端地址映射中获取地址
	cAddr, ok := s.clientAddrs[r.ClientID]
	if !ok {
		// 如果找不到客户端地址，使用localhost作为默认值
		cAddr = &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 8088,
		}
		fmt.Printf("Warning: Client address not found for %s, using localhost\n", r.ClientID)
	}

	no, err := s.SrvHub.WriteToUDP(bs, cAddr)
	if err != nil {
		fmt.Printf("Reply client failed:%s\n", err)
		return err
	}
	fmt.Printf("Reply Directly Success!:%d seq=%d\n", no, r.SeqID)
	return nil
}
