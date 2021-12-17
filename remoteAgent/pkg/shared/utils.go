package shared

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const ConnectorIdHeader = "x-alation-connector"
const AgentIdHeader = "x-alation-agent"
const JobIdHeader = "x-alation-job"

func ExtractConnectorId(ctx context.Context) (uint64, error) {
	return ExtractUint64Header(ConnectorIdHeader, ctx)
}

func ExtractAgentId(ctx context.Context) (uint64, error) {
	return ExtractUint64Header(AgentIdHeader, ctx)
}

func ExtractJobId(ctx context.Context) (uint64, error) {
	return ExtractUint64Header(JobIdHeader, ctx)
}

func ExtractUint64Header(header string, ctx context.Context) (uint64, error) {
	header, err := ExtractHeader(header, ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(header, 10, 64)
}

func ExtractHeader(header string, ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("expected the header `%s`, however no headers were present", header)
	}
	values := md.Get(header)
	if len(values) == 0 {
		return "", fmt.Errorf("expected the header `%s`, however it was not present", header)
	}
	if len(values) > 1 {
		// @TODO
		//logrus.Warnf("expected a single value for the header '%s', however the following were found: %v", header, values)
	}
	return values[0], nil
}

type HeaderBuilder struct {
	Headers map[string]string
}

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{Headers: map[string]string{}}
}

func (h *HeaderBuilder) Build(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.New(h.Headers))
}

func (h *HeaderBuilder) SetConnectorId(value uint64) *HeaderBuilder {
	return h.SetUint64Header(ConnectorIdHeader, value)
}

func (h *HeaderBuilder) SetAgentId(value uint64) *HeaderBuilder {
	return h.SetUint64Header(AgentIdHeader, value)
}

func (h *HeaderBuilder) SetJobId(value uint64) *HeaderBuilder {
	return h.SetUint64Header(JobIdHeader, value)
}

func (h *HeaderBuilder) SetUint64Header(header string, value uint64) *HeaderBuilder {
	return h.SetHeader(header, strconv.FormatUint(value, 10))
}

func (h *HeaderBuilder) SetHeader(header, value string) *HeaderBuilder {
	h.Headers[header] = value
	return h
}
