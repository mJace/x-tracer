#!/usr/bin/python
#
# cachestat     Count cache kernel function calls.
#               For Linux, uses BCC, eBPF. See .c file.
#
# USAGE: cachestat
# Taken from funccount by Brendan Gregg
# This is a rewrite of cachestat from perf to bcc
# https://github.com/brendangregg/perf-tools/blob/master/fs/cachestat
#
# Copyright (c) 2016 Allan McAleavy.
# Copyright (c) 2015 Brendan Gregg.
# Licensed under the Apache License, Version 2.0 (the "License")
#
# 09-Sep-2015   Brendan Gregg   Created this.
# 21-May-2020   Sheenam Pathak        Modified to use without curse, displaying terminal output 

from __future__ import print_function
from bcc import BPF
from time import sleep, strftime
import argparse
import signal
import re
import pwd
from sys import argv
from collections import defaultdict

# signal handler
def signal_ignore(signal, frame):
    print()

# Function to gather data from /proc/meminfo
# return dictionary for quicker lookup of both values
def get_meminfo():
    result = dict()

    for line in open('/proc/meminfo'):
        k = line.split(':', 3)
        v = k[1].split()
        result[k[0]] = int(v[0])
    return result


# arguments
parser = argparse.ArgumentParser(
    description="Count cache kernel function calls",
    formatter_class=argparse.RawDescriptionHelpFormatter)
parser.add_argument("-T", "--timestamp", action="store_true",
    help="include timestamp on output")
parser.add_argument("interval", nargs="?", default=1,
    help="output interval, in seconds")
parser.add_argument("count", nargs="?", default=-1,
    help="number of outputs")
parser.add_argument("--ebpf", action="store_true",
    help=argparse.SUPPRESS)
args = parser.parse_args()
count = int(args.count)
tstamp = args.timestamp
interval = int(args.interval)


# define BPF program
bpf_text = """
#include <uapi/linux/ptrace.h>
struct key_t {
    u64 ip;
    u32 pid;
    u32 uid;
    char comm[16];
};

BPF_HASH(counts, struct key_t);

int do_count(struct pt_regs *ctx) {
    struct key_t key = {};
    u64 pid = bpf_get_current_pid_tgid();
    u32 uid = bpf_get_current_uid_gid();

    key.ip = PT_REGS_IP(ctx);
    key.pid = pid & 0xFFFFFFFF;
    key.uid = uid & 0xFFFFFFFF;
    bpf_get_current_comm(&(key.comm), 16);
    counts.increment(key); // update counter
    return 0;
}

"""


def get_processes_stats(bpf):
        counts = bpf.get_table("counts")
	stats = defaultdict(lambda: defaultdict(int))
        
	for k, v in counts.items():
                  
	   stats["%d-%d-%s" % (k.pid, k.uid, k.comm.decode('utf-8', 'replace'))][k.ip] = v.value
	stats_list = []

	for pid, count in sorted(stats.items(), key=lambda stat: stat[0]):
	    rtaccess = 0
	    wtaccess = 0
	    mpa = 0
	    mbd = 0
	    apcl = 0
	    apd = 0
	    access = 0
	    misses = 0
	    rhits = 0
	    whits = 0
	    
	    for k, v in count.items():
		if re.match(b'mark_page_accessed', bpf.ksym(k)) is not None:
		    mpa = max(0, v)
		if re.match(b'mark_buffer_dirty', bpf.ksym(k)) is not None:
		    mbd = max(0, v)

		if re.match(b'add_to_page_cache_lru', bpf.ksym(k)) is not None:
		    apcl = max(0, v)

		if re.match(b'account_page_dirtied', bpf.ksym(k)) is not None:
		    apd = max(0, v)
		  # access = total cache access incl. reads(mpa) and writes(mbd)
		    # misses = total of add to lru which we do when we write(mbd)
		    # and also the mark the page dirty(same as mbd)
		access = (mpa + mbd)
		misses = (apcl + apd)

		    # rtaccess is the read hit % during the sample period.
		    # wtaccess is the write hit % during the smaple period.
		if mpa > 0:
		    rtaccess = float(mpa) / (access + misses)
		if apcl > 0:
		    wtaccess = float(apcl) / (access + misses)

		if wtaccess != 0:
		    whits = 100 * wtaccess
		if rtaccess != 0:
		    rhits = 100 * rtaccess
	    _pid, uid, comm = pid.split('-', 2)
	    stats_list.append(
            (int(_pid), uid, comm,
            access, misses, mbd,
            rhits, whits))
        stats_list = sorted(
        stats_list, key=lambda stat: stat[3])

	counts.clear()
	return stats_list





#if debug or args.ebpf:
#    print(bpf_text)
#    if args.ebpf:
#        exit()


# load BPF program
b = BPF(text=bpf_text)
b.attach_kprobe(event="add_to_page_cache_lru", fn_name="do_count")
b.attach_kprobe(event="mark_page_accessed", fn_name="do_count")
b.attach_kprobe(event="account_page_dirtied", fn_name="do_count")
b.attach_kprobe(event="mark_buffer_dirty", fn_name="do_count")


# header
if tstamp:
    print("%-8s " % "-TIME-", end="")
print("%6s %16s %16s %8s %8s %10s %14s %14s" %
     ("PID","UID","CMD","HITS", "MISSES", "DIRTIES", "READ_HIT%", "WRITE_HIT%"))


loop = 0
exiting = 0
while 1:
    if count > 0:
        loop += 1
        if loop > count:
            exit()

    try:
        sleep(interval)
    except KeyboardInterrupt:
        exiting = 1
    # as cleanup can take many seconds, trap Ctrl-C:
        signal.signal(signal.SIGINT, signal_ignore)

# Get memory info
    mem = get_meminfo()
    cached = int(mem["Cached"]) / 1024
    buff = int(mem["Buffers"]) / 1024


    process_stats = get_processes_stats(b)
    for i, stat in enumerate(process_stats):
		    uid = int(stat[1])
		    try: 
		        username = pwd.getpwuid(uid)[0]
		    except KeyError:
		        # `pwd` throws a KeyError if the user cannot be found. This can
		        # happen e.g. when the process is running in a cgroup that has
		        # different users from the host.
		        username = 'UNKNOWN({})'.format(uid)
		    print("%6s %16s %16s %8s %8s %8s %12.0f%% %10.0f%%" %
		    (stat[0], username, stat[2], stat[3], stat[4], stat[5], stat[6], stat[7]))

	#if tstamp:
	 #   print("%-8s--- " % strftime("%H:%M:%S"), end="")
	#print("%8d %8d %8d %8d %8d %8d %12.0f %10.0f" %
	#        (pid, uid, cmd, hits, misses, dirties, whits, rhits))

    rtaccess = wtaccess = mpa = mbd = apcl = apd = access = misses = rhits = whits = 0

    if exiting:
        print("Detaching...")
        exit()
