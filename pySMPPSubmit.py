r,d,m="447700111222","447799222333","Hello world!"
i,p,h,o="SYSTEMID","PASSWORD",'smscsim.smpp.org',2775

import socket as s, struct as st
c=s.socket();c.connect((h,o))
c.sendall(st.pack(f"!IIII{len(i)+1}s{len(p)+1}s5B",16+len(i)+len(p)+7,2,0,1,i.encode(),p.encode(),0,0x34,1,1,0))
if c.recv(1024)[8:12]==b'\x00'*4:c.sendall(st.pack(f"!IIII3B{len(r)+1}s2B{len(d)+1}s9B1B{len(m)}s", 33+len(r)+len(d)+len(m),4,0,2,0,1,1,r.encode(),1,1,d.encode(),0,0,0,0,0,0,0,0,0, len(m),m.encode())); print("Message ID:",c.recv(1024)[16:].decode().split('\x00')[0])
c.close()
