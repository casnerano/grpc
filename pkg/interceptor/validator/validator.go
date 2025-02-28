package validator

import (
    "context"
    "google.golang.org/grpc"
)

// UnaryServerInterceptor проверка входных запросов сервера
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        if err := validation(req); err != nil {
            return nil, err
        }
        return handler(ctx, req)
    }
}

// UnaryClientInterceptor проверка входных запросов клиента
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        if err := validation(req); err != nil {
            return err
        }
        return invoker(ctx, method, req, reply, cc, opts...)
    }
}
