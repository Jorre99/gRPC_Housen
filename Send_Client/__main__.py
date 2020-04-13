from __future__ import print_function
import logging

import grpc
from server_fllower_house.proto import HouseServer_pb2_grpc
from server_fllower_house.proto import HouseServer_pb2


def run():
    with grpc.insecure_channel('bink.us-east.chatter.housen.tech:5050') as channel:
        stub = HouseServer_pb2_grpc.BroadcastStub(channel)
        id=input("uw naam: ")
        peer_user=input("na wie wil je een berricht sturen? ")
        while True:
            response = stub.BroadcastMessage(HouseServer_pb2.Message(content=input("type your message: "),
                                                                     id=id, peer_user=peer_user))


if __name__ == '__main__':
    logging.basicConfig()
    run()