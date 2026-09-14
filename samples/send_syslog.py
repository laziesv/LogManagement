"""Send a current RFC3164 sample over UDP or newline-framed TCP."""
import argparse
import datetime
import socket

p = argparse.ArgumentParser()
p.add_argument('--host', default='127.0.0.1')
p.add_argument('--port', type=int, default=5514)
p.add_argument('--tcp', action='store_true')
args = p.parse_args()
now = datetime.datetime.now(datetime.timezone.utc)
months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
stamp = f'{months[now.month-1]} {now.day:2d} {now:%H:%M:%S}'
message = f'<134>{stamp} fw01 vendor=demo product=ngfw action=deny src=10.0.1.10 dst=8.8.8.8 spt=5353 dpt=53 proto=udp msg="DNS blocked"\n'.encode()
with socket.socket(socket.AF_INET, socket.SOCK_STREAM if args.tcp else socket.SOCK_DGRAM) as s:
    s.settimeout(5)
    if args.tcp:
        s.connect((args.host, args.port))
        s.sendall(message)
    else:
        s.sendto(message, (args.host, args.port))
print('Syslog sent. Refresh Log explorer to confirm ingestion.')
