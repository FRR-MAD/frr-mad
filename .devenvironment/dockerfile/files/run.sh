#!/bin/sh

BINARY_PATH="/app/tmp/analyzer_frr"
LAST_MODIFIED=""
INIT_START=1

echo "Starting: $BINARY_PATH debug"

while true; do
    if [ -f "$BINARY_PATH" ]; then
        CURRENT_MODIFIED=$(stat -c %Y "$BINARY_PATH" 2>/dev/null)

        if [ "$CURRENT_MODIFIED" != "$LAST_MODIFIED" ]; then
            
            # Kill previous instance if running
            if [ ! -z "$PID" ]; then
                #kill $PID 2>/dev/null
                $BINARY_PATH stop
                wait $PID 2>/dev/null
            fi
            if [ $INIT_START -eq 0 ]; then
                echo "Restarting $BINARY_PATH debug..."
            fi
            # Start new instance
            $BINARY_PATH debug &
            PID=$!
            LAST_MODIFIED="$CURRENT_MODIFIED"
        fi
    fi
    INIT_START=0
    
    sleep 2
done