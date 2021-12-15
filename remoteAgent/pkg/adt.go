package grpcinverter

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentId = uint64
type JobId = uint64

type ListenerMap struct {
	listeners map[AgentId]chan *GrpcCall
	lock      sync.Mutex
}

func NewListenerMap() *ListenerMap {
	return &ListenerMap{
		listeners: map[AgentId]chan *GrpcCall{},
		lock:      sync.Mutex{},
	}
}

func (a *ListenerMap) Register(agent AgentId) (chan *GrpcCall, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, alreadyListening := a.listeners[agent]; alreadyListening {
		return nil, true
	}
	channel := make(chan *GrpcCall)
	a.listeners[agent] = channel
	return channel, false
}

func (a *ListenerMap) Unregister(agent AgentId) {
	a.lock.Lock()
	delete(a.listeners, agent)
	a.lock.Unlock()
}

func (a *ListenerMap) Listening(agent AgentId) bool {
	a.lock.Lock()
	_, listening := a.listeners[agent]
	a.lock.Unlock()
	return listening
}

// Submit is used by the upstream recipient of Alation's gRPC requests to submit
// the job to the agent listening channel.
//
// If no such agent is currently connected then an error is immediately returned describing as such.
func (a *ListenerMap) Submit(call *GrpcCall) error {
	a.lock.Lock()
	listener, listening := a.listeners[call.Agent]
	a.lock.Unlock()
	if !listening {
		return status.Error(codes.Unavailable, fmt.Sprintf("agent %d is not currently connected", call.Agent))
	}
	listener <- call
	return nil
}

type PendingJobs struct {
	pending map[AgentId]*EndangeredJob
	lock    sync.Mutex
}

func NewPendingJobs() *PendingJobs {
	return &PendingJobs{
		pending: map[AgentId]*EndangeredJob{},
		lock:    sync.Mutex{},
	}
}

// Submit
func (a *PendingJobs) Submit(call *GrpcCall) {
	endangeredJob := NewEndangeredJob(call)
	a.lock.Lock()
	a.pending[call.JobId] = endangeredJob
	a.lock.Unlock()
	go endangeredJob.CountdownToDeath(a)
}

// Retrieve removes and returns the GrpcCall associated with the agent ID and job ID.
//
// If no such job ID exists (possibly because the EndangeredJob timed out) then nil is returned.
// If the job exists, but it is not associated with the given agent ID then nil is returned.
func (a *PendingJobs) Retrieve(agentId AgentId, jobId JobId) *GrpcCall {
	a.lock.Lock()
	defer a.lock.Unlock()
	job, ok := a.pending[jobId]
	if !ok {
		return nil
	}
	if job.call.Agent != agentId {
		logrus.Warnf("agent %d attempted to access job %d, however job %d belongs to agent %d!",
			agentId, jobId, jobId, job.call.Agent)
		return nil
	}
	job.Rescue()
	delete(a.pending, jobId)
	return job.call
}

// EndangeredJob are jobs that the upstream Alation has requested and that the forward proxy has fed
// into the agent. HOWEVER, there is a brief period of time in which we are waiting for the agent to establish
// a new stream to send the results over. During that time, we do not know whether-or-not the agent will ever
// successfully call us back to being streaming!
//
// As such, we record the upstream GrpcCall, which is keeping the stream to Alation open, and we patiently wait
// for the agent to call back. If the agent does not callback within the configured tolerance then Alation
// will be sent a plain gRPC error message informing it on what had happened.
type EndangeredJob struct {
	call          *GrpcCall
	stopCountdown chan struct{}
}

func NewEndangeredJob(call *GrpcCall) *EndangeredJob {
	return &EndangeredJob{
		call:          call,
		stopCountdown: make(chan struct{}),
	}
}

// CountdownToDeath is a dramatic name for a simple concept - the agent has two minutes to call back
// or else it's job will be killed by the forward proxy and Alation will be sent a message as such.
//
// This grizzly fate cane be avoided by calling Rescue on the EndangeredJob.
func (n *EndangeredJob) CountdownToDeath(parent *PendingJobs) {
	stopwatch := time.NewTicker(time.Minute * 2)
	select {
	case <-n.stopCountdown:
	case <-stopwatch.C:
		parent.lock.Lock()
		defer parent.lock.Unlock()
		_, present := parent.pending[n.call.JobId]
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
		delete(parent.pending, n.call.JobId)
		n.call.Alation.SendError(status.New(codes.Unavailable, "agent failed to callback within time tolerance"))
	}
}
func (n *EndangeredJob) Rescue() {
	close(n.stopCountdown)
}
