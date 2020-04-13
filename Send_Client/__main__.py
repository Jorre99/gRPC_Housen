from __future__ import print_function
import logging

import grpc
from server_fllower_house.proto import HouseServer_pb2_grpc
from server_fllower_house.proto import HouseServer_pb2


def run():
    with grpc.insecure_channel('localhost:5050') as channel:
        stub = HouseServer_pb2_grpc.BroadcastStub(channel)
        response1 = stub.BroadcastMessage(HouseServer_pb2.Message(id=input("uw naam: ")))
        while True:
            response = stub.BroadcastMessage(HouseServer_pb2.Message(content=input("type your message: ")))


if __name__ == '__main__':
    logging.basicConfig()
    run()