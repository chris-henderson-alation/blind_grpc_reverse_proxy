package grpcinverter

import (
	"fmt"
	"sync"
	"time"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentId = uint64
type JobId = uint64

type Agents struct {
	listeners    map[AgentId]chan *GrpcCall
	awaitingJobs map[AgentId]*EndangeredJob
	lock         sync.Mutex
}

func NewAgents() *Agents {
	return &Agents{
		listeners:    map[AgentId]chan *GrpcCall{},
		awaitingJobs: map[AgentId]*EndangeredJob{},
		lock:         sync.Mutex{},
	}
}

func (a *Agents) Register(agent AgentId) (chan *GrpcCall, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, alreadyListening := a.listeners[agent]; alreadyListening {
		return nil, true
	}
	channel := make(chan *GrpcCall)
	a.listeners[agent] = channel
	return channel, false
}

func (a *Agents) Unregister(agent AgentId) {
	a.lock.Lock()
	delete(a.listeners, agent)
	a.lock.Unlock()
}

func (a *Agents) Listening(agent AgentId) bool {
	a.lock.Lock()
	_, listening := a.listeners[agent]
	a.lock.Unlock()
	return listening
}

func (a *Agents) Submit(call *GrpcCall) error {
	a.lock.Lock()
	listener, listening := a.listeners[call.Agent]
	a.lock.Unlock()
	if !listening {
		return status.Error(codes.Unavailable, fmt.Sprintf("agent %d is not currently connected", call.Agent))
	}
	listener <- call
	return nil
}

func (a *Agents) Enqueue(call *GrpcCall, bookmark *ioc.Message) {
	endangeredJob := NewEndangeredJob(call, bookmark)
	a.lock.Lock()
	a.awaitingJobs[call.JobId] = endangeredJob
	a.lock.Unlock()
	go endangeredJob.CountdownToDeath(a)
}

func (a *Agents) Retrieve(agentId AgentId, jobId JobId) (*GrpcCall, *ioc.Message) {
	a.lock.Lock()
	defer a.lock.Unlock()
	job, ok := a.awaitingJobs[jobId]
	if !ok {
		// @TODO
		return nil, nil
	}
	if job.call.Agent != agentId {
		// @TODO
		return nil, nil
	}
	job.Rescue()
	delete(a.awaitingJobs, jobId)
	return job.call, job.bookmark
}

type EndangeredJob struct {
	call          *GrpcCall
	bookmark      *ioc.Message
	stopCountdown chan struct{}
}

func NewEndangeredJob(call *GrpcCall, bookmark *ioc.Message) *EndangeredJob {
	return &EndangeredJob{
		call:          call,
		stopCountdown: make(chan struct{}),
		bookmark:      bookmark,
	}
}

func (n *EndangeredJob) CountdownToDeath(parent *Agents) {
	stopwatch := time.NewTicker(time.Second * 30)
	select {
	case <-n.stopCountdown:
	case <-stopwatch.C:
		parent.lock.Lock()
		defer parent.lock.Unlock()
		_, present := parent.awaitingJobs[n.call.JobId]
		if !present {
			// This is technically possible if
			// someone manages to snipe a call to Retrieve
			// in the nanoseconds between this branch
			// beginning execution and the capture of
			// the lock.
			//
			// This is a very action hero type scene
			// where our protagonist escapes death.
			return
		}
		delete(parent.awaitingJobs, n.call.JobId)
		n.call.SendError(status.New(codes.Unavailable, "agent failed to callback within time tolerance"))
	}
}
func (n *EndangeredJob) Rescue() {
	close(n.stopCountdown)
}
