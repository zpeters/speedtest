speedtest
=========
This is a quick client for speedtest.net in go.  Patterned after https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

[![Build Status](https://drone.io/github.com/zpeters/speedtest/status.png)](https://drone.io/github.com/zpeters/speedtest/latest)

License
=======
Licensed under GPLv3 (See COPYING and LICENSE)

Version
=======
0.09-dev

Bugs
======
- Latency seems very incorrect
- Upload speeds roughly half of speedtest.net 
- Implement - flag.StringVar(&ALGOTYPE, "a", "max", "\tSpecify the measurement method to use ('max', 'avg')") and 'min' for latency

Features to Add / Improvements
==============================
- move more stuff from speedtest.go to their own functions
- verify latency - first test is always higher
- currently we are using a very "dumb" way of testing speed (just downloading files and timing them).  review speedtest explanation and use a more sophistocated (faster) method of testing
- bump version and tag
- clean up switches they kind of a mess
- update binary links 
- add a thank you  for groach
- add "methodology" (re emails)

Download
========
- Linux - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.09/Linux/
- Windows - http://media.thehelpfulhacker.net/index.php?dir=speedtest/v0.09/Windows/

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

speedtest.exe -r -- Runs speedtest in "reporting" mode (useful for Labtec, Excel spreadsheets, etc)
speedtest.exe -r -rc="," -- Use a different separator (default is '|')
Report Fields: Server ID, Server Name (Location), Ping time in ms, Download speed in kbps, Upload speed in kbps
```shell
1752|5NINES(Madison, WI,United States)|36.18|19452|4053
```

```shell
speedtest.exe -h
Usage of bin\speedtest.exe:
  -d=false: Turn on debugging
  -dc=false: Turn on debugging and just dump config
  -l=false: List servers (hint use 'grep' or 'findstr' to locate a server ID to use for '-s'
  -nc=3: Number of geographically close servers to test to find the optimal serv er
  -nl=3: Number of latency tests to perform to determine which server is the fastest
  -q=false: Quiet Mode. Only output server and results
  -r=false: 'Reporting mode' output, minimal output with '|' for separators, use '-rc' to change separator characters. Reports the following: Server ID, Server Name (Location), Ping time in ms, Download speed in kbps, Upload speed in kbps
  -rc="|": Character to use to separate fields in report mode (-r)
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
- Cory Lievers - Testing and feedback. Suggestions for formatting to make this more useful for labtec - https://github.com/zpeters/speedtest/issues/9

Why don't my speeds match those reported from the speedtest.net website?
========================================================================
The calculation that is used for testing download speeds is literally measuring the amount of data we are downloading (we request a "random" image and count how many bytes are received) and how long it takes to download.  We multiply by the correct factors to get from bytes to megabits. I consider this to be an honest and accurate measurement.

In speedtest.net's reference documentation they describe doing a lot of manipulation to the restults to return an "ideal" measurement (https://support.speedtest.net/entries/20862782-How-does-the-test-itself-work-How-is-the-result-calculated-). This, to me, is trading accuracy for speed and not what I'm looking for out of a testing tool.

For confirmation that my download calculations are correct I have tested against a few other speed testing sites, specifically http://testmy.net ("What makes TestMy.net better") who appear to use an "unfiltered" method of calculating bandwidth speeds.  These results typically match up with speedtest.net cli


Reference
=========
- how does it work - https://support.speedtest.net/entries/20862782-How-does-the-test-itself-work-How-is-the-result-calculated-
- why actual speedtest.net results may be innaccurate - http://testmy.net/
