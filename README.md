[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 

[![Build Status](https://travis-ci.com/Sheenam3/x-tracer.svg?branch=master)](https://travis-ci.com/Sheenam3/x-tracer)

# x-tracer
Hi
Server Streaming the filtered PID information from testpod's as follows in Real Time:

<pre>
Choose pod : 2
---------------------------------------------
The pod you chose is testpod
Container ID is ...
19fb910a711f5eabf2cad6569a01db3702752e2f8155059654130711b8bf2c8f
2020/04/01 12:45:25 Hostname :  dad
Start Agent Pod
Start Agent Service

{Probe:tcptracer |Sys_Time: 04:03:39 |T: 25.724 | PID:20656 | PNAME:iperf3 |IP->4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | SPORT:42334 | DPORT:6001 
{Probe:tcpconnect |Sys_Time: 04:03:40 |T: 28.857 | PID:20656 | PNAME:iperf3 | IP:4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | DPORT:6001 
{Probe:tcptracer |Sys_Time: 04:03:39 |T: 25.724 | PID:8592 | PNAME:iperf3 |IP->6 | SADDR:[::] | DADDR:[0:ffff:7f00:1::] | SPORT:0 | DPORT:65535 
{Probe:tcpaccept |Sys_Time: 04:03:40 |T: 28.863 | PID:8592 | PNAME:iperf3 | IP:6 | RADDR:::ffff:127.0.0.1 | RPORT:42336 | LADDR:::ffff:127.0.0.1 | LPORT:6001 
{Probe:tcptracer |Sys_Time: 04:03:40 |T: 25.767 | PID:20656 | PNAME:iperf3 |IP->4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | SPORT:42336 | DPORT:6001 
</pre>
