Just exploring fuse.
Mounting array of array of strings as a file system.
Directory and file names are array indexes.
```
$ umount /tmp/test ;  ./arrayfs.exe -s -d /tmp/test
$ tree /tmp/test
/tmp/test
├── 0
│   ├── 0
│   └── 5
├── 1
│   └── 1
├── 2
│   └── 1
├── 3
│   └── 9
├── 4
│   └── 7
└── 8
    └── 1

6 directories, 7 files

$ find /tmp/test -type f -exec rg . {} \;
1:Hi.
1:this
1:is
1:only
1:a
1:test
1:!
