#!/bin/bash
# TARGZ_FILES=$(ls ./data/fb_used/*.tar.gz)
# TARGZ_TO="./data/fb_used"

POWER_FILES=$(ls /data/power/*.tar.gz)
POWER_TO="/data/Trace_code/lj/GPU_data/data/power"

MEM_FILES=$(ls /data/fb_used/*.tar.gz)
MEM_TO="/data/Trace_code/lj/GPU_data/data/fb_used"

UTIL_FILES=$(ls /data/gpu_utilization/*.tar.gz)
UTIL_TO="/data/Trace_code/lj/GPU_data/data/gpu_utilization"

for power_file in $POWER_FILES; do	
	tar -zxvf $POWER_file -C $POWER_TO 
	# rm -rf $power_file
done

for mem_file in $MEM_FILES; do	
	tar -zxvf $mem_file -C $MEM_TO 	
done

for util_file in $UTIL_FILES; do	
	tar -zxvf $util_file -C $UTIL_TO 	
done