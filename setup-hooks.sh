#!/bin/bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
HOOKS_SOURCE_DIR=".hooks"
HOOKS_TARGET_DIR=".git/hooks"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    print_error "Not in a git repository root directory"
    exit 1
fi

# Check if .hooks directory exists
if [ ! -d "$HOOKS_SOURCE_DIR" ]; then
    print_error "Source hooks directory '.hooks' not found"
    exit 1
fi

# Create .git/hooks directory if it doesn't exist
if [ ! -d "$HOOKS_TARGET_DIR" ]; then
    print_status "Creating hooks directory: $HOOKS_TARGET_DIR"
    mkdir -p "$HOOKS_TARGET_DIR"
fi

# Initialize counters
installed_count=0
updated_count=0
skipped_count=0

# Process each file in the .hooks directory
for hook_file in "$HOOKS_SOURCE_DIR"/*; do
    print_status "--- Processing: $hook_file ---"

    # Skip if no files found (glob didn't match)
    if [ ! -e "$hook_file" ]; then
        print_error "File doesn't exist, skipping"
        continue
    fi

    # Get just the filename
    hook_name=$(basename "$hook_file")
    target_path="$HOOKS_TARGET_DIR/$hook_name"

    # Skip non-executable files and directories
    if [ -d "$hook_file" ]; then
        print_warning "Skipping directory: $hook_name"
        skipped_count=$((skipped_count+1))
        continue
    fi

    # Skip files that don't look like hooks
    if [[ "$hook_name" == *.* ]] && [[ "$hook_name" != *.sample ]]; then
        print_warning "Skipping non-hook file: $hook_name"
        skipped_count=$((skipped_count+1))
        continue
    fi

    # Make source file executable
    print_status "Making executable: $hook_file"
    if chmod +x "$hook_file"; then
        print_success "✓ Made executable: $hook_file"
    else
        print_error "✗ Failed to make executable: $hook_file"
        continue
    fi

    # Handle existing hooks
    if [ -L "$target_path" ]; then
        print_status "Target is a symlink"
        current_target=$(readlink "$target_path")
        print_status "Current target: $current_target"
        print_status "Expected target: ../../.hooks/$hook_name"

        if [ "$current_target" = "../../$HOOKS_SOURCE_DIR/$hook_name" ]; then
            print_success "✓ Hook '$hook_name' already properly linked"
            updated_count=$((updated_count+1))
        else
            print_warning "Hook '$hook_name' links to different location: $current_target"

            if ! rm "$target_path"; then
                print_error "Failed to remove existing symlink: $target_path"
                continue
            fi
            if ! ln -s "../../$HOOKS_SOURCE_DIR/$hook_name" "$target_path"; then
              print_error "Failed to remove existing symlink: $target_path"
              continue
            fi
            updated_count=$((updated_count+1))
            print_success "✓ Updated symlink"
        fi
    elif [ -f "$target_path" ]; then
        print_warning "Target is a regular file, backing up existing hook to: $backup_path"
        backup_path="${target_path}.backup.$(date +%Y%m%d_%H%M%S)"
        if ! mv "$target_path" "$backup_path"; then
            print_error "Failed to backup existing hook: $target_path"
            continue
        fi
        if ln -s "../../$HOOKS_SOURCE_DIR/$hook_name" "$target_path"; then
            print_error "Failed to create symlink: $target_path"
            continue
        fi
        print_success "✓ Backed up to $backup_path and created symlink"
        installed_count=$((installed_count+1))
    else
        # No existing hook - create symlink
        print_status "Installing hook: $hook_name"
        if ! ln -s "../../$HOOKS_SOURCE_DIR/$hook_name" "$target_path"; then
            print_error "Failed to create symlink: $target_path"
            continue
        fi
        print_success "✓ Created symlink: $target_path -> ../../$HOOKS_SOURCE_DIR/$hook_name"
        installed_count=$((installed_count+1))
    fi

    print_status "--- Finished processing $hook_name ---"
done

#echo ""
#print_status "=== FINAL STATE ==="
#print_status "Contents of .git/hooks directory (after):"
#ls -la .git/hooks/
#
#print_status ""
#print_status "Symlink verification:"
#for link in .git/hooks/*; do
#    if [ -L "$link" ]; then
#        echo "$(basename "$link") -> $(readlink "$link")"
#    fi
#done
#
#print_status ""
#print_status "=== DEBUG COMPLETE ==="

# Summary
echo
print_status "Summary:"
print_status "  - Installed: $installed_count hooks"
print_status "  - Updated: $updated_count hooks"
print_status "  - Skipped: $skipped_count files"

if [ $((installed_count + updated_count)) -eq 0 ]; then
    print_warning "No hooks were installed. Make sure you have valid hook files in $HOOKS_SOURCE_DIR/"
    print_status "Valid hook names: applypatch-msg, pre-applypatch, post-applypatch,"
    print_status "                  pre-commit, prepare-commit-msg, commit-msg,"
    print_status "                  post-commit, pre-rebase, post-checkout,"
    print_status "                  post-merge, pre-push, pre-receive, update, post-receive"
else
    echo
    print_success "Hooks are now active! Your versioned hooks from $HOOKS_SOURCE_DIR/ are linked to $HOOKS_TARGET_DIR/"
    print_status "To disable a hook temporarily, remove its symlink from $HOOKS_TARGET_DIR/"
    print_status "To update hooks, just modify files in $HOOKS_SOURCE_DIR/ and run this script again"
fi
