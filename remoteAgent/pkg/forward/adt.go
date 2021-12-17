package forward

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentId = uint64
type JobId = uint64

// ListenerMap records all agents who are currently registered with this proxy
// and provides a mapping of their agent IDs to a channel on which that agent
// may receive new jobs.
type ListenerMap struct {
	listeners map[AgentId]chan *GrpcCall
	lock      sync.Mutex
}

// NewListenerMap instantiates an empty ListenerMap.
func NewListenerMap() *ListenerMap {
	return &ListenerMap{
		listeners: map[AgentId]chan *GrpcCall{},
		lock:      sync.Mutex{},
	}
}

// Register registers a new listener for the given agent ID and returns a channel
// which may be used to receive incoming jobs for that agent ID.
//
// If a listener is already registered at the provided ID then nil is returned.
// Agents requesting double registration SHOULD be turned away.
//
// A call to Submit(call) corresponds to the listener returned by Register.
// That is to say, a channel returned by Register(32) will receive all incoming jobs
// sent via Submit(call) where the provided call has a GrpcCall.Agent == 32.
func (l *ListenerMap) Register(agent AgentId) chan *GrpcCall {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, alreadyListening := l.listeners[agent]; alreadyListening {
		return nil
	}
	channel := make(chan *GrpcCall)
	l.listeners[agent] = channel
	return channel
}

// Unregister removes the listener for the provided agent ID.
// This typically corresponds with agent disconnect events.
func (l *ListenerMap) Unregister(agent AgentId) {
	l.lock.Lock()
	delete(l.listeners, agent)
	l.lock.Unlock()
}

// Listening returns whether-or-not the given agent is listening already.
// This is currently only used in unit tests to prevent race conditions.
func (l *ListenerMap) Listening(agent AgentId) bool {
	l.lock.Lock()
	_, listening := l.listeners[agent]
	l.lock.Unlock()
	return listening
}

// Submit is used by the upstream recipient of Alation's gRPC requests to submit
// the job to the agent listening channel.
//
// If no such agent is currently connected then an error is immediately returned describing as such.
func (l *ListenerMap) Submit(call *GrpcCall) error {
	l.lock.Lock()
	listener, listening := l.listeners[call.Agent]
	l.lock.Unlock()
	if !listening {
		return status.Error(codes.Unavailable, fmt.Sprintf("agent %d is not currently connected", call.Agent))
	}
	listener <- call
	return nil
}

// PendingJobs holds all jobs which have been submitted to an agent for progress, but for which
// the agent has not yet called back home in order to stream the results of that job.
//
// This is useful for tracking (and killing) jobs that have been submitted to agents who may have
// silently died during that time inbetween the moment the job is sent and the moment an upstream
// reconnect is established.
type PendingJobs struct {
	pending map[AgentId]*PendingJob
	lock    sync.Mutex
}

// NewPendingJobs instantiates an empty PendingJobs.
func NewPendingJobs() *PendingJobs {
	return &PendingJobs{
		pending: map[AgentId]*PendingJob{},
		lock:    sync.Mutex{},
	}
}

// Submit submits a pending job. Doing so immediately begins the countdown towards the deadline
// for an agent to call back home.
func (p *PendingJobs) Submit(call *GrpcCall) {
	endangeredJob := NewPendingJob(call)
	p.lock.Lock()
	p.pending[call.JobId] = endangeredJob
	p.lock.Unlock()
	go endangeredJob.CountdownToDeath(p)
}

// Retrieve removes and returns the GrpcCall associated with the agent ID and job ID.
//
// If no such job ID exists (possibly because the PendingJob timed out) then nil is returned.
// If the job exists, but it is not associated with the given agent ID then nil is returned.
func (p *PendingJobs) Retrieve(agentId AgentId, jobId JobId) *GrpcCall {
	p.lock.Lock()
	defer p.lock.Unlock()
	job, ok := p.pending[jobId]
	if !ok {
		return nil
	}
	if job.call.Agent != agentId {
		// @TODO
		//logrus.Warnf("agent %d attempted to access job %d, however job %d belongs to agent %d!",
		//	agentId, jobId, jobId, job.call.Agent)
		return nil
	}
	job.Rescue()
	delete(p.pending, jobId)
	return job.call
}

// PendingJob are jobs that the upstream Alation has requested and that the forward proxy has fed
// into the agent. HOWEVER, there is a brief period of time in which we are waiting for the agent to establish
// a new stream to send the results over. During that time, we do not know whether-or-not the agent will ever
// successfully call us back to being streaming!
//
// As such, we record the upstream GrpcCall, which is keeping the stream to Alation open, and we patiently wait
// for the agent to call back. If the agent does not callback within the configured tolerance then Alation
// will be sent a plain gRPC error message informing it on what had happened.
type PendingJob struct {
	call          *GrpcCall
	stopCountdown chan struct{}
}

// NewPendingJob instantiates a new PendingJob. The countdown to the deadline does NOT
// begin within the invocatation of this constructor. To begin the countdown, one must call CountDownToDeath.
func NewPendingJob(call *GrpcCall) *PendingJob {
	return &PendingJob{
		call:          call,
		stopCountdown: make(chan struct{}),
	}
}

// CountdownToDeath is a dramatic name for a simple concept - the agent has two minutes to call back
// or else it's job will be killed by the forward proxy and Alation will be sent a message as such.
//
// This grizzly fate cane be avoided by calling Rescue on the PendingJob.
func (p *PendingJob) CountdownToDeath(parent *PendingJobs) {
	stopwatch := time.NewTicker(time.Minute * 2)
	select {
	case <-p.stopCountdown:
	case <-stopwatch.C:
		parent.lock.Lock()
		defer parent.lock.Unlock()
		_, present := parent.pending[p.call.JobId]
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
		delete(parent.pending, p.call.JobId)
		p.call.Alation.SendError(status.New(codes.Unavailable, "agent failed to callback within time tolerance"))
	}
}

// Rescue stops the countdown to the deadline for this PendingJob.
func (p *PendingJob) Rescue() {
	close(p.stopCountdown)
}
