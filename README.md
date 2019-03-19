2019 Update
==========
I don't have a lot of time to work on this code anymore. I may post updates from time to time, but at the moment this software is mostly abandonded.  I am working on a Rust implementation and will post more details here in the future. See https://github.com/zpeters/speedtestr

Thank you for all of the fun and support over the years.

-zach


VERSION 2.0 Testing
===================
Initial testing release of v2.0 is out for the testing. See "releases" for downloads.  The current "test" is hard coded and there are no options at the moment. Please send me any feedback at zpeters@gmail.com or through the issues.

The Unofficial Speedtest CLI
============================
The Unofficial Speedtest CLI is a command-line program to test
bandwidth in situations where you don't have access to a full GUI
environment and web browser.

In [2013 I was feeling guilty](http://thehelpfulhacker.net/2013/07/29/giving-something-back/)
about using Open Source software for most of my life without giving
anything back in return.  I decided to create this project to my part
to help the IT community.

A lot of the initial algorithms here are based on different scripts I
found when I was studying how speedtest.net works.  Mainly, @sivel's
[speedtest-cli](https://github.com/sivel/speedtest-cli),
thanks for your work!

**master branch**
[![Go Report Card](https://goreportcard.com/badge/github.com/zpeters/speedtest)](https://goreportcard.com/report/github.com/zpeters/speedtest)
[![Github All Releases](https://img.shields.io/github/downloads/zpeters/speedtest/total.svg?style=plastic)](https://www.somsubhra.com/github-release-stats/?username=zpeters&repository=speedtest)
[![Build Status](https://travis-ci.org/zpeters/speedtest.svg?branch=master)](https://travis-ci.org/zpeters/speedtest)
[![GoDoc](https://godoc.org/github.com/zpeters/speedtest?status.svg)](https://godoc.org/github.com/zpeters/speedtest)

**development branch**
[![Build Status](https://travis-ci.org/zpeters/speedtest.svg?branch=develop)](https://travis-ci.org/zpeters/speedtest)

[![Sparkline](https://stars.medv.io/zpeters/speedtest.svg)](https://stars.medv.io/zpeters/speedtest)

License
=======
Licensed under GPLv3 (See COPYING and LICENSE)

Download
========
- Github (Windows/Linux/Mac) - https://github.com/zpeters/speedtest/releases
- Mirror (Windows/Linux/Mac) - http://media.thehelpfulhacker.net/index.php?dir=speedtest/

Build
=====
See [Build Instructions](https://github.com/zpeters/speedtest/wiki/Build-Instructions)

Bugs, Features and Contributing
=================
See github issues tracker - https://github.com/zpeters/speedtest/issues

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

speedtest.exe -b 1234 -b 5678 -- Run the test blacklisting servers 1234 and 5678
speedtest.exe -r -- Runs speedtest in "reporting" mode (useful for Labtec, Excel spreadsheets, etc)
speedtest.exe -r -rc="," -- Use a different separator (default is '|')
Report Fields: Server ID, Server Name (Location), Ping time in ms, Download speed in kbps, Upload speed in kbps
```shell
1752|5NINES(Madison, WI,United States)|36.18|19452|4053
```

```shell
COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --algo value, -a value          Specify the measurement method to use ('max', 'avg')
   --debug, -d                     Turn on debugging
   --list, -l                      List available servers
   --update, -u                    Check for a new version of speedtest
   --ping, -p                      Ping only mode
   --quiet, -q                     Quiet mode
   --report, -r                    Reporting mode output, minimal output with '|' for separators, use '--rc'
                                     to change separator characters. Reports the following: Server ID,
                                     Server Name (Location), Ping time in ms, Download speed in kbps, Upload speed in kbps
   --downloadonly, --do            Only perform download test
   --uploadonly, --uo              Only perform upload test
   --reportchar value, --rc value  Set the report separator. Example: --rc=','
   --server value, -s value        Use a specific server
   --blacklist value, -b value     Blacklist a server.  Use this multiple times for more than one server
   --mini value, -m value          URL of speedtest mini server
   --useragent value, --ua value   Specify a useragent string
   --numclosest value, --nc value  Number of 'closest' servers to find (default: 3)
   --httptimeout value, -t value   Timeout (seconds) for http connections (default: 15)
   --numlatency value, --nl value  Number of latency tests to perform (default: 5)
   --interface value, -I value     Source IP address or name of an interface
   --help, -h                      show help
   --version, -v                   print the version
```

Thank You
=========
- Jacob McDonald - jmc734 - Cleaned up printing and formatting.  Added parameter passing to run.sh - https://github.com/zpeters/speedtest/pull/4
- Cory Lievers - Testing and feedback. Suggestions for formatting to make this more useful for labtec - https://github.com/zpeters/speedtest/issues/9
- Paul Baker (Network Manager - BMS Telecorp) - Located a bug in the speedtest.net server list generation and found the correct 'static' url
- Graham Roach (Contact Info?) - Extensive user testing to help determine issues with latency and accuracy of upload and download speeds - #11 (and others)
- @larray - slightly obscure issues with http caches interfering with test results - #20
- Noric - reporting and help with testing issues with report formatting - #32
- @jannson - submitting patch to reduce memory usage on download test - #37
- @vendion - teaching me how to import packages the corret way - #67
- @invalid-email-address - various formatting
- @l2dy - cleaned up README and broken links
- @m01 - speed test mini support
- @pra85 - fixed types in README
- @schweikert - for adding the interface selection code

Why don't my speeds match those reported from the speedtest.net website?
========================================================================
The calculation that is used for testing download speeds is literally measuring the amount of data we are downloading (we request a "random" image and count how many bytes are received) and how long it takes to download.  We multiply by the correct factors to get from bytes to megabits. I consider this to be an honest and accurate measurement.

In speedtest.net's reference documentation they describe doing a lot of manipulation to the results to return an "ideal" measurement (https://support.speedtest.net/entries/20862782-How-does-the-test-itself-work-How-is-the-result-calculated-). This, to me, is trading accuracy for speed and not what I'm looking for out of a testing tool.

For confirmation that my download calculations are correct I have tested against a few other speed testing sites, specifically http://testmy.net ("What makes TestMy.net better") who appear to use an "unfiltered" method of calculating bandwidth speeds.  These results typically match up with speedtest.net cli


Reference
=========
- how does it work - https://support.speedtest.net/entries/20862782-How-does-the-test-itself-work-How-is-the-result-calculated-
- why actual speedtest.net results may be inaccurate - http://testmy.net/
   
