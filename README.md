# james-d-elliott/go-hekatechkbkp

[![Build Status](https://travis-ci.com/james-d-elliott/go-hekatechkbkp.svg?branch=master)](https://travis-ci.com/james-d-elliott/go-hekatechkbkp)

CLI utility to verify backups that have accompanying .sha256sums files, specifically designed to work with hekate style backups on the Switch.

## Features

- Verifies chunks of files against sha256 sums
- Logs details about the verification
- Outputs .invalid_chunks files allowing the possibility to redump them individually in the future

## Usage

This tool is a CLI tool. It will run if you just open it, but it is not the intention of the tool. Additionally if you do this it immidiately closes. You can still view the log file if you wish to see the output however. The following flags exist and will customize how the tool is used. The dir flag currently has only been tested against relative paths.

  -help

        shows CLI help

  -debug int

        sets the debug level which affects how many things are sent to the log

  -dir string

        set the source directory of the files to validate (default "./")

  -log

        disables/enables the log sending output to the log file (default true)

  -log-console

        disables/enables the log sending output to the console (default true)

  -log-file string
  
        changes the filename of the log file (default "verification.log")

## To Do

- Improve logging and error handling
- Implement rehashing (flag to rehash files as the process occurs)
- Implement part file combining (flag to combine part files for fat32, should also combine any checksum list files)
- Implement a --no-verify flag so users can skip verification (so they can just combine and hash at the end)
- Implement combining sha256sums files when combining occurs
- Implement test units
- Implement logic to handle crc32, sha1, and/or sha256 (will have to refactor to check for existence of all files to save resources, we only want to read once, we're not IO barbarians)