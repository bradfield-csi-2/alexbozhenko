import argparse
import http.server
import json
from socketserver import ThreadingMixIn


class JSONHeaderReporter(http.server.BaseHTTPRequestHandler):
    """
    A simple HTTP server which simply returns a JSON representation
    of the request headers.
    """
    def do_GET(self):
        self.send_response(http.server.HTTPStatus.OK)
        body = json.dumps(dict(self.headers), indent=4).encode('utf8')
        self.send_header('Content-Length', str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    do_POST = do_GET


class JSONHeaderReporterKeepalive(JSONHeaderReporter):
    protocol_version = 'HTTP/1.1'


class ThreadingHTTPServer(ThreadingMixIn, http.server.HTTPServer):
    pass


def run_server(host, port, keepalive):
    handler = JSONHeaderReporterKeepalive if keepalive else JSONHeaderReporter
    httpd = ThreadingHTTPServer((host, port), handler)
    version = '1.1' if keepalive else '1.0'
    print(f'Running on {host}:{port}, with HTTP/{version}')
    httpd.serve_forever()


if __name__ == '__main__':
    parser = argparse.ArgumentParser('python3 server.py')
    parser.add_argument('--host', default='localhost',
                        help='Local interface, default "localhost"')
    parser.add_argument('--port', default='9000',
                        help='Port to listen on, default 9000')
    parser.add_argument('--keepalive', default=False,
                        help='Run HTTP/1.1, so support persistent connections')
    args = parser.parse_args()
    run_server(args.host, int(args.port), args.keepalive)
