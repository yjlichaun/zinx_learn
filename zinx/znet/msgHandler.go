package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	//Apis save every msgId -> process method
	Apis map[uint32]ziface.IRouter
	//The message queue responsible for the worker fetch task
	TaskQueue []chan ziface.IRequest
	//The number of workers in the pool of business workers
	WorkerPoolSize uint32
}

//NewMsgHandler create a new MsgHandler
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //Obtained from the global configuration
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//DoMsgHandler Scheduling/executing the corresponding routine message processing method
func (h *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	//find msgId from Request
	handle, ok := h.Apis[req.GetMsgId()]
	if !ok {
		fmt.Printf("No such msgId: %d\n ,need register", req.GetMsgId())
		return
	}
	//Scheduling router with messageID
	handle.PreHandle(req)
	handle.Handle(req)
	handle.PostHandle(req)
}

//AddRouter Add specific processing logic for the message
func (h *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	//judge if the msgId:IRouter already exists
	if _, ok := h.Apis[msgId]; ok {
		fmt.Println("repeat api , msgId : ", strconv.Itoa(int(msgId)))
	}
	h.Apis[msgId] = router
	fmt.Println("add api, msgId : ", msgId, "success")
}

//StartWorkerPool start a Worker pool(only once)
func (h *MsgHandler) StartWorkerPool() {
	//open a Worker with the worker pool size,every Worker supports with a goroutine
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		//start a worker
		//create a new channel for this worker
		h.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go h.StartOneWorker(i, h.TaskQueue[i])
	}
}

//StartOneWorker start a Worker work flow
func (h *MsgHandler) StartOneWorker(WorkerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerId: ", WorkerId, "start")
	//blocking queue to wait it`s Queue`s msg
	for {
		select {
		//if msg is come , pop a Client Request, and execute the corresponding routine
		case req := <-taskQueue:
			h.DoMsgHandler(req)
		}
	}
}

//SendMsgToTaskQueue msg -> TaskQueue ,Processed by Worker
func (h *MsgHandler) SendMsgToTaskQueue(req ziface.IRequest) {
	//Distribute equally msg to the UnPass worker
	workerId := req.GetConnection().GetConnID() % h.WorkerPoolSize
	fmt.Println("Add ConnID: ", req.GetConnection().GetConnID(),
		"req MsgId", req.GetMsgId(),
		" to WorkerId: ", workerId)
	//msg -> worker TaskQueue
	h.TaskQueue[workerId] <- req
}
