speedtest
=========
This is a quick client for speedtest.net in go.  Patterned after https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

License
=======
Licensed under GPLv3 (See COPYING and LICENSE)

Version
=======
0.05

Download
========
- Linux - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.04/Linux/
- Windows - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.04/Windows/

Usage
=====
speedtest.exe -- normal run, will automatically select the closests/fastest server to test against
```shell
$ bin/speedtest.exe 
Finding fastest server..
1752 | 5NINES (Madison, WI, United States)
Testing download speed......
Testing upload speed......
Ping: 53.613233ms | Download: 13.34 Mbps | Upload: 3.89 Mbps
```

speedtest.exe -l -- List servers
```shell
$ bin/speedtest.ext -l
1724 | CityNet (Zaporizhzhya, Ukraine)
2966 | FUSION MEDIA Kft. (Budapest, Hungary)
3634 | Paul Bunyan Communications (Bemidji, MN, United States
...

```

speedtest.exe -s 1724 -- Run against a specific server
```shell
$ bin/speedtest.exe -s 1724
1724 | CityNet (Zaporizhzhya, Ukraine)
Testing latency...
Testing download speed......
Testing upload speed......
Ping: 982.913566ms | Download: 0.91 Mbps | Upload: 1.25 Mbps
```

```shell
speedtest.exe -h
  -d=false: Turn on debugging
  -l=false: List servers (hint use 'grep' or 'findstr' to locate a server ID to use for '-s'
  -nc=3: Number of geographically close servers to test to find the optimal server
  -nl=3: Number of latency tests to perform to determine which server is the fastest
  -q=false: Quiet Mode. Only output server and results
  -s="": Specify a server ID to use
  -v=false: Display version
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
