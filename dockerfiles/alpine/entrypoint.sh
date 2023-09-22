#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -o errexit
# Exit immediately if an unset variable is used
set -o nounset

# Define the directory containing entrypoint scripts
ENTRYPOINT_DIR="/docker-entrypoint.d"

# If the entrypoint directory exists and contains scripts
if [ -d "$ENTRYPOINT_DIR" ] && [ -n "$(ls -A "$ENTRYPOINT_DIR")" ]; then
    echo "Running entrypoint scripts in $ENTRYPOINT_DIR"

    # Sort the scripts in numerical order based on their names
    # and execute them in ascending order, skipping the script itself (if it's present)
    for script in $(ls -v "$ENTRYPOINT_DIR"/*); do
        if [ -x "$script" ] && [ ! "${script##*/}" = "$(basename "$0")" ]; then
            "$script"
        fi
    done
else
    echo "no scripts found in $ENTRYPOINT_DIR"
fi

# Check if the $ENTRYPOINT variable is not empty
if [[ -n "${ENTRYPOINT+x}" ]]; then
    echo "$0: Fetching entrypoint files from $ENTRYPOINT"

    # Create a directory called /entrypoint.d/ if it doesn't exist
    mkdir -p /entrypoint.d/

    # Copy files recursively from the specified S3 bucket ($ENTRYPOINT) to the local directory /entrypoint.d/
    aws s3 cp s3://$ENTRYPOINT/ /itoentrypoint.d/ --recursive

    # Give executable permissions (+x) to all the shell script files in /entrypoint.d/
    chmod +x /entrypoint.d/*.sh

    # Check if there are any files in /entrypoint.d/
    if /usr/bin/find /entrypoint.d/ -mindepth 1 -maxdepth 1 -type f -print -quit 2>/dev/null | read v; then
        echo "$0: /entrypoint.d/ is not empty, will attempt to perform configuration"

        # Look for shell scripts in /entrypoint.d/ and process them in version order
        echo "$0: Looking for shell scripts in /entrypoint.d/"
        find /entrypoint.d/ -follow -type f -print | sort -V | while read -r f; do
            case "$f" in
            # If the file has a .sh extension and is executable, launch it
            *.sh)
                if [ -x "$f" ]; then
                    echo "$0: Launching $f"
                    "$f"
                else
                    # If the file is not executable, ignore it and issue a warning
                    echo "$0: Ignoring $f, not executable"
                fi
                ;;
            # If the file doesn't have a .sh extension, ignore it
            *) echo "$0: Ignoring $f" ;;
            esac
        done
        echo "$0: Configuration complete; continuing start up"
    else
        echo "$0: No files found in /entrypoint.d/, skipping configuration"
    fi
else
    echo "$0: No custom entrypoint specified, skipping configuration"
fi

# Execute the command passed to the entrypoint
exec "$@"
