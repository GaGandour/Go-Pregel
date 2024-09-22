#!/bin/bash

# Get the file path argument
file_path="$1"

# Remove the optional leading slash and the prefix (assumed to be the first part after /)
cleaned_path="${file_path#/}"      # Remove the leading slash if it exists
cleaned_path="${cleaned_path#*/}"  # Remove the prefix (e.g., 'graphs')

# Extract the folder structure by removing the filename
folder_structure=$(dirname "$cleaned_path")

# Set the target directory where you want to create the folder structure
target_directory="."

# Create the full folder structure in the target directory
mkdir -p "$target_directory/$folder_structure"

# Print a message to confirm
echo "Created folder structure: $target_directory/$folder_structure"
