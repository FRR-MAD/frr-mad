#!/bin/bash

# Initialize variables
FAILED_TESTS=0
NOT_100PERCENT_SUCCESS=0
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

# vvvvvvvvvvvvvvvvvvvvvvvvvv PING TESTS vvvvvvvvvvvvvvvvvvvvvvvvvv #

# ========================== #
# Test pings from PC101      #
# ========================== #
# test_ping "clab-frr01-pc101" "10.0.0.100" "PC101-PC101"
test_ping "clab-frr01-pc101" "10.0.2.100" "PC101-PC102"
test_ping "clab-frr01-pc101" "10.1.0.100" "PC101-PC111"
test_ping "clab-frr01-pc101" "10.2.0.100" "PC101-PC121"
test_ping "clab-frr01-pc101" "10.3.0.100" "PC101-PC131"
test_ping "clab-frr01-pc101" "192.168.1.100" "PC101-PC191"
test_ping "clab-frr01-pc101" "10.20.0.100" "PC101-PC201"
test_ping "clab-frr01-pc101" "10.20.3.100" "PC101-PC203"
test_ping "clab-frr01-pc101" "10.20.4.100" "PC101-PC204"
test_ping "clab-frr01-pc101" "10.30.0.100" "PC101-PC301"
test_ping "clab-frr01-pc101" "192.168.4.4" "PC101-CR104"
test_ping "clab-frr01-pc101" "192.168.5.5" "PC101-CR105"
test_ping "clab-frr01-pc101" "192.168.6.6" "PC101-CR106"
test_ping "clab-frr01-pc101" "192.168.7.7" "PC101-CR107"
test_ping "clab-frr01-pc101" "192.168.8.8" "PC101-CR108"
test_ping "clab-frr01-pc101" "192.168.9.9" "PC101-CR109"
# does not work because it doesn’t find the way back,
# CPC104 has no route back to ospf network
test_no_ping "clab-frr01-pc101" "192.168.4.100" "PC101-CPC104"
# works because this PC has a default route to 192.168.5.1,
# therefore only CPC105 finds the way back.
test_ping "clab-frr01-pc101" "192.168.5.100" "PC101-CPC105"
test_no_ping "clab-frr01-pc101" "192.168.6.100" "PC101-CPC106"
test_no_ping "clab-frr01-pc101" "192.168.7.100" "PC101-CPC107"
test_no_ping "clab-frr01-pc101" "192.168.8.100" "PC101-CPC108"
test_no_ping "clab-frr01-pc101" "192.168.9.100" "PC101-CPC109"

# ========================== #
# Test pings from PC102      #
# ========================== #
test_ping "clab-frr01-pc102" "10.0.0.100" "PC102-PC101"
test_ping "clab-frr01-pc102" "10.1.0.100" "PC102-PC111"
test_ping "clab-frr01-pc102" "10.1.1.100" "PC102-PC112"
test_ping "clab-frr01-pc102" "10.2.0.100" "PC102-PC121"
test_ping "clab-frr01-pc102" "10.3.0.100" "PC102-PC131"
test_ping "clab-frr01-pc102" "192.168.1.100" "PC102-PC191"
test_ping "clab-frr01-pc102" "10.20.0.100" "PC102-PC201"
test_ping "clab-frr01-pc102" "10.30.0.100" "PC102-PC301"
test_ping "clab-frr01-pc102" "192.168.4.4" "PC102-CR104"
test_ping "clab-frr01-pc102" "192.168.5.5" "PC102-CR105"
test_no_ping "clab-frr01-pc102" "192.168.4.100" "PC102-CPC104"
test_ping "clab-frr01-pc102" "192.168.5.100" "PC102-CPC105"

# ========================== #
# Test pings from PC103      #
# ========================== #
test_ping "clab-frr01-pc103" "10.0.0.100" "PC103-PC101"
test_ping "clab-frr01-pc103" "10.0.2.100" "PC103-PC102"
test_ping "clab-frr01-pc103" "10.1.0.100" "PC103-PC111"
test_ping "clab-frr01-pc103" "10.1.1.100" "PC103-PC112"
test_ping "clab-frr01-pc103" "10.2.0.100" "PC103-PC121"
test_ping "clab-frr01-pc103" "10.3.0.100" "PC103-PC131"
test_ping "clab-frr01-pc103" "192.168.1.100" "PC103-PC191"
test_ping "clab-frr01-pc103" "10.20.0.100" "PC103-PC201"
test_ping "clab-frr01-pc103" "10.30.0.100" "PC103-PC301"
test_ping "clab-frr01-pc103" "192.168.4.4" "PC103-CR104"
test_ping "clab-frr01-pc103" "192.168.5.5" "PC103-CR105"
test_no_ping "clab-frr01-pc103" "192.168.4.100" "PC103-CPC104"
test_ping "clab-frr01-pc103" "192.168.5.100" "PC103-CPC105"

