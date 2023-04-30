package providers

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

type QueueServiceProvider struct {
	*support.BaseServiceProvider
}

func (receiver *QueueServiceProvider) Register() {
	facades.Queue.Register(receiver.Jobs())
}

func (receiver *QueueServiceProvider) Boot() {

}

func (receiver *QueueServiceProvider) Name() string {
	return "QueueServiceProvider"
}

func (receiver *QueueServiceProvider) Jobs() []queue.Job {
	return []queue.Job{}
}
