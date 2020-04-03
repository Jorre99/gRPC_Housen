from __future__ import print_function
import logging

import grpc
import HouseServer_pb2
import HouseServer_pb2_grpc


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel('localhost:5050') as channel:
        stub = HouseServer_pb2_grpc.BroadcastStub(channel)
        user = HouseServer_pb2.User(id='bla', name='alb')
        stream = stub.CreateStream(HouseServer_pb2.Connect(user=user, active=True))

        for resp in stream:
            print(resp)


if __name__ == '__main__':
    logging.basicConfig()
    run()