# ========================== #
# Test pings from PC111      #
# ========================== #
test_ping "clab-frr01-pc111" "10.0.0.100" "PC111-PC101"
test_ping "clab-frr01-pc111" "10.0.3.100" "PC111-PC103"
test_ping "clab-frr01-pc111" "10.1.1.100" "PC111-PC112"
test_ping "clab-frr01-pc111" "10.2.0.100" "PC111-PC121"
test_ping "clab-frr01-pc111" "10.3.0.100" "PC111-PC131"
test_ping "clab-frr01-pc111" "192.168.1.100" "PC111-PC191"
test_ping "clab-frr01-pc111" "10.20.0.100" "PC111-PC201"
test_ping "clab-frr01-pc111" "10.30.0.100" "PC111-PC301"
test_ping "clab-frr01-pc111" "192.168.4.4" "PC111-CR104"
test_ping "clab-frr01-pc111" "192.168.5.5" "PC111-CR105"
test_no_ping "clab-frr01-pc111" "192.168.4.100" "PC111-CPC104"
test_ping "clab-frr01-pc111" "192.168.5.100" "PC111-CPC105"

# ========================== #
# Test pings from PC112      #
# ========================== #
test_ping "clab-frr01-pc112" "10.0.0.100" "PC112-PC101"
test_ping "clab-frr01-pc112" "10.2.0.100" "PC112-PC102"
test_ping "clab-frr01-pc112" "10.1.0.100" "PC112-PC111"
test_ping "clab-frr01-pc112" "10.2.0.100" "PC112-PC121"
test_ping "clab-frr01-pc112" "10.3.0.100" "PC112-PC131"
test_ping "clab-frr01-pc112" "192.168.1.100" "PC112-PC191"
test_ping "clab-frr01-pc112" "10.20.0.100" "PC112-PC201"
test_ping "clab-frr01-pc112" "10.30.0.100" "PC112-PC301"
test_ping "clab-frr01-pc112" "192.168.4.4" "PC112-CR104"
test_ping "clab-frr01-pc112" "192.168.5.5" "PC112-CR105"
test_no_ping "clab-frr01-pc112" "192.168.4.100" "PC112-CPC104"
test_ping "clab-frr01-pc112" "192.168.5.100" "PC112-CPC105"

# ========================== #
# Test pings from r111       #
# ========================== #
# fails: because router 191 only knows the client networks, causing that ping doesn’t find return path
test_no_ping "clab-frr01-r111" "192.168.1.100" "R111-PC191"

# ========================== #
# Test pings from PC121      #
# ========================== #
test_ping "clab-frr01-pc121" "10.0.0.100" "PC121-PC101"
test_ping "clab-frr01-pc121" "10.1.0.100" "PC121-PC111"
# test_ping "clab-frr01-pc121" "10.2.0.100" "PC121-PC121"
test_ping "clab-frr01-pc121" "10.3.0.100" "PC121-PC131"
test_ping "clab-frr01-pc121" "192.168.1.100" "PC121-PC191"
test_ping "clab-frr01-pc121" "192.168.11.100" "PC121-PC193"
test_ping "clab-frr01-pc121" "10.20.0.100" "PC121-PC201"
test_ping "clab-frr01-pc121" "10.20.3.100" "PC121-PC203"
test_ping "clab-frr01-pc121" "10.20.4.100" "PC121-PC204"
test_ping "clab-frr01-pc121" "10.30.0.100" "PC121-PC301"
test_ping "clab-frr01-pc121" "192.168.4.4" "PC121-CR104"
test_ping "clab-frr01-pc121" "192.168.5.5" "PC121-CR105"
test_no_ping "clab-frr01-pc121" "192.168.4.100" "PC121-CPC104"
test_ping "clab-frr01-pc121" "192.168.5.100" "PC121-CPC105"

# ========================== #
# Test pings from PC131      #
# ========================== #
test_ping "clab-frr01-pc131" "10.0.0.100" "PC131-PC101"
test_ping "clab-frr01-pc131" "10.0.3.100" "PC131-PC103"
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

