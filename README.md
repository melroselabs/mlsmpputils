# mlsmpputils
Code for working with SMPP

## pySMPPSubmit.py
Compact and minimal code that binds to an SMPP server and submits a message. No external dependencies.

Update first line (r,d,m variables ) with source address, destination address and message content.

Update second line (i,p,h,o) with SMPP account credentials (system ID, password, host/IP and port).

```
% python3 pySMPPSubmit.py 
Message ID: 80a905d0483cd2dc7e695620ea22e7892bed
```
