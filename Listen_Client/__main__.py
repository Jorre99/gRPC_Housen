from __future__ import print_function
import logging

import grpc
from server_fllower_house.proto import HouseServer_pb2_grpc
from server_fllower_house.proto import HouseServer_pb2


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel('localhost:5050') as channel:
        stub = HouseServer_pb2_grpc.BroadcastStub(channel)
        stream = stub.CreateStream(HouseServer_pb2.Connect(id=input("your name: ")))

        for resp in stream:
            print(resp)


if __name__ == '__main__':
    logging.basicConfig()
    run()
