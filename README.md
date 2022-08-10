discord-attachment-backker-upper
===

Takes your Discord Privacy export, and goes and also fetches all of your attachments.

## Motivation

I send a lot of screenshots, and I've realised that if I lose screenshots I've lost a lot of context and history to chats.

Since discord does not export this data (even though services like Twitter and co do...), this tool does it.

## Usage:

```
$ ./discord-attachment-backker-upper ~/Downloads/Discord-10-Aug-2022.zip 
2022/08/10 11:39:51 Processing messages/c689.../messages.csv
2022/08/10 11:39:51 Processing messages/c325.../messages.csv
2022/08/10 11:44:00 Processing messages/c996.../messages.csv
```

It then puts all the files in the `-output.dir` location (default ./dump/), sorted by day.

```
./dump
./dump/2018-11-25
./dump/2018-11-25/5cd7994d154f22e191b0655d5e8f231c.png
./dump/2018-11-25/10fcb4536629c575370a10d9fcb0c77f.png
```
