package wrapper

import (
	"../../proto"
	responseCode "../../util"
	"../KeyValueStore"
	"context"
	"fmt"
)

type Server struct{}

func (s *Server) PUT(ctx context.Context, request *proto.PutRequest) (*proto.Response, error)  {
	key := request.Key
	value := request.Value

	if key == "" {
		return &proto.Response{ResponseCode:responseCode.INVALID_INPUT, Message:"Key cannot be null or empty. PUT is not initiated"}, nil
	}

	if value == "" {
		return &proto.Response{ResponseCode:responseCode.INVALID_INPUT, Message:"Value cannot be null or empty. PUT is not initiated"}, nil
	}

	KeyValueStore.PUT(key, value)

	return &proto.Response{ResponseCode:responseCode.SUCCESS, Message:fmt.Sprintf("PUT has been successful.  Key-Value pair %s-%s is updated into the DB successfully", key, value)}, nil
}

func (s *Server) GET(ctx context.Context, request *proto.GetAndDeleteRequest) (*proto.Response, error)  {

	key := request.Key

	if key == "" {
		return &proto.Response{ResponseCode:responseCode.INVALID_INPUT, Message:"Key cannot be null or empty. GET is not initiated"}, nil
	}

	value := KeyValueStore.GET(key)

	if value == "" {
		return &proto.Response{ResponseCode:responseCode.INVALID_INPUT, Message:"Value is not found for the given key"}, nil
	}

	return &proto.Response{ResponseCode:responseCode.SUCCESS, Message:fmt.Sprintf("GET has been successful. Value for the key %s is %s", key, value)}, nil

}

func (s *Server) DELETE(ctx context.Context, request *proto.GetAndDeleteRequest) (*proto.Response, error)  {

	key := request.Key

	if key == "" {
		return &proto.Response{ResponseCode:responseCode.INVALID_INPUT, Message:"Key cannot be null or empty. DELETE is not initiated"}, nil
	}

	value := KeyValueStore.DELETE(key)

	if value == "" {
		return &proto.Response{ResponseCode:responseCode.INVALID_INPUT, Message:"Value is not found for the given key. DELETE is not initiated"}, nil
	}

	return &proto.Response{ResponseCode:responseCode.SUCCESS, Message:fmt.Sprint("DELETE has been successful. value for the key %s is %s.", key, value)}, nil
}