#!/bin/sh
# POSIX-compliant script to log command details to a Unix socket

# Parameters:
# $1 - Execution phase (start/end)
# $2 - Command to log
# $3 - Working directory
# $4 - User who executed the command
# $5 - Timestamp of command execution (start time, Unix timestamp)
# For end: $6 - End time, $7 - Duration

# UNIX socket path
SOCKET_PATH="{{.SocketPath}}"

# Construct the log message with start time, end time, and duration
if [ "$1" = "start" ]; then
    LOG_MESSAGE="start|$2|$3|$4|$5"
elif [ "$1" = "end" ]; then
    LOG_MESSAGE="end|$2|$3|$4|$5|$6|$7"
fi

# Function to check command existence
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Python function to communicate via socket
send_via_python() {
    python -c "import socket; import sys; \
    s = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM); \
    s.connect('$SOCKET_PATH'); \
    s.sendall('$LOG_MESSAGE'.encode('utf-8')); \
    s.close()" 2>/dev/null
}

# Perl function to communicate via socket
send_via_perl() {
    perl -e "use IO::Socket::UNIX; \
    \$sock = IO::Socket::UNIX->new(Type => SOCK_STREAM, Peer => '$SOCKET_PATH'); \
    print \$sock '$LOG_MESSAGE'; \
    close(\$sock);" 2>/dev/null
}

# Send the log message to the Go application via UNIX socket
if command_exists nc; then
    echo "$LOG_MESSAGE" | nc -U "$SOCKET_PATH"
elif command_exists socat; then
    echo "$LOG_MESSAGE" | socat - UNIX-CONNECT:"$SOCKET_PATH"
elif command_exists python; then
    send_via_python
elif command_exists perl; then
    send_via_perl
else
    # TODO: Implement direct command sending
    echo "Neither nc, socat, python or perl are available on this system."
    exit 1
fi
