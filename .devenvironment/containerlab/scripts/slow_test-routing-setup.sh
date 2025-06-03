#!/bin/bash

# Initialize variables
FAILED_TESTS=0
LOG_FILE="/tmp/ping_test-env_log.txt"

# Clear the previous log file
> "$LOG_FILE"

# Function to execute a ping inside a container
test_ping() {
    local container="$1"
    local target_ip="$2"
    local description="$3"

    echo "🔄 Testing connectivity $description ($target_ip)..."

    docker exec -ti "$container" bash -c "ping -c 2 $target_ip" > /tmp/ping_output.log 2>&1

    cat "/tmp/ping_output.log"

    if grep -q " 0% packet loss" /tmp/ping_output.log; then
        echo "✅ Success: Ping from $description ($target_ip) passed!"
    else
        echo "❌ FAILURE: Ping from $description ($target_ip) failed!"
        echo "❌ FAILURE: Ping from $description ($target_ip) failed!" >> "$LOG_FILE"
        cat /tmp/ping_output.log >> "$LOG_FILE"
        FAILED_TESTS=1  # Mark that at least one test has failed
    fi
}

test_no_ping() {
    local container="$1"
    local target_ip="$2"
    local description="$3"

    echo "🔄 Ensuring no active connection from $description ($target_ip)..."

    docker exec -ti "$container" bash -c "ping -c 2 $target_ip" > /tmp/ping_output.log 2>&1

    if grep -q "100% packet loss" /tmp/ping_output.log; then
        echo "✅ Success: Ping from $description ($target_ip) failed (as intended)!"
    else
        echo "❌ FAILURE: Ping from $description ($target_ip) passed but it should fail!"
        echo "❌ FAILURE: Ping from $description ($target_ip) passed but it should fail!" >> "$LOG_FILE"
        cat /tmp/ping_output.log >> "$LOG_FILE"
        FAILED_TESTS=1  # Mark that at least one test has failed
    fi
}

# Test pings from PC101
# test_ping "clab-frr01-pc101" "10.0.0.100" "PC101-PC101"
test_ping "clab-frr01-pc101" "10.1.0.100" "PC101-PC111"
test_ping "clab-frr01-pc101" "10.2.0.100" "PC101-PC121"
test_ping "clab-frr01-pc101" "10.3.0.100" "PC101-PC131"
test_ping "clab-frr01-pc101" "192.168.1.100" "PC101-PC191"
test_ping "clab-frr01-pc101" "10.20.0.100" "PC101-PC201"
test_ping "clab-frr01-pc101" "10.30.0.100" "PC101-PC301"

# Test pings from PC111
test_ping "clab-frr01-pc111" "10.0.0.100" "PC111-PC101"
# test_ping "clab-frr01-pc111" "10.1.0.100" "PC111-PC111"
test_ping "clab-frr01-pc111" "10.2.0.100" "PC111-PC121"
test_ping "clab-frr01-pc111" "10.3.0.100" "PC111-PC131"
test_no_ping "clab-frr01-pc111" "192.168.1.100" "PC111-PC121"
test_no_ping "clab-frr01-pc111" "10.20.0.100" "PC111-PC201"
test_ping "clab-frr01-pc111" "10.30.0.100" "PC111-PC301"

# Test pings from PC121
test_ping "clab-frr01-pc121" "10.0.0.100" "PC121-PC101"
test_ping "clab-frr01-pc121" "10.1.0.100" "PC121-PC111"
# test_ping "clab-frr01-pc121" "10.2.0.100" "PC121-PC121"
test_ping "clab-frr01-pc121" "10.3.0.100" "PC121-PC131"
test_ping "clab-frr01-pc121" "192.168.1.100" "PC121-PC191"
test_ping "clab-frr01-pc121" "10.20.0.100" "PC121-PC201"
test_ping "clab-frr01-pc121" "10.30.0.100" "PC121-PC301"

# Test pings from PC131
test_ping "clab-frr01-pc131" "10.0.0.100" "PC131-PC101"
test_ping "clab-frr01-pc131" "10.1.0.100" "PC131-PC111"
test_ping "clab-frr01-pc131" "10.2.0.100" "PC131-PC121"
# test_ping "clab-frr01-pc131" "10.3.0.100" "PC131-PC131"
test_no_ping "clab-frr01-pc131" "192.168.1.100" "PC131-PC191"
test_no_ping "clab-frr01-pc131" "10.20.0.100" "PC131-PC201"
test_no_ping "clab-frr01-pc131" "10.30.0.100" "PC131-PC301"

# Test pings from PC191
test_ping "clab-frr01-pc191" "10.0.0.100" "PC191-PC101"
test_no_ping "clab-frr01-pc191" "10.1.0.100" "PC191-PC111"
test_ping "clab-frr01-pc191" "10.2.0.100" "PC191-PC121"
test_no_ping "clab-frr01-pc191" "10.3.0.100" "PC191-PC131"
# test_ping "clab-frr01-pc191" "192.168.1.100" "PC191-PC191"
test_no_ping "clab-frr01-pc191" "10.20.0.100" "PC191-PC201"
test_no_ping "clab-frr01-pc191" "10.30.0.100" "PC191-PC301"

# Test pings from PC201
test_ping "clab-frr01-pc201" "10.0.0.100" "PC201-PC101"
test_no_ping "clab-frr01-pc201" "10.1.0.100" "PC201-PC111"
test_ping "clab-frr01-pc201" "10.2.0.100" "PC201-PC121"
test_no_ping "clab-frr01-pc201" "10.3.0.100" "PC201-PC131"
test_ping "clab-frr01-pc201" "192.168.1.100" "PC201-PC191"
# test_ping "clab-frr01-pc201" "10.20.0.100" "PC201-PC201"
test_no_ping "clab-frr01-pc201" "10.30.0.100" "PC201-PC301"

# Test pings from PC301
test_ping "clab-frr01-pc301" "10.0.0.100" "PC301-PC101"
test_ping "clab-frr01-pc301" "10.1.0.100" "PC301-PC111"
test_ping "clab-frr01-pc301" "10.2.0.100" "PC301-PC121"
test_no_ping "clab-frr01-pc301" "10.3.0.100" "PC301-PC131"
test_ping "clab-frr01-pc301" "192.168.30.100" "PC301-PC391"
test_no_ping "clab-frr01-pc301" "10.20.0.100" "PC301-PC201"
# test_ping "clab-frr01-pc301" "10.30.0.100" "PC301-PC301"

# Test pings from PC391
test_no_ping "clab-frr01-pc391" "10.0.0.100" "PC391-PC101"
test_no_ping "clab-frr01-pc391" "10.1.0.100" "PC391-PC111"
test_no_ping "clab-frr01-pc391" "10.2.0.100" "PC391-PC121"
test_no_ping "clab-frr01-pc391" "10.3.0.100" "PC391-PC131"
# test_ping "clab-frr01-pc391" "192.168.30.100" "PC391-PC391"
test_no_ping "clab-frr01-pc391" "10.20.0.100" "PC391-PC201"
test_ping "clab-frr01-pc391" "10.30.0.100" "PC391-PC301"


# Check if any test failed and print the log
if [ "$FAILED_TESTS" -ne 0 ]; then
    echo "⚠️ Some ping tests failed. Here are the details: $LOG_FILE"
    # cat "$LOG_FILE"
    exit 1
else
    echo "🎉 All connectivity tests passed successfully!"
    exit 0
fi

