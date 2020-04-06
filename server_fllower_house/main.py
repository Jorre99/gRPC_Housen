from concurrent import futures
import grpc
import HouseServer_pb2
import HouseServer_pb2_grpc
import time
import threading


class Listener(HouseServer_pb2_grpc.BroadcastServicer):
    def __init__(self):
        self.counter = 0
        self.clients = []
        self.client_lock = threading.Lock()
        self.lastPrintTime = time.time()

    def CreateStream(self, request, context):
        with self.client_lock:
            print('new client')
            new_client = Connection(self)
            self.clients.append(new_client)
            return new_client

    def BroadcastMessage(self, request, context):
        with self.client_lock:
            print(f'broadcasting to {len(self.clients)} clients')
            for client in self.clients:
                client.sendmsg(request)
            print(request)
            return HouseServer_pb2.Close()

    def delete(self, connection):
        with self.client_lock:
            print(f'removing {connection}')
            self.clients.remove(connection)


class Connection:
    def __init__(self, listener):
        self.cond = threading.Condition()
        self.messages = []
        self.listener = listener
        self.timer = threading.Timer(15, self.quit)
        self.timer.start()
        self.quitting = False

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
        self.listener.delete(self)
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
