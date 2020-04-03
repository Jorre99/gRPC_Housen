# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import HouseServer_pb2 as HouseServer__pb2


class BroadcastStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.CreateStream = channel.unary_stream(
        '/proto.Broadcast/CreateStream',
        request_serializer=HouseServer__pb2.Connect.SerializeToString,
        response_deserializer=HouseServer__pb2.Message.FromString,
        )
    self.BroadcastMessage = channel.unary_unary(
        '/proto.Broadcast/BroadcastMessage',
        request_serializer=HouseServer__pb2.Message.SerializeToString,
        response_deserializer=HouseServer__pb2.Close.FromString,
        )


class BroadcastServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def CreateStream(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def BroadcastMessage(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_BroadcastServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'CreateStream': grpc.unary_stream_rpc_method_handler(
          servicer.CreateStream,
          request_deserializer=HouseServer__pb2.Connect.FromString,
          response_serializer=HouseServer__pb2.Message.SerializeToString,
      ),
      'BroadcastMessage': grpc.unary_unary_rpc_method_handler(
          servicer.BroadcastMessage,
          request_deserializer=HouseServer__pb2.Message.FromString,
          response_serializer=HouseServer__pb2.Close.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'proto.Broadcast', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
