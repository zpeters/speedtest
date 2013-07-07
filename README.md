speedtest
=========
This is a quick client for speedtest.net in go.  Patterned after https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

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
```

TODO
====
- [x] More code cleanup
- [x] move stuff in main into it's own functions
- [x] move some part into their own packages
- [x] submit to github
- [ ] test download speeds against speedtest.net to make sure measurements are correct, dl ususally seems slower
- [ ] add more timeout/error checking around servers
https://github.com/zpeters/speedtest/issues/1
- [ ] specify server
http://www.reddit.com/r/sysadmin/comments/1ht86k/command_line_interface_to_speedtestnet/caxrn65
- [ ] list servers
- [ ] add license?

WISHLIST
=======
- [ ] semi-automate builds and new releases with git hooks on push
- [ ] daemon mode that does a continuous test/recording for graphing
- [ ] config switches for amount to download / number of downloads to perform