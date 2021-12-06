package grpcinverter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/metadata"
)

func ExtractConnectorId(ctx context.Context) (uint64, error) {
	return ExtractUint64Header("x-alation-connector", ctx)
}

func ExtractAgentId(ctx context.Context) (uint64, error) {
	return ExtractUint64Header("x-alation-agent", ctx)
}

func ExtractJobId(ctx context.Context) (uint64, error) {
	return ExtractUint64Header("x-alation-job", ctx)
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
		logrus.Warnf("expected a single value for the header '%s', however the following were found: %v", header, values)
	}
	return values[0], nil
}
