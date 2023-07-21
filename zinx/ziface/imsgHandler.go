package ziface

type IMsgHandler interface {
	//DoMsgHandler Scheduling/executing the corresponding routine message processing method
	DoMsgHandler(request IRequest)
	//AddRouter Add specific processing logic for the message
	AddRouter(msgId uint32, router IRouter)
	//StartWorkerPool Start the worker pool
	StartWorkerPool()
	//SendMsgToTaskQueue msg -> TaskQueue ,Processed by Worker
	SendMsgToTaskQueue(request IRequest)
}
