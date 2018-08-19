From: https://gist.github.com/sdstrowes/411fca9d900a846a704f68547941eb97

Source: https://web.archive.org/web/20141216073338/https://gkbrk.com/blog/read?name=reverse_engineering_the_speedtest_net_protocol
Author: Gökberk Yaltıraklı

# Reverse Engineering the Speedtest.net Protocol

After finishing my command line speed tester written in Rust, I didn't have a proper blog to document this process. A few days ago I wrapped up a simple blogging script in Python so hopefully it works good enough to explain how everything works.

By now I have already figured out the whole protocol for performing a speed test but I will write all the steps that I took so you can learn how to reverse engineer a simple protocol.

The code that I wrote can be found at https://github.com/gkbrk/speedtest-rust.

# Finding the TCP stream in Wireshark

First of all lets open Wireshark and start sniffing. This allows us to find and view any connection made from our machine.
(Wireshark Screenshot)

After getting Wireshark ready I will go to speedtest.net and start a speed test. Then I will start checking Wireshark for possible connections.

By trial and error, I found the connection in Wireshark. It was connecting to speedtest.turk.net (a Turkish ISP, which makes sense since I'm in Turkey).
(DNS response and TCP stream screenshot)

After finding the connection, we can see that it uses port 8080, which is called http-alt in Wireshark.

If we right click the connection and do "Follow TCP Stream". If we found the right stream, we should see the data and we can start decoding the protocol.
(Follow TCP Stream dialog screenshot)

# Understanding The Protocol - The HI Command

By just looking at the data, we can tell that it's a plaintext protocol. That means instead of binary data, it uses text. This makes things a lot easier for us.

Here's the first part of the data. C is client and S is the server

```
C: HI
S: HELLO 2.1 2013-08-14.01
```

When the client sends the HI command, the server responds with its version.

The best part of plaintext protocols is that you can test them with a tool like netcat, without writing a single line of code. To test the HI command, let's open a terminal and run nc speedtest.turk.net 8080.

This creates a TCP stream that we can use to test the HI command. If we write HI and press enter; the server sends us its version, thinking that we are the speedtest.net app.

```
$ nc speedtest.turk.net 8080
HI
HELLO 2.1 2013-08-14.01
```

# Understanding The Protocol - The PING Command

Looking at the rest of the data, we see that the client sends a ping command followed by a timestamp. In response, the server sends "PONG timestamp". It goes like this.

```
C: PING 1418661866099
S: PONG 1418661866349
```

This happens 20 times. The original implementation takes the highest ping and displays it, mine gets the average and displays it.

This can be tested with netcat just like the HI command.

# Understanding The Protocol - The DOWNLOAD Command

After doing the handshake (the HI command) and the ping test, it is time to actually take a look at the download function.

It works like this; the client sends the download command followed by the number of bytes it wants to download, and the server responds with that many bytes, including the newline at the end of the response.

*NOTE*: The original implementation always responds with the first 8 bytes "DOWNLOAD", a space, and repeating "JABCDEFGHI".

It goes like this:

```
C: DOWNLOAD 14
S: DOWNLOAD JABC

C: DOWNLOAD 5
S: DOWN

C: DOWNLOAD 50
S: DOWNLOAD JABCDEFGHIJABCDEFGHIJABCDEFGHIJABCDEFGHI
```

Normally you ask for much more data to do a reliable speed test. You can test the DOWNLOAD function in netcat with small numbers like the ones above.

# Understanding The Protocol - The UPLOAD Command

The upload test is similar to the download test. To upload data to the server, you first send the UPLOAD command, followed by the number of bytes to upload and a zero. Then you upload the bytes. The server responds with "OK BYTES 0".

*NOTE*: There is a small detail that I missed. I actually fixed the bug in my code while writing this. When you send the data, instead of sending the number of bytes you said in the command, you need to send the number of bytes minus the number of bytes in the command itself.

With the note in mind, it goes like this:

```
C: UPLOAD 20 0
C: RANDOM1
S: OK 20 0
```

That is 7 bytes. That is because the length of "UPLOAD 20 0" is 11. If we add the newline after the upload command and the newline after the bytes to it, that makes 13. And 20-13 is 7.

# Understanding The Protocol - The QUIT Command

There isn't much to explain about the QUIT command. You send it, the server closes the connection without a response. You don't need to send it, but to comply with the original implementation, I do.

```
C:QUIT
```

After sending this command, the server terminates the connection and you can't read or write to it anymore.

# Implementing the Actual Speed Testing

Now that you know the whole network protocol of the speed testing servers, you can implement this protocol in your programs.

The servers send timestamps with the responses to PING and UPLOAD commands. I don't know if the original client uses those values for calculating the speed but I just use time::precise_time_ns() from the Rust time library before and after downloading/uploading the data and subtract those to get the time it took to download/upload it.

Using this method I got accurate results that matched with those on speedtest.net, so this method works.

For gathering the actual results, my implementation does:

* 20 pings and gets the average
* Downloads 1 MB 4 times and gets the average
* Uploads 1 MB 2 times and gets the average

# Getting the Serverlist

To get the coordinates(lat, lon) of the user, check out http://www.speedtest.net/speedtest-config.php.

To get a list of all the servers and their coordinates, check out http://www.speedtest.net/speedtest-servers-static.php.

You just need to parse these pages with an XML parser or use some regex trickery to get the server list data.
