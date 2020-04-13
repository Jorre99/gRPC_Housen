from concurrent import futures
import grpc
from server_fllower_house.proto import HouseServer_pb2_grpc
from server_fllower_house.proto import HouseServer_pb2
import time
import threading


# run from gRPC_Housen
# generate proto: python3 -m grpc_tools.protoc -I . --python_out=. --grpc_python_out=. server_fllower_house/proto/HouseServer.proto


class Listener(HouseServer_pb2_grpc.BroadcastServicer):
    def __init__(self):
        self.counter = 0
        self.clients = {}
        self.client_lock = threading.Lock()
        self.lastPrintTime = time.time()

    def CreateStream(self, connect, context):
        with self.client_lock:
            print('new client', connect.id)
            if connect.id in self.clients:
                return []
            new_client = Connection(self, connect.id)
            self.clients[connect.id] = new_client
            return new_client

    def BroadcastMessage(self, message, context):
        with self.client_lock:
            print('sending to {}'.format(message.peer_user))
            if message.peer_user in self.clients:
                self.clients[message.peer_user].sendmsg(message)
                return HouseServer_pb2.Close()
            for client in self.clients.values():
                client.sendmsg(message)
            print(message)
            return HouseServer_pb2.Close()

    def delete(self, connection_id):
        with self.client_lock:
            print('removing {}'.format(connection_id))
            del self.clients[connection_id]
            # self.clients.remove(connection)


class Connection:
    def __init__(self, listener, id):
        self.cond = threading.Condition()
        self.messages = []
        self.listener = listener
        self.timer = threading.Timer(15, self.quit)
        self.timer.start()
        self.quitting = False
        self.id = id

    def __iter__(self):
        return self

    def __next__(self):
        with self.cond:
            self.cond.wait_for(lambda: self.messages or self.quitting)
            if self.quitting:
                raise StopIteration()
            self.timer.cancel()
            self.timer = threading.Timer(15, self.quit)
            self.timer.start()
            return self.messages.pop(0)

    def sendmsg(self, message):
        with self.cond:
            self.messages.append(message)
            self.cond.notify()

    def quit(self):
        self.listener.delete(self.id)
        with self.cond:
            self.quitting = True
            self.cond.notify()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=20))
    HouseServer_pb2_grpc.add_BroadcastServicer_to_server(Listener(), server)
    server.add_insecure_port('[::]:5050')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
