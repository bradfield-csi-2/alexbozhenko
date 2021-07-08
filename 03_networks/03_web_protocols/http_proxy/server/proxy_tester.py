import argparse
import json
import signal
import socket
import sys
import unittest
import urllib.request


class TestTimeout(Exception):
    pass


class timeout_after:
    """
    A context manager for timing out steps in a test case
    """
    def __init__(self, seconds):
        self.seconds = seconds

    def __enter__(self):
        signal.signal(signal.SIGALRM, self.handle_timeout)
        signal.alarm(self.seconds)

    def __exit__(self, *args, **kwargs):
        signal.alarm(0)

    def handle_timeout(self, signum, frame):
        raise TestTimeout(f'Test timed out after {self.seconds} seconds')


class BaseTest(unittest.TestCase):
    PROXY_LOCATION = None  # set (host, port) before running


class HttpRequestTest(BaseTest):
    """
    Test that we can send an HTTP request, and receive a response
    from the final server (via the proxy)
    """
    def test_open_url(self):
        http_addr = 'http://' + ':'.join(self.PROXY_LOCATION)
        with urllib.request.urlopen(http_addr, timeout=5) as f:
            response = json.loads(f.read().decode('utf-8'))
            # for this test, ignore connection header
            del response['Connection']
            self.assertDictEqual(response, {
                'Accept-Encoding': 'identity',
                'Host': 'localhost:8000',
                'User-Agent': 'Python-urllib/3.6'
            })


class KeepAliveTest(BaseTest):
    """
    Test that we can send an HTTP request, and that the proxy
    sends a `Connection: Keep-Alive` header
    """
    def test_open_url(self):
        http_addr = 'http://' + ':'.join(self.PROXY_LOCATION)
        with urllib.request.urlopen(http_addr, timeout=5) as f:
            response = json.loads(f.read().decode('utf-8'))
            self.assertDictEqual(response, {
                'Accept-Encoding': 'identity',
                'Host': 'localhost:8000',
                'User-Agent': 'Python-urllib/3.6',
                'Connection': 'Keep-Alive'
            })


class ConcurrentRequestTest(BaseTest):
    """
    Test that the proxy can handle concurrent requests (stretch goal)
    """
    def test_concurrent_requests(self):
        # we'll send a message in parts, concurrently between multiple
        # connections. a correctly implemented proxy won't block
        # waiting for the first to finish
        message = (
            b'GET / HTTP/1.0\r\n',
            b'Foo: Bar\r\n',
            b'\r\n'
        )
        n_socks = 3
        sockets = [
            socket.create_connection(self.PROXY_LOCATION)
            for _ in range(n_socks)
        ]

        for part in message:
            for s in sockets:
                s.send(part)

        for s in sockets:
            data = s.recv(4096)
            self.assertTrue(data)

        for s in sockets:
            s.close()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(add_help=False)
    parser.add_argument('--host', default='localhost')
    parser.add_argument('--port', default='8000')
    options, args = parser.parse_known_args()
    BaseTest.PROXY_LOCATION = (options.host, options.port)

    unittest.main(verbosity=2, argv=sys.argv[:1] + args)
