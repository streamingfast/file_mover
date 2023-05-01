# File mover

## Description

This script moves files from one directory to another. It is currently limited to `jpg` files.
It is listening to the kernel events for file changes in the source directory and moves the files to the destination directory. So there is not delay has soon a file appear in the src folder the transfert to the destination folder is started.

It is also monitoring the destination folder size and file are deleted from folder if the max size for that folder is reach. Each time the max size is reached 10% of the files are deleted.

## Usage

```bash
file_mover /tmp/src1 /tmp/dest1 1000000 /tmp/src2 /tmp/dest2 1000000
```

