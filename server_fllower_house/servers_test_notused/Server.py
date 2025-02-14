import socket
import pickle

HEADER_SIZE = 10


s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind((socket.gethostname(), 1235))
s.listen(5)

while True:
    client_socket, address = s.accept()
    print(f'Connection from {address} has been established')

    d = {1: 'lel', 2: 'lit'}
    msg = pickle.dumps(d)

    msg = bytes(f'{len(msg):<{HEADER_SIZE}}', 'utf-8') + msg

    client_socket.send(msg)

