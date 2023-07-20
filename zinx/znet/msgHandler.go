package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandler struct {
	//Apis save every msgId -> process method
	Apis map[uint32]ziface.IRouter
}

//NewMsgHandler create a new MsgHandler
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
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
