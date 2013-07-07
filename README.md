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
- [ ] test download speeds against speedtest.net to make sure measurements are correct, dl ususally seems slower
- [x] submit to github
