#!/bin/bash

# Initialize variables
FAILED_TESTS=0
LOG_FILE="/tmp/ping_test-env_log.txt"

# Clear the previous log file
> "$LOG_FILE"

# Function to execute a ping inside a container (Parallel Execution)
test_ping() {
    local container="$1"
    local target_ip="$2"
    local description="$3"
    local log_file="/tmp/ping_output_$(echo "$description" | tr -d '[:space:]/').log"

    # Run ping in background and store process ID
    docker exec -i "$container" bash -c "ping -c 4 $target_ip" > "$log_file" 2>&1 &

    # Store background process ID, log file, and description for later check
    PING_PIDS+=($!)
    PING_LOGS+=("$log_file")
    PING_SOURCE+=("$container")
    PING_DESCRIPTIONS+=("$description")
    PING_EXPECT_FAILURE+=("false") # Expect success
}

test_no_ping() {
    local container="$1"
    local target_ip="$2"
    local description="$3"
    local log_file="/tmp/ping_output_$(echo "$description" | tr -d '[:space:]/').log"

    # Run ping in background and store process ID
    docker exec -i "$container" bash -c "ping -c 4 $target_ip" > "$log_file" 2>&1 &

    # Store background process ID, log file, and description for later check
    PING_PIDS+=($!)
    PING_LOGS+=("$log_file")
    PING_SOURCE+=("$container")
    PING_DESCRIPTIONS+=("$description")
    PING_EXPECT_FAILURE+=("true") # Expect failure
}

# Arrays to store background processes, logs, and descriptions
PING_PIDS=()
PING_LOGS=()
PING_DESCRIPTIONS=()
PING_EXPECT_FAILURE=() # Track whether we expect success or failure

# Test pings from PC101
# test_ping "clab-frr01-pc101" "10.0.0.100" "PC101-PC101"
test_ping "clab-frr01-pc101" "10.1.0.100" "PC101-PC111"
test_ping "clab-frr01-pc101" "10.2.0.100" "PC101-PC121"
test_ping "clab-frr01-pc101" "10.3.0.100" "PC101-PC131"
test_ping "clab-frr01-pc101" "192.168.1.100" "PC101-PC191"
test_ping "clab-frr01-pc101" "10.20.0.100" "PC101-PC201"
test_ping "clab-frr01-pc101" "10.30.0.100" "PC101-PC301"
test_ping "clab-frr01-pc101" "192.168.4.4" "PC101-CR104"
test_ping "clab-frr01-pc101" "192.168.5.5" "PC101-CR105"
test_ping "clab-frr01-pc101" "192.168.6.6" "PC101-CR106"
test_ping "clab-frr01-pc101" "192.168.7.7" "PC101-CR107"
test_ping "clab-frr01-pc101" "192.168.8.8" "PC101-CR108"
test_ping "clab-frr01-pc101" "192.168.9.9" "PC101-CR109"
test_no_ping "clab-frr01-pc101" "192.168.4.100" "PC101-CPC104"
test_ping "clab-frr01-pc101" "192.168.5.100" "PC101-CPC105"
test_no_ping "clab-frr01-pc101" "192.168.6.100" "PC101-CPC106"
test_no_ping "clab-frr01-pc101" "192.168.7.100" "PC101-CPC107"
test_no_ping "clab-frr01-pc101" "192.168.8.100" "PC101-CPC108"
test_no_ping "clab-frr01-pc101" "192.168.9.100" "PC101-CPC109"

