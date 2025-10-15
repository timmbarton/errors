package errs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GetGRPCInterceptor(serviceID int) grpc.UnaryServerInterceptor {
	if serviceID < 10 || serviceID > 99 {
		panic("errs.GetGRPCInterceptor(): invalid serviceID")
	}

	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			customErr := (*Err)(nil)

			if errors.As(err, &customErr) {
				err = ToGRPC(customErr)
			} else {
				err = ToGRPC(New(ErrCodeUnknown, serviceID*10000, err.Error()).(*Err))
			}

			return resp, err
		}

		return resp, err
	}
}

type logger interface {
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

func GetLoggingInterceptor(logger logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			infoFullMethod := "unknown"
			if info != nil {
				infoFullMethod = info.FullMethod
			}

			reqJSON, jsonErr := json.MarshalIndent(req, "", "  ")
			if jsonErr != nil {
				reqJSON = []byte(fmt.Sprintf("%#+v", req))
			}

			respJSON, jsonErr := json.MarshalIndent(resp, "", "  ")
			if jsonErr != nil {
				respJSON = []byte(fmt.Sprintf("%#+v", resp))
			}

			logger.Error(
				ctx,
				infoFullMethod,
				zap.Error(err),
				zap.String("request", string(reqJSON)),
				zap.String("response", string(respJSON)),
			)
		}

		return resp, err
	}
}

func ClientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	customErr, ok := Parse(err)
	if ok {
		return customErr
	}

	return err
}
