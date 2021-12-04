import grpc

from test_pb2 import *
from test_pb2_grpc import *
from pstream import Stream

PORT = "8080"
channel = grpc.insecure_channel(f"0.0.0.0:{PORT}", options=[('grpc.keepalive_time_ms', 30000)])
stub = TestStub(channel)
metadata=(("x-alation-connector", "1"),)

def unary():
    response = stub.Unary(request=String(message="hi from unary"), metadata=metadata)
    print(f"unary: {response}")

def server_stream():
    stream = stub.ServerStream(request=String(message="hi from server stream"), metadata=metadata)
    for i, object in enumerate(stream):
        print(f"Server stream {i} {object}")

def client_stream():
    iterator = Stream([String(message="1"), String(message="2")])
    stub.ClientStream(request_iterator=iterator, metadata=metadata)
    print("client stream done")

def bidirectional_stream():
    iterator = Stream([String(message="1"), String(message="2")])
    stream = stub.Bidirectional(request_iterator=iterator, metadata=metadata)
    first = stream.next()
    print(f"bidi 1 {first}")
    second = stream.next()
    print(f"bidi 1 {second}")

def intentional_error():
    error = None
    try:
        for _ in stub.IntentionalError(request=String(message="explode for me"), metadata=metadata):
            pass
    except Exception as e:
        error = e
    if error is None:
        raise Exception("wanted exception")
    print(f"correctly got exception {error}")


if __name__ == "__main__":
    unary()
    server_stream()
    client_stream()
    bidirectional_stream()
    intentional_error()
