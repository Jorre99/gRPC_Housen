from __future__ import print_function
import logging

import grpc
import HouseServer_pb2
import HouseServer_pb2_grpc


def run():
    with grpc.insecure_channel('localhost:5050') as channel:
        stub = HouseServer_pb2_grpc.BroadcastStub(channel)
        while True:
            response = stub.BroadcastMessage(HouseServer_pb2.Message(content=input("type your message: ")))


if __name__ == '__main__':
    logging.basicConfig()
    run()