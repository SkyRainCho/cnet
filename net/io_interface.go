package cnet

type ReadInterface interface {
	Read(*Session, []byte) (interface{}, int, error)
}

type WriteInterface interface {
	Write(*Session, interface{}) error
}

type IOHandler interface {
	ReadInterface
	WriteInterface
}

type EventHandler interface {
	//建立连接触发的回调函数
	OnConnect(*Session) error
	//断开连接触发的回调函数
	OnDisconnect(*Session)
	//连接被异常中断的回调函数
	OnAbortConnect(*Session, error)
	//连接上的心跳回调函数
	OnHeartbeat(*Session)
	//收到一条完整消息的回调函数ele
	OnHandleMsg(*Session, interface{})
}
