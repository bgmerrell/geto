#!/usr/bin/env python

import argparse
import sys

from jsonclient import JSONClient


def test_echo(client, args):
    print "Sending: {}".format(args.string)
    print "Received: {}".format(
        client.call("GetoRPC.Echo", args.string))


def test_host_connection(args):
    # TODO: Implement
    pass


def main():
    parser = argparse.ArgumentParser(
        description='Test geto RPCs.',
        epilog='Type \'%(prog)s <command> --help\' for help on a specific '
               'command.')
    subparsers = parser.add_subparsers()
    parser_echo = subparsers.add_parser(
            'echo', help='Test the echo RPC.')
    parser_echo.set_defaults(func=test_echo)
    parser_con_test = subparsers.add_parser(
            'connection-test', help='Call the TestHostConnection RPC.')
    parser_con_test.set_defaults(func=test_host_connection)

    parser.add_argument(
            '--server', '-s', metavar='ADDR', type=str, default='localhost',
            help="The RPC server (hostname, IP, FQDN, etc)")
    parser.add_argument(
            '--port', '-p', metavar='PORT', type=int, default=11102,
            help="The RPC server port")

    parser_echo.add_argument(
            '--string', metavar='STRING', type=str,
            help='The string to echo', default="test")

    args = parser.parse_args()

    client = JSONClient((args.server, args.port))
    args.func(client, args)

    return 0


if __name__ == '__main__':
    sys.exit(main())