# Test pings from PC111
test_ping "clab-frr01-pc111" "10.0.0.100" "PC111-PC101"
# test_ping "clab-frr01-pc111" "10.1.0.100" "PC111-PC111"
test_ping "clab-frr01-pc111" "10.2.0.100" "PC111-PC121"
test_ping "clab-frr01-pc111" "10.3.0.100" "PC111-PC131"
test_ping "clab-frr01-pc111" "192.168.1.100" "PC111-PC191"
test_ping "clab-frr01-pc111" "10.20.0.100" "PC111-PC201"
test_ping "clab-frr01-pc111" "10.30.0.100" "PC111-PC301"
test_ping "clab-frr01-pc111" "192.168.4.4" "PC111-CR104"
test_ping "clab-frr01-pc111" "192.168.5.5" "PC111-CR105"
test_no_ping "clab-frr01-pc111" "192.168.4.100" "PC111-CPC104"
test_ping "clab-frr01-pc111" "192.168.5.100" "PC111-CPC105"

# Test pings from r111
# fails: because router 191 only knows the client networks, causing that ping doesn’t find return path
test_no_ping "clab-frr01-r111" "192.168.1.100" "R111-PC191"

# Test pings from PC121
test_ping "clab-frr01-pc121" "10.0.0.100" "PC121-PC101"
test_ping "clab-frr01-pc121" "10.1.0.100" "PC121-PC111"
# test_ping "clab-frr01-pc121" "10.2.0.100" "PC121-PC121"
test_ping "clab-frr01-pc121" "10.3.0.100" "PC121-PC131"
test_ping "clab-frr01-pc121" "192.168.1.100" "PC121-PC191"
test_ping "clab-frr01-pc121" "10.20.0.100" "PC121-PC201"
test_ping "clab-frr01-pc121" "10.30.0.100" "PC121-PC301"
test_ping "clab-frr01-pc121" "192.168.4.4" "PC121-CR104"
test_ping "clab-frr01-pc121" "192.168.5.5" "PC121-CR105"
test_no_ping "clab-frr01-pc121" "192.168.4.100" "PC121-CPC104"
test_ping "clab-frr01-pc121" "192.168.5.100" "PC121-CPC105"

# Test pings from PC131
test_ping "clab-frr01-pc131" "10.0.0.100" "PC131-PC101"
test_ping "clab-frr01-pc131" "10.1.0.100" "PC131-PC111"
test_ping "clab-frr01-pc131" "10.2.0.100" "PC131-PC121"
# test_ping "clab-frr01-pc131" "10.3.0.100" "PC131-PC131"
test_ping "clab-frr01-pc131" "192.168.1.100" "PC131-PC191"
test_ping "clab-frr01-pc131" "10.20.0.100" "PC131-PC201"
test_ping "clab-frr01-pc131" "10.30.0.100" "PC131-PC301"
test_ping "clab-frr01-pc131" "192.168.4.4" "PC131-CR104"
test_ping "clab-frr01-pc131" "192.168.5.5" "PC131-CR105"
test_no_ping "clab-frr01-pc131" "192.168.4.100" "PC131-CPC104"
test_ping "clab-frr01-pc131" "192.168.5.100" "PC131-CPC105"

# Test pings from PC191
test_ping "clab-frr01-pc191" "10.0.0.100" "PC191-PC101"
test_ping "clab-frr01-pc191" "10.1.0.100" "PC191-PC111"
test_ping "clab-frr01-pc191" "10.2.0.100" "PC191-PC121"
test_ping "clab-frr01-pc191" "10.3.0.100" "PC191-PC131"
# test_ping "clab-frr01-pc191" "192.168.1.100" "PC191-PC191"
test_ping "clab-frr01-pc191" "10.20.0.100" "PC191-PC201"
test_no_ping "clab-frr01-pc191" "10.30.0.100" "PC191-PC301"
# the next four test have to fail because no matching route exists on R191
test_no_ping "clab-frr01-pc191" "192.168.4.4" "PC191-CR104"
test_no_ping "clab-frr01-pc191" "192.168.5.5" "PC191-CR105"
test_no_ping "clab-frr01-pc191" "192.168.4.100" "PC191-CPC104"
test_no_ping "clab-frr01-pc191" "192.168.5.100" "PC191-CPC105"

