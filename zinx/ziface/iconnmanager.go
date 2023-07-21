package ziface

//IConnManager model to manage connections (abstract)
type IConnManager interface {
	//AddConn add connection
	AddConn(conn IConnection)
	//DeleteConn delete connection
	DeleteConn(conn IConnection)
	//GetConn get connection with connId
	GetConn(connId uint32) (IConnection, error)
	//GetConnNum get number of connection
	GetConnNum() int
	//CleanAllConn clean all connection
	CleanAllConn()
}
