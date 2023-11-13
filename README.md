# mlsmpputils
Code for working with SMPP

## pySMPPSubmit.py
Compact and minimal SMS client written in Python that binds to an SMPP server and submits a message. No external dependencies.

Update first line (r,d,m variables ) with source address, destination address and message content.

Update second line (i,p,h,o) with SMPP account credentials (system ID, password, host/IP and port).

```
% python3 pySMPPSubmit.py 
Message ID: 80a905d0483cd2dc7e695620ea22e7892bed
```

## pySMPPReceive.py
Compact SMS client written in Python that binds as a receive to an SMPP server and receives messages. No external dependencies.

Update system_id, password, host and port variables with SMPP account credentials (system ID, password, host/IP and port).

```
% python3 pySMPPReceive.py
Successfully bound as receiver.
Received message: 447700111222447944654716
                                          Hello world!
Received message: 447700111222447944654716
                                          Hello world!
```