# Test pings from PC201
test_ping "clab-frr01-pc201" "10.0.0.100" "PC201-PC101"
test_ping "clab-frr01-pc201" "10.1.0.100" "PC201-PC111"
test_ping "clab-frr01-pc201" "10.2.0.100" "PC201-PC121"
test_ping "clab-frr01-pc201" "10.3.0.100" "PC201-PC131"
test_ping "clab-frr01-pc201" "192.168.1.100" "PC201-PC191"
# test_ping "clab-frr01-pc201" "10.20.0.100" "PC201-PC201"
test_no_ping "clab-frr01-pc201" "10.30.0.100" "PC201-PC301"
test_ping "clab-frr01-pc201" "192.168.4.4" "PC201-CR104"
test_ping "clab-frr01-pc201" "192.168.5.5" "PC201-CR105"
test_no_ping "clab-frr01-pc201" "192.168.4.100" "PC201-CPC104"
test_ping "clab-frr01-pc201" "192.168.5.100" "PC201-CPC105"

# Test pings from PC301
test_ping "clab-frr01-pc301" "10.0.0.100" "PC301-PC101"
test_ping "clab-frr01-pc301" "10.1.0.100" "PC301-PC111"
test_ping "clab-frr01-pc301" "10.2.0.100" "PC301-PC121"
test_ping "clab-frr01-pc301" "10.3.0.100" "PC301-PC131"
test_no_ping "clab-frr01-pc301" "192.168.1.100" "PC301-PC191"
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
# only this test will work because R391 has only one static route (this)
test_ping "clab-frr01-pc391" "10.30.0.100" "PC391-PC301"

# Wait for all background pings to finish
wait "${PING_PIDS[@]}"

echo "✅ All ping tests completed!"

# Now check each log file for success/failure
source_device="none"
for i in "${!PING_LOGS[@]}"; do
    log_file="${PING_LOGS[i]}"
    description="${PING_DESCRIPTIONS[i]}"
    expect_failure="${PING_EXPECT_FAILURE[i]}"
    current_source_device="${PING_SOURCE[i]}"

    # Check if the log file is empty
    if [ ! -s "$log_file" ]; then
        echo "⚠️ WARNING: Log file for $description is empty! Something went wrong!"
        echo "⚠️ WARNING: Log file for $description is empty!" >> "$LOG_FILE"
        FAILED_TESTS=1  # Mark failure
        continue
    fi

    # Print source device header if it changes
    if [ "$current_source_device" != "$source_device" ]; then
        echo -e "\nPings from $current_source_device"
        source_device="$current_source_device"
    fi

    # Check for success or failure
    if [ "$expect_failure" == "false" ]; then
        # Expecting success: Check for "0% packet loss"
        if grep -q " 0% packet loss" "$log_file"; then
            echo "✅ Success: Ping from $description passed!"
        else
            echo "❌ FAILURE: Ping from $description failed!"
            echo "❌ FAILURE: Ping from $description failed!" >> "$LOG_FILE"
            cat "$log_file" >> "$LOG_FILE"
            FAILED_TESTS=1  # Mark that at least one test has failed
        fi
    elif [ "$expect_failure" == "true" ]; then
        # Expecting failure: Check for "100% packet loss"
        if grep -q "100% packet loss" "$log_file"; then
            echo "✅ Success: Ping from $description failed (as expected)!"
        else
            echo "❌ FAILURE: Ping from $description succeeded but should have failed!"
            echo "❌ FAILURE: Ping from $description succeeded but should have failed!" >> "$LOG_FILE"
            cat "$log_file" >> "$LOG_FILE"
            FAILED_TESTS=1  # Mark that at least one test has failed
        fi
    else
        echo "⚠️ WARNING: unknown error"
    fi
done


# Final Check
if [ "$FAILED_TESTS" -ne 0 ]; then
    echo "⚠️ Some ping tests failed or unexpected successes occurred. Check the log file: $LOG_FILE"
    exit 1
else
    echo "🎉 All connectivity tests passed successfully!"
    exit 0
fi
