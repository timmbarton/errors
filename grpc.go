package errs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

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

			isCustomErr := errors.As(err, &customErr)
			if isCustomErr {
				err = customErr
			} else {
				err = New(ErrCodeUnknown, serviceID*10000, err.Error())
			}

			return resp, ToGRPC(customErr)
		}

		return resp, err
	}
}

func LoggingInterceptor(
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

		log.Printf(
			"error on %s:\nrequest:\n%s\nerror: %s\nresponse:\n%s\n",
			infoFullMethod,
			string(reqJSON),
			err.Error(),
			string(respJSON),
		)
	}

	return resp, err
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