# ========================== #
# Test pings from PC191      #
# ========================== #
test_ping "clab-frr01-pc191" "10.0.0.100" "PC191-PC101"
test_ping "clab-frr01-pc191" "10.0.2.100" "PC191-PC102"
test_ping "clab-frr01-pc191" "10.0.3.100" "PC191-PC103"
test_ping "clab-frr01-pc191" "10.1.0.100" "PC191-PC111"
test_ping "clab-frr01-pc191" "10.2.0.100" "PC191-PC121"
test_ping "clab-frr01-pc191" "10.3.0.100" "PC191-PC131"
# test_ping "clab-frr01-pc191" "192.168.1.100" "PC191-PC191"
test_ping "clab-frr01-pc191" "10.20.0.100" "PC191-PC201"
test_ping "clab-frr01-pc191" "10.20.3.100" "PC191-PC203"
test_ping "clab-frr01-pc191" "10.20.4.100" "PC191-PC204"
test_no_ping "clab-frr01-pc191" "10.30.0.100" "PC191-PC301"
# the next four test have to fail because no matching route exists on R191
test_no_ping "clab-frr01-pc191" "192.168.4.4" "PC191-CR104"
test_no_ping "clab-frr01-pc191" "192.168.5.5" "PC191-CR105"
test_no_ping "clab-frr01-pc191" "192.168.4.100" "PC191-CPC104"
test_no_ping "clab-frr01-pc191" "192.168.5.100" "PC191-CPC105"

# ========================== #
# Test pings from PC192      #
# ========================== #
test_ping "clab-frr01-pc192" "10.0.0.100" "PC192-PC101"
test_ping "clab-frr01-pc192" "10.1.0.100" "PC192-PC111"
test_ping "clab-frr01-pc192" "10.1.1.100" "PC192-PC112"
test_ping "clab-frr01-pc192" "10.2.0.100" "PC192-PC121"
test_ping "clab-frr01-pc192" "10.3.0.100" "PC192-PC131"
# Because no static route exists (both directions)
test_no_ping "clab-frr01-pc192" "192.168.1.100" "PC192-PC191"
test_no_ping "clab-frr01-pc192" "10.20.0.100" "PC192-PC201"
# Because static is not redistributed into bgp
test_no_ping "clab-frr01-pc192" "192.168.32.100" "PC192-PC392"
test_no_ping "clab-frr01-pc192" "192.168.33.100" "PC192-PC393"
test_no_ping "clab-frr01-pc192" "192.168.34.100" "PC192-PC394"
# Because static is redistributed into bgp
test_ping "clab-frr01-pc192" "10.30.0.100" "PC192-PC301"

# ========================== #
# Test pings from PC193      #
# ========================== #
test_ping "clab-frr01-pc193" "10.0.0.100" "PC193-PC101"
test_ping "clab-frr01-pc193" "10.1.0.100" "PC193-PC111"
test_ping "clab-frr01-pc193" "10.1.1.100" "PC193-PC112"
test_ping "clab-frr01-pc193" "10.2.0.100" "PC193-PC121"
test_ping "clab-frr01-pc193" "10.3.0.100" "PC193-PC131"
test_ping "clab-frr01-pc193" "192.168.1.100" "PC193-PC191"
test_ping "clab-frr01-pc193" "10.20.0.100" "PC193-PC201"
# Because static is not redistributed into bgp
test_no_ping "clab-frr01-pc193" "192.168.32.100" "PC193-PC392"
test_no_ping "clab-frr01-pc193" "192.168.33.100" "PC193-PC393"
test_no_ping "clab-frr01-pc193" "192.168.34.100" "PC193-PC394"
# Because static is redistributed into bgp
test_ping "clab-frr01-pc193" "10.30.0.100" "PC193-PC301"

# ========================== #
# Test pings from PC201      #
# ========================== #
test_ping "clab-frr01-pc201" "10.0.0.100" "PC201-PC101"
test_ping "clab-frr01-pc201" "10.1.0.100" "PC201-PC111"
test_ping "clab-frr01-pc201" "10.1.1.100" "PC201-PC112"
test_ping "clab-frr01-pc201" "10.2.0.100" "PC201-PC121"
test_ping "clab-frr01-pc201" "10.3.0.100" "PC201-PC131"
test_ping "clab-frr01-pc201" "192.168.1.100" "PC201-PC191"
test_ping "clab-frr01-pc201" "10.20.3.100" "PC201-PC203"
test_ping "clab-frr01-pc201" "10.20.4.100" "PC201-PC204"
test_no_ping "clab-frr01-pc201" "10.30.0.100" "PC201-PC301"
test_ping "clab-frr01-pc201" "192.168.4.4" "PC201-CR104"
test_ping "clab-frr01-pc201" "192.168.5.5" "PC201-CR105"
test_no_ping "clab-frr01-pc201" "192.168.4.100" "PC201-CPC104"
test_ping "clab-frr01-pc201" "192.168.5.100" "PC201-CPC105"

# ========================== #
# Test pings from PC203      #
# ========================== #
test_ping "clab-frr01-pc203" "10.0.0.100" "PC203-PC101"
test_ping "clab-frr01-pc203" "10.1.0.100" "PC203-PC111"
test_ping "clab-frr01-pc203" "10.1.1.100" "PC203-PC112"
test_ping "clab-frr01-pc203" "10.2.0.100" "PC203-PC121"
test_ping "clab-frr01-pc203" "10.3.0.100" "PC203-PC131"
test_ping "clab-frr01-pc203" "192.168.1.100" "PC203-PC191"
test_ping "clab-frr01-pc203" "10.20.0.100" "PC203-PC201"
test_ping "clab-frr01-pc203" "10.20.4.100" "PC203-PC204"
test_no_ping "clab-frr01-pc203" "10.30.0.100" "PC203-PC301"

