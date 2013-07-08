speedtest
=========
This is a quick client for speedtest.net in go.  Patterned after https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

License
=======
Licensed under GPLv3 (See COPYING and LICENSE)

Version
=======
0.03

Download
========
- Linux - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.03/Linux/
- Windows - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.03/Windows/

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

Todo
====
- [ ] try to get rid of globals...ick
- [ ] review FIXME's and clean up ugly code

Wishlist
=======
- [ ] semi-automate builds and new releases with git hooks on push
- [ ] daemon mode that does a continuous test/recording for graphing
- [ ] config switches for amount to download / number of downloads to perform

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
