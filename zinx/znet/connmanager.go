package znet

import (
	"fmt"
	"sync"
	"zinx/ziface"
)

//ConnManager model to manage connections (implement)
type ConnManager struct {
	//connMap map to save already created connections
	connMap map[uint32]ziface.IConnection
	//connLock to protect ConnMap
	connLock sync.RWMutex
}

//NewConnManager creates a new ConnManager
func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]ziface.IConnection),
	}
}

//AddConn add connection
func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	// protect map and lock
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// add connection to map
	cm.connMap[conn.GetConnID()] = conn
	fmt.Println("Connection add to map successfully : conn num = ", cm.GetConnNum(), "conn id = ", conn.GetConnID())
}

//DeleteConn delete connection
func (cm *ConnManager) DeleteConn(conn ziface.IConnection) {
	// protect map and lock
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// delete connection from map
	delete(cm.connMap, conn.GetConnID())
	fmt.Println("Connection delete from map successfully : conn num = ", cm.GetConnNum(), "conn id = ", conn.GetConnID())
}

//GetConn get connection with connId
func (cm *ConnManager) GetConn(connId uint32) (ziface.IConnection, error) {
	//protect map and add read lock
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	// get connection from map
	if conn, ok := cm.connMap[connId]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("connection not found")
}

//GetConnNum get number of connection
func (cm *ConnManager) GetConnNum() int {
	return len(cm.connMap)
}

//CleanAllConn clean all connection
func (cm *ConnManager) CleanAllConn() {
	//protect map add write	lock
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//delete conn and stop connection work
	for connId, conn := range cm.connMap {
		//stop connection
		conn.Stop()
		//delete connection from map
		delete(cm.connMap, connId)
	}
	fmt.Println("All connection clean successfully : conn num = ", cm.GetConnNum())
}