# ========================== #
# Test pings from PC204      #
# ========================== #
test_ping "clab-frr01-pc204" "10.0.0.100" "PC204-PC101"
test_ping "clab-frr01-pc204" "10.1.0.100" "PC204-PC111"
test_ping "clab-frr01-pc204" "10.1.1.100" "PC204-PC112"
test_ping "clab-frr01-pc204" "10.2.0.100" "PC204-PC121"
test_ping "clab-frr01-pc204" "10.3.0.100" "PC204-PC131"
test_ping "clab-frr01-pc204" "192.168.1.100" "PC204-PC191"
test_ping "clab-frr01-pc204" "10.20.0.100" "PC204-PC201"
test_ping "clab-frr01-pc204" "10.20.3.100" "PC204-PC203"
test_no_ping "clab-frr01-pc204" "10.30.0.100" "PC204-PC301"
test_ping "clab-frr01-pc204" "192.168.4.4" "PC204-CR104"
test_ping "clab-frr01-pc204" "192.168.5.5" "PC204-CR105"
test_no_ping "clab-frr01-pc204" "192.168.4.100" "PC204-CPC104"
test_ping "clab-frr01-pc204" "192.168.5.100" "PC204-CPC105"

# ========================== #
# Test pings from PC301      #
# ========================== #
test_ping "clab-frr01-pc301" "10.0.0.100" "PC301-PC101"
test_ping "clab-frr01-pc301" "10.1.0.100" "PC301-PC111"
test_ping "clab-frr01-pc301" "10.1.1.100" "PC301-PC112"
test_ping "clab-frr01-pc301" "10.2.0.100" "PC301-PC121"
test_ping "clab-frr01-pc301" "10.3.0.100" "PC301-PC131"
test_no_ping "clab-frr01-pc301" "192.168.1.100" "PC301-PC191"
test_ping "clab-frr01-pc301" "192.168.33.100" "PC301-PC393"
test_ping "clab-frr01-pc301" "192.168.34.100" "PC301-PC394"
test_no_ping "clab-frr01-pc301" "10.20.0.100" "PC301-PC201"
# test_ping "clab-frr01-pc301" "10.30.0.100" "PC301-PC301"
test_ping "clab-frr01-pc301" "192.168.32.100" "PC301-PC392"

# ========================== #
# Test pings from PC393      #
# ========================== #
test_no_ping "clab-frr01-pc393" "10.0.0.100" "PC393-PC101"
test_no_ping "clab-frr01-pc393" "10.1.0.100" "PC393-PC111"
test_no_ping "clab-frr01-pc393" "10.2.0.100" "PC393-PC121"
test_no_ping "clab-frr01-pc393" "10.3.0.100" "PC393-PC131"
# test_ping "clab-frr01-pc393" "192.168.30.100" "PC393-PC391"
test_no_ping "clab-frr01-pc393" "10.20.0.100" "PC393-PC201"
# only this test will work because R393 has only one static route (this)
test_ping "clab-frr01-pc393" "10.30.0.100" "PC393-PC301"

# ^^^^^^^^^^^^^^^^^^^^^^^^^^ PING TESTS ^^^^^^^^^^^^^^^^^^^^^^^^^^ #

# Wait for all background pings to finish
wait "${PING_PIDS[@]}"

echo "⬇️ All ping tests completed! ⬇️"

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
        elif ! grep -q "100% packet loss" "$log_file" && ! grep -q " 0% packet loss" "$log_file"; then
          echo "✅ Success: Ping from $description passed! ❌ BUT NOT 100% success rate"
          echo "❌ FAILURE: Ping from $description did not return 0% packet loss" >> "$LOG_FILE"
          cat "$log_file" >> "$LOG_FILE"
          NOT_100PERCENT_SUCCESS=1
        else
          echo "❌ FAILURE: Ping from $description failed!"
          echo "❌ FAILURE: Ping from $description failed!" >> "$LOG_FILE"
          cat "$log_file" >> "$LOG_FILE"
          FAILED_TESTS=1
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
    echo -e "\n\n⚠️ Some ping tests failed or unexpected successes occurred. Check the log file: $LOG_FILE\n"
    exit 1
elif [ "$NOT_100PERCENT_SUCCESS" -ne 0 ];then
    echo -e "\n\n🎉⚠️ Some ping tests had not 100% success rate. Check the log file: $LOG_FILE\n"
    exit 1
else
    echo -e "\n\n🎉 All connectivity tests passed successfully!\n"
    exit 0
fi
