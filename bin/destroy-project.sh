#!/bin/ash

# If the argument is empty, then exit with error message
if [ -z "$1" ]; then
	echo "Client directory name is required"
	exit 1
fi
CLIENT_DIR_NAME="$1"
if [ -z "$2" ]; then
	echo "Project directory name is required"
	exit 1
fi
PROJECT_DIR_NAME="$2"

DOCUMENT_ROOT="/usr/local/apache2/htdocs/"

TARGET_DIR="${DOCUMENT_ROOT}${CLIENT_DIR_NAME}/${PROJECT_DIR_NAME}"

# If the target directory exists, then delete it and return success message.
if [ -d "${TARGET_DIR}" ]; then
	rm -rf "${TARGET_DIR}"
	echo "Project ${CLIENT_DIR_NAME}/${PROJECT_DIR_NAME} has been destroyed"
	# if the client directory is empty, then delete it
	if [ -z "$(ls -A ${DOCUMENT_ROOT}${CLIENT_DIR_NAME})" ]; then
		rm -rf "${DOCUMENT_ROOT}${CLIENT_DIR_NAME}"
		echo "Client ${CLIENT_DIR_NAME} has been destroyed"
	fi
	exit
fi
echo "Project ${CLIENT_DIR_NAME}/${PROJECT_DIR_NAME} has been destroyed already"
exit
