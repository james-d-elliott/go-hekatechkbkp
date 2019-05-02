# james-d-elliott/go-hekatechkbkp

Utility to verify backups that have accompanying .sha1sums files, specifically designed to work with hekate style backups on the Switch.

## Features

- Verifies chunks of files against sha1 sums
- Logs details about the verification
- Outputs .invalid_chunks files allowing the possibility to redump them individually in the future

## To Do

- Implement rehashing (flag to rehash files as the process occurs)
- Implement part file combining (flag to combine part files for fat32)
- Implement combining sha1sums files when combining occurs
- Improve logging and error handling
- Implement test units