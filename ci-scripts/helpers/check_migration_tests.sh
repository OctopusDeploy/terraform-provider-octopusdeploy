#!/bin/bash

# Check for missing migration tests when SDKv2 resources are replaced with TPF resources
# Usage: check_migration_tests.sh <base_ref>

set -e

BASE_REF=${1:-origin/main}
CURRENT_REF=${GITHUB_REF:-HEAD}

echo "Checking for migration tests between $BASE_REF and $CURRENT_REF"

# Get deleted SDKv2 and new TPF resources
deleted_sdkv2=$(git diff --name-status $BASE_REF..$CURRENT_REF | grep '^D.*octopusdeploy/resource_.*\.go$' | grep -v '_test\.go$' || true)
new_tpf=$(git diff --name-status $BASE_REF..$CURRENT_REF | grep '^A.*octopusdeploy_framework/resource_.*\.go$' | grep -v '_test\.go$' | grep -v '_migration_test\.go$' || true)

if [ ! -z "$deleted_sdkv2" ] && [ ! -z "$new_tpf" ]; then
  missing_tests=""
  
  # Extract resource names and check for matches
  for deleted_file in $deleted_sdkv2; do
    deleted_resource=$(echo "$deleted_file" | sed 's/.*resource_\(.*\)\.go$/\1/')
    
    for new_file in $new_tpf; do
      new_resource=$(echo "$new_file" | sed 's/.*resource_\(.*\)\.go$/\1/')
      
      if [ "$deleted_resource" = "$new_resource" ]; then
        migration_test="octopusdeploy_framework/resource_${deleted_resource}_migration_test.go"
        
        # Check if migration test exists or is being added
        if ! git diff --name-status $BASE_REF..$CURRENT_REF | grep -q "$migration_test" && [ ! -f "$migration_test" ]; then
          missing_tests="$missing_tests $deleted_resource"
        fi
      fi
    done
  done
  
  if [ ! -z "$missing_tests" ]; then
    echo "::error::Migration tests required when replacing SDKV2 resources with TPF resources."
    echo "Missing migration tests for:$missing_tests"
    echo "Required files:"
    for resource in $missing_tests; do
      echo "  - octopusdeploy_framework/resource_${resource}_migration_test.go"
    done
    exit 1
  fi
fi

echo "::notice::Migration test requirements satisfied."