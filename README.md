# james-d-elliott/go-hekatechkbkp
[![Build Status](https://travis-ci.com/james-d-elliott/go-hekatechkbkp.svg?branch=master)](https://travis-ci.com/james-d-elliott/go-hekatechkbkp)

CLI utility to verify backups that have accompanying .sha256sums files, specifically designed to work with hekate style backups on the Switch.

## Features

- Verifies chunks of files against sha256 sums
- Logs details about the verification
- Outputs .invalid_chunks files allowing the possibility to redump them individually in the future

## To Do

- Improve logging and error handling
- Document the CLI flags
- Implement rehashing (flag to rehash files as the process occurs)
- Implement part file combining (flag to combine part files for fat32, should also combine any checksum list files)
- Implement a --no-verify flag so users can skip verification (so they can just combine and hash at the end)
- Implement combining sha256sums files when combining occurs
- Implement test units
- Implement logic to handle crc32, sha1, and/or sha256 (will have to refactor to check for existence of all files to save resources, we only want to read once, we're not IO barbarians)