speedtest
=========
This is a quick client for speedtest.net in go.  Patterned after https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

License
=======
Licensed under GPLv3 (See COPYING and LICENSE)

Version
=======
0.04

Download
========
- Linux - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.04/Linux/
- Windows - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.04/Windows/

Usage
=====
```shell
speedtest.exe - normal run
speedtest.exe -d - turn on debugging
speedtest.exe -v - show version
speedtest.exe -l - list servers (hint, you can "grep" (or "findstr" in Windows) by location, name, etc to get the URL)
```

Feedback / Contributing
=======================
Contact zpeters@gmail.com for general feedback

For Bug reports please use the Github issue tracker for this project

To contribute please see CONTRIBUTING.md

Thank You
=========
- Jacob McDonald - jmc734 - Cleaned up printing and formatting.  Added parameter passing to run.sh - https://github.com/zpeters/speedtest/pull/4

Todo / Wishlist
===============
Migrated to github "issues" - https://github.com/zpeters/speedtest/issues?state=open

Done
====
- [x] add more timeout/error checking around servers
https://github.com/zpeters/speedtest/issues/1
- [x] More code cleanup
- [x] move stuff in main into it's own functions
- [x] move some part into their own packages
- [x] submit to github
- [X] list servers
- [X] add license?
- [x] specify server
http://www.reddit.com/r/sysadmin/comments/1ht86k/command_line_interface_to_speedtestnet/caxrn65
