# https://github.com/melroselabs/mlsmpputils

system_id, password, host, port = "SYSTEMID", "PASSWORD", 'smscsim.smpp.org', 2775

import socket as s, struct as st

def create_pdu(command_id, sequence_number, message=b''):
    pdu_header = st.pack("!IIII", 16 + len(message), command_id, 0, sequence_number)
    return pdu_header + message

with s.socket() as cs:
    cs.connect((host, port))
    bi = (system_id.encode() + b'\x00', password.encode() + b'\x00', b'\x00\x34\x01\x01\x00')
    cs.sendall(create_pdu(0x00000001, 1, b''.join(bi)))
    if cs.recv(1024)[8:12] == b'\x00' * 4:
        print("Successfully bound as receiver.")
        try:
            while True:
                pdu = cs.recv(1024)
                if pdu and len(pdu) >= 16:
                    command_id = st.unpack('!I', pdu[4:8])[0]
                    sequence_number = st.unpack('!I', pdu[12:16])[0]
                    if command_id == 5:
                        ml = st.unpack('!I', pdu[:4])[0] - 16
                        message = pdu[16:16+ml].decode()
                        print("Received message:", message)
                        cs.sendall(create_pdu(0x80000005, sequence_number,b'\0'))
                    elif command_id == 21:
                        cs.sendall(create_pdu(0x80000015, sequence_number))
        except KeyboardInterrupt:
            print("Stopping receiver...")
    else:
        print("Bind failed.")
