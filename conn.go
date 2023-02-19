package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"github.com/pingcap/errors"
	"net"
	"strings"
	"time"
)

const defaultAuthPluginName = AUTH_NATIVE_PASSWORD

type Conn struct {
	*Packet

	user      string
	password  string
	db        string
	tlsConfig *tls.Config
	proto     string

	// server capabilities
	capability uint32
	// client-set capabilities only
	ccaps      uint32
	attributes map[string]string

	charset        string
	status         uint16
	serverVersion  string
	connectionID   uint32
	salt           []byte
	authPluginName string
}

func (c *Conn) Ping() error {
	if err := c.writeCommand(COM_PING); err != nil {
		return errors.Trace(err)
	}

	if _, err := c.readOK(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) writeCommand(command byte) error {
	c.ResetSequence()

	return c.WritePacket([]byte{
		0x01, // 1 bytes long
		0x00,
		0x00,
		0x00, // sequence
		command,
	})
}

func (c *Conn) ResetSequence() {
	c.Sequence = 0
}

func (c *Conn) readOK() (*Result, error) {
	data, err := c.ReadPacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	if data[0] == OK_HEADER {
		return c.handleOKPacket(data)
	}

	if data[0] == ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	}

	return nil, errors.New("invalid ok packet")
}

func (c *Conn) handleOKPacket(data []byte) (*Result, error) {
	var n int
	var pos = 1

	r := new(Result)

	r.AffectedRows, _, n = LengthEncodedInt(data[pos:])
	pos += n
	r.InsertId, _, n = LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&CLIENT_PROTOCOL_41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2

		//todo:strict_mode, check warnings as error
		r.Warnings = binary.LittleEndian.Uint16(data[pos:])
		// pos += 2
	} else if c.capability&CLIENT_TRANSACTIONS > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		// pos += 2
	}

	//new ok package will check CLIENT_SESSION_TRACK too, but I don't support it now.

	//skip info
	return r, nil
}

func (c *Conn) handleErrorPacket(data []byte) error {
	e := new(MyError)

	var pos = 1

	// 2 bytes long, bytes 2 and 3
	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		// State 5 bytes long
		e.State = BytesToString(data[pos : pos+5])
		pos += 5
	}

	e.Message = BytesToString(data[pos:])

	return e
}

// Connect to a MySQL server, addr can be ip:port, or a unix socket domain like /var/sock.
// Accepts a series of configuration functions as a variadic argument.
func Connect(addr string, user string, password string, dbName string, options ...func(*Conn)) (*Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dialer := &net.Dialer{}

	return ConnectWithDialer(ctx, "", addr, user, password, dbName, dialer.DialContext, options...)
}

// Dialer connects to the address on the named network using the provided context.
type Dialer func(ctx context.Context, network, address string) (net.Conn, error)

// ConnectWithDialer Connect to a MySQL server using the given Dialer.
func ConnectWithDialer(ctx context.Context, network string, addr string, user string, password string, dbName string, dialer Dialer, options ...func(*Conn)) (*Conn, error) {
	c := new(Conn)

	if network == "" {
		network = getNetProto(addr)
	}

	var err error
	conn, err := dialer(ctx, network, addr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	c.user = user
	c.password = password
	c.db = dbName
	c.proto = network
	c.Packet = NewPacket(conn)

	// use default charset here, utf-8
	c.charset = DEFAULT_CHARSET

	// Apply configuration functions.
	for i := range options {
		options[i](c)
	}

	// if c.tlsConfig != nil {
	// 	seq := c.Packet.Sequence
	// 	c.Packet = NewTLSConn(conn)
	// 	c.Packet.Sequence = seq
	// }

	if err = c.handshake(); err != nil {
		return nil, errors.Trace(err)
	}

	return c, nil
}

func getNetProto(addr string) string {
	proto := "tcp"
	if strings.Contains(addr, "/") {
		proto = "unix"
	}
	return proto
}

func (c *Conn) handshake() error {
	var err error
	if err = c.readInitialHandshake(); err != nil {
		c.Close()
		return errors.Trace(err)
	}

	if err := c.writeAuthHandshake(); err != nil {
		c.Close()

		return errors.Trace(err)
	}

	if err := c.handleAuthResult(); err != nil {
		c.Close()
		return errors.Trace(err)
	}

	return nil
}

func (c *Conn) Close() error {
	return c.Packet.Close()
}

// See: http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake
func (c *Conn) readInitialHandshake() error {
	data, err := c.ReadPacket()
	if err != nil {
		return errors.Trace(err)
	}

	if data[0] == ERR_HEADER {
		return errors.Annotate(c.handleErrorPacket(data), "read initial handshake error")
	}

	if data[0] < MinProtocolVersion {
		return errors.Errorf("invalid protocol version %d, must >= 10", data[0])
	}

	// skip mysql version
	// mysql version end with 0x00
	version := data[1 : bytes.IndexByte(data[1:], 0x00)+1]
	c.serverVersion = string(version)
	pos := 1 + len(version)

	// connection id length is 4
	c.connectionID = binary.LittleEndian.Uint32(data[pos : pos+4])
	pos += 4

	c.salt = []byte{}
	c.salt = append(c.salt, data[pos:pos+8]...)

	// skip filter
	pos += 8 + 1

	// capability lower 2 bytes
	c.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))
	// check protocol
	if c.capability&CLIENT_PROTOCOL_41 == 0 {
		return errors.New("the MySQL server can not support protocol 41 and above required by the client")
	}
	if c.capability&CLIENT_SSL == 0 && c.tlsConfig != nil {
		return errors.New("the MySQL Server does not support TLS required by the client")
	}
	pos += 2

	if len(data) > pos {
		// skip server charset
		// c.charset = data[pos]
		pos += 1

		c.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
		// capability flags (upper 2 bytes)
		c.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | c.capability
		pos += 2

		// auth_data is end with 0x00, min data length is 13 + 8 = 21
		// ref to https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake
		maxAuthDataLen := 21
		if c.capability&CLIENT_PLUGIN_AUTH != 0 && int(data[pos]) > maxAuthDataLen {
			maxAuthDataLen = int(data[pos])
		}

		// skip reserved (all [00])
		pos += 10 + 1

		// auth_data is end with 0x00, so we need to trim 0x00
		resetOfAuthDataEndPos := pos + maxAuthDataLen - 8 - 1
		c.salt = append(c.salt, data[pos:resetOfAuthDataEndPos]...)

		// skip reset of end pos
		pos = resetOfAuthDataEndPos + 1

		if c.capability&CLIENT_PLUGIN_AUTH != 0 {
			c.authPluginName = string(data[pos : len(data)-1])
		}
	}

	// if server gives no default auth plugin name, use a client default
	if c.authPluginName == "" {
		c.authPluginName = defaultAuthPluginName
	}

	return nil
}

// See: http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::HandshakeResponse
func (c *Conn) writeAuthHandshake() error {
	if !authPluginAllowed(c.authPluginName) {
		return fmt.Errorf("unknow auth plugin name '%s'", c.authPluginName)
	}

	// Set default client capabilities that reflect the abilities of this library
	capability := CLIENT_PROTOCOL_41 | CLIENT_SECURE_CONNECTION |
		CLIENT_LONG_PASSWORD | CLIENT_TRANSACTIONS | CLIENT_PLUGIN_AUTH
	// Adjust client capability flags based on server support
	capability |= c.capability & CLIENT_LONG_FLAG
	// Adjust client capability flags on specific client requests
	// Only flags that would make any sense setting and aren't handled elsewhere
	// in the library are supported here
	capability |= c.ccaps&CLIENT_FOUND_ROWS | c.ccaps&CLIENT_IGNORE_SPACE |
		c.ccaps&CLIENT_MULTI_STATEMENTS | c.ccaps&CLIENT_MULTI_RESULTS |
		c.ccaps&CLIENT_PS_MULTI_RESULTS | c.ccaps&CLIENT_CONNECT_ATTRS

	// To enable TLS / SSL
	if c.tlsConfig != nil {
		capability |= CLIENT_SSL
	}

	auth, addNull, err := c.genAuthResponse(c.salt)
	if err != nil {
		return err
	}

	// encode length of the auth plugin data
	// here we use the Length-Encoded-Integer(LEI) as the data length may not fit into one byte
	// see: https://dev.mysql.com/doc/internals/en/integer.html#length-encoded-integer
	var authRespLEIBuf [9]byte
	authRespLEI := AppendLengthEncodedInteger(authRespLEIBuf[:0], uint64(len(auth)))
	if len(authRespLEI) > 1 {
		// if the length can not be written in 1 byte, it must be written as a
		// length encoded integer
		capability |= CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
	}

	// packet length
	// capability 4
	// max-packet size 4
	// charset 1
	// reserved all[0] 23
	// username
	// auth
	// mysql_native_password + null-terminated
	length := 4 + 4 + 1 + 23 + len(c.user) + 1 + len(authRespLEI) + len(auth) + 21 + 1
	if addNull {
		length++
	}
	// db name
	if len(c.db) > 0 {
		capability |= CLIENT_CONNECT_WITH_DB
		length += len(c.db) + 1
	}
	// connection attributes
	attrData := c.genAttributes()
	if len(attrData) > 0 {
		capability |= CLIENT_CONNECT_ATTRS
		length += len(attrData)
	}

	data := make([]byte, length+4)

	// capability [32 bit]
	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	// MaxPacketSize [32 bit] (none)
	data[8] = 0x00
	data[9] = 0x00
	data[10] = 0x00
	data[11] = 0x00

	// Charset [1 byte]
	// use default collation id 33 here, is utf-8
	data[12] = DEFAULT_COLLATION_ID

	// SSL Connection Request Packet
	// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::SSLRequest
	if c.tlsConfig != nil {
		// Send TLS / SSL request packet
		if err := c.WritePacket(data[:(4+4+1+23)+4]); err != nil {
			return err
		}

		// Switch to TLS
		tlsConn := tls.Client(c.Packet.Conn, c.tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			return err
		}

		currentSequence := c.Sequence
		c.Packet = NewPacket(tlsConn)
		c.Sequence = currentSequence
	}

	// Filler [23 bytes] (all 0x00)
	pos := 13
	for ; pos < 13+23; pos++ {
		data[pos] = 0
	}

	// User [null terminated string]
	if len(c.user) > 0 {
		pos += copy(data[pos:], c.user)
	}
	data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	pos += copy(data[pos:], authRespLEI)
	pos += copy(data[pos:], auth)
	if addNull {
		data[pos] = 0x00
		pos++
	}

	// db [null terminated string]
	if len(c.db) > 0 {
		pos += copy(data[pos:], c.db)
		data[pos] = 0x00
		pos++
	}

	// Assume native client during response
	pos += copy(data[pos:], c.authPluginName)
	data[pos] = 0x00
	pos++

	// connection attributes
	if len(attrData) > 0 {
		copy(data[pos:], attrData)
	}

	return c.WritePacket(data)
}

func (c *Conn) handleAuthResult() error {
	data, switchToPlugin, err := c.readAuthResult()
	if err != nil {
		return err
	}

	// handle auth switch, only support 'sha256_password', and 'caching_sha2_password'
	if switchToPlugin != "" {
		//fmt.Printf("now switching auth plugin to '%s'\n", switchToPlugin)
		if data == nil {
			data = c.salt
		} else {
			copy(c.salt, data)
		}
		c.authPluginName = switchToPlugin
		auth, addNull, err := c.genAuthResponse(data)
		if err != nil {
			return err
		}

		if err = c.WriteAuthSwitchPacket(auth, addNull); err != nil {
			return err
		}

		// Read Result Packet
		data, switchToPlugin, err = c.readAuthResult()
		if err != nil {
			return err
		}

		// Do not allow to change the auth plugin more than once
		if switchToPlugin != "" {
			return errors.Errorf("can not switch auth plugin more than once")
		}
	}

	// handle caching_sha2_password
	if c.authPluginName == AUTH_CACHING_SHA2_PASSWORD {
		if data == nil {
			return nil // auth already succeeded
		}
		if data[0] == CACHE_SHA2_FAST_AUTH {
			_, err = c.readOK()
			return err
		} else if data[0] == CACHE_SHA2_FULL_AUTH {
			// need full authentication
			if c.tlsConfig != nil || c.proto == "unix" {
				if err = c.WriteClearAuthPacket(c.password); err != nil {
					return err
				}
			} else {
				if err = c.WritePublicKeyAuthPacket(c.password, c.salt); err != nil {
					return err
				}
			}
			_, err = c.readOK()
			return err
		} else {
			return errors.Errorf("invalid packet %x", data[0])
		}
	} else if c.authPluginName == AUTH_SHA256_PASSWORD {
		if len(data) == 0 {
			return nil // auth already succeeded
		}
		block, _ := pem.Decode(data)
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return err
		}
		// send encrypted password
		err = c.WriteEncryptedPassword(c.password, c.salt, pub.(*rsa.PublicKey))
		if err != nil {
			return err
		}
		_, err = c.readOK()
		return err
	}
	return nil
}

// helper function to determine what auth methods are allowed by this client
func authPluginAllowed(pluginName string) bool {
	for _, p := range supportedAuthPlugins {
		if pluginName == p {
			return true
		}
	}
	return false
}

// defines the supported auth plugins
var supportedAuthPlugins = []string{AUTH_NATIVE_PASSWORD, AUTH_SHA256_PASSWORD, AUTH_CACHING_SHA2_PASSWORD}

// generate auth response data according to auth plugin
//
// NOTE: the returned boolean value indicates whether to add a \NUL to the end of data.
// it is quite tricky because MySQL server expects different formats of responses in different auth situations.
// here the \NUL needs to be added when sending back the empty password or cleartext password in 'sha256_password'
// authentication.
func (c *Conn) genAuthResponse(authData []byte) ([]byte, bool, error) {
	// password hashing
	switch c.authPluginName {
	case AUTH_NATIVE_PASSWORD:
		return CalcPassword(authData[:20], []byte(c.password)), false, nil
	case AUTH_CACHING_SHA2_PASSWORD:
		return CalcCachingSha2Password(authData, c.password), false, nil
	case AUTH_CLEAR_PASSWORD:
		return []byte(c.password), true, nil
	case AUTH_SHA256_PASSWORD:
		if len(c.password) == 0 {
			return nil, true, nil
		}
		if c.tlsConfig != nil || c.proto == "unix" {
			// write cleartext auth packet
			// see: https://dev.mysql.com/doc/refman/8.0/en/sha256-pluggable-authentication.html
			return []byte(c.password), true, nil
		} else {
			// request public key from server
			// see: https://dev.mysql.com/doc/internals/en/public-key-retrieval.html
			return []byte{1}, false, nil
		}
	default:
		// not reachable
		return nil, false, fmt.Errorf("auth plugin '%s' is not supported", c.authPluginName)
	}
}

// generate connection attributes data
func (c *Conn) genAttributes() []byte {
	if len(c.attributes) == 0 {
		return nil
	}

	attrData := make([]byte, 0)
	for k, v := range c.attributes {
		attrData = append(attrData, PutLengthEncodedString([]byte(k))...)
		attrData = append(attrData, PutLengthEncodedString([]byte(v))...)
	}
	return append(PutLengthEncodedInt(uint64(len(attrData))), attrData...)
}

func (c *Conn) readAuthResult() ([]byte, string, error) {
	data, err := c.ReadPacket()
	if err != nil {
		return nil, "", err
	}

	// see: https://insidemysql.com/preparing-your-community-connector-for-mysql-8-part-2-sha256/
	// packet indicator
	switch data[0] {
	case OK_HEADER:
		_, err := c.handleOKPacket(data)
		return nil, "", err

	case MORE_DATE_HEADER:
		return data[1:], "", err

	case EOF_HEADER:
		// server wants to switch auth
		if len(data) < 1 {
			// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::OldAuthSwitchRequest
			return nil, AUTH_MYSQL_OLD_PASSWORD, nil
		}
		pluginEndIndex := bytes.IndexByte(data, 0x00)
		if pluginEndIndex < 0 {
			return nil, "", errors.New("invalid packet")
		}
		plugin := string(data[1:pluginEndIndex])
		authData := data[pluginEndIndex+1:]
		return authData, plugin, nil

	default: // Error otherwise
		return nil, "", c.handleErrorPacket(data)
	}
}

// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::AuthSwitchResponse
func (c *Conn) WriteAuthSwitchPacket(authData []byte, addNUL bool) error {
	pktLen := 4 + len(authData)
	if addNUL {
		pktLen++
	}
	data := make([]byte, pktLen)

	// Add the auth data [EOF]
	copy(data[4:], authData)
	if addNUL {
		data[pktLen-1] = 0x00
	}

	return errors.Wrap(c.WritePacket(data), "WritePacket failed")
}

// WriteClearAuthPacket: Client clear text authentication packet
// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::AuthSwitchResponse
func (c *Conn) WriteClearAuthPacket(password string) error {
	// Calculate the packet length and add a tailing 0
	pktLen := len(password) + 1
	data := make([]byte, 4+pktLen)

	// Add the clear password [null terminated string]
	copy(data[4:], password)
	data[4+pktLen-1] = 0x00

	return errors.Wrap(c.WritePacket(data), "WritePacket failed")
}

// WritePublicKeyAuthPacket: Caching sha2 authentication. Public key request and send encrypted password
// http://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::AuthSwitchResponse
func (c *Conn) WritePublicKeyAuthPacket(password string, cipher []byte) error {
	// request public key
	data := make([]byte, 4+1)
	data[4] = 2 // cachingSha2PasswordRequestPublicKey
	if err := c.WritePacket(data); err != nil {
		return errors.Wrap(err, "WritePacket(single byte) failed")
	}

	data, err := c.ReadPacket()
	if err != nil {
		return errors.Wrap(err, "ReadPacket failed")
	}

	block, _ := pem.Decode(data[1:])
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "x509.ParsePKIXPublicKey failed")
	}

	plain := make([]byte, len(password)+1)
	copy(plain, password)
	for i := range plain {
		j := i % len(cipher)
		plain[i] ^= cipher[j]
	}
	sha1v := sha1.New()
	enc, _ := rsa.EncryptOAEP(sha1v, rand.Reader, pub.(*rsa.PublicKey), plain, nil)
	data = make([]byte, 4+len(enc))
	copy(data[4:], enc)
	return errors.Wrap(c.WritePacket(data), "WritePacket failed")
}

func (c *Conn) WriteEncryptedPassword(password string, seed []byte, pub *rsa.PublicKey) error {
	enc, err := EncryptPassword(password, seed, pub)
	if err != nil {
		return errors.Wrap(err, "EncryptPassword failed")
	}
	return errors.Wrap(c.WriteAuthSwitchPacket(enc, false), "WriteAuthSwitchPacket failed")
}

func (c *Conn) Execute(command string, args ...interface{}) (*Result, error) {
	if len(args) == 0 {
		return c.exec(command)
	}

	s, err := c.Prepare(command)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var r *Result
	r, err = s.Execute(args...)
	s.Close()
	return r, err
}

func (c *Conn) exec(query string) (*Result, error) {
	if err := c.writeCommandStr(COM_QUERY, query); err != nil {
		return nil, errors.Trace(err)
	}

	return c.readResult(false)
}

func (c *Conn) writeCommandStr(command byte, arg string) error {
	return c.writeCommandBuf(command, StringToBytes(arg))
}

func (c *Conn) writeCommandBuf(command byte, arg []byte) error {
	c.ResetSequence()

	length := len(arg) + 1
	data := GetByteSlice(length + 4)
	data.B[4] = command

	copy(data.B[5:], arg)

	err := c.WritePacket(data.B)

	PutByteSlice(data)

	return err
}

func (c *Conn) readResult(binary bool) (*Result, error) {
	bs := GetByteSlice(16)
	defer PutByteSlice(bs)
	var err error
	bs.B, err = c.ReadPacketReuseMem(bs.B[:0])
	if err != nil {
		return nil, errors.Trace(err)
	}

	switch bs.B[0] {
	case OK_HEADER:
		return c.handleOKPacket(bs.B)
	case ERR_HEADER:
		return nil, c.handleErrorPacket(bytes.Repeat(bs.B, 1))
	case LocalInFile_HEADER:
		return nil, ErrMalformPacket
	default:
		return c.readResultset(bs.B, binary)
	}
}

func (c *Conn) readResultset(data []byte, binary bool) (*Result, error) {
	// column count
	count, _, n := LengthEncodedInt(data)

	if n-len(data) != 0 {
		return nil, ErrMalformPacket
	}

	result := &Result{
		Resultset: NewResultset(int(count)),
	}

	if err := c.readResultColumns(result); err != nil {
		return nil, errors.Trace(err)
	}

	if err := c.readResultRows(result, binary); err != nil {
		return nil, errors.Trace(err)
	}

	return result, nil
}

func (c *Conn) readResultColumns(result *Result) (err error) {
	var i = 0
	var data []byte

	for {
		rawPkgLen := len(result.RawPkg)
		result.RawPkg, err = c.ReadPacketReuseMem(result.RawPkg)
		if err != nil {
			return
		}
		data = result.RawPkg[rawPkgLen:]

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&CLIENT_PROTOCOL_41 > 0 {
				result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			if i != len(result.Fields) {
				err = ErrMalformPacket
			}

			return
		}

		if result.Fields[i] == nil {
			result.Fields[i] = &Field{}
		}
		err = result.Fields[i].Parse(data)
		if err != nil {
			return
		}

		result.FieldNames[BytesToString(result.Fields[i].Name)] = i

		i++
	}
}

func (c *Conn) readResultRows(result *Result, isBinary bool) (err error) {
	var data []byte

	for {
		rawPkgLen := len(result.RawPkg)
		result.RawPkg, err = c.ReadPacketReuseMem(result.RawPkg)
		if err != nil {
			return
		}
		data = result.RawPkg[rawPkgLen:]

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&CLIENT_PROTOCOL_41 > 0 {
				result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			break
		}

		if data[0] == ERR_HEADER {
			return c.handleErrorPacket(data)
		}

		result.RowDatas = append(result.RowDatas, data)
	}

	if cap(result.Values) < len(result.RowDatas) {
		result.Values = make([][]FieldValue, len(result.RowDatas))
	} else {
		result.Values = result.Values[:len(result.RowDatas)]
	}

	for i := range result.Values {
		result.Values[i], err = result.RowDatas[i].Parse(result.Fields, isBinary, result.Values[i])

		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (c *Conn) isEOFPacket(data []byte) bool {
	return data[0] == EOF_HEADER && len(data) <= 5
}

func (c *Conn) Prepare(query string) (*Stmt, error) {
	if err := c.writeCommandStr(COM_STMT_PREPARE, query); err != nil {
		return nil, errors.Trace(err)
	}

	data, err := c.ReadPacket()
	if err != nil {
		return nil, errors.Trace(err)
	}

	if data[0] == ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	}
	if data[0] != OK_HEADER {
		return nil, ErrMalformPacket
	}

	s := new(Stmt)
	s.conn = c

	pos := 1

	//for statement id
	s.id = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	//number columns
	s.columns = int(binary.LittleEndian.Uint16(data[pos:]))
	pos += 2

	//number params
	s.params = int(binary.LittleEndian.Uint16(data[pos:]))
	pos += 2

	//warnings
	s.warnings = int(binary.LittleEndian.Uint16(data[pos:]))
	// pos += 2

	if s.params > 0 {
		if err := s.conn.readUntilEOF(); err != nil {
			return nil, errors.Trace(err)
		}
	}

	if s.columns > 0 {
		if err := s.conn.readUntilEOF(); err != nil {
			return nil, errors.Trace(err)
		}
	}

	return s, nil
}

func (c *Conn) readUntilEOF() (err error) {
	var data []byte

	for {
		data, err = c.ReadPacket()

		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			return
		}
	}
}

func (c *Conn) writeCommandUint32(command byte, arg uint32) error {
	c.ResetSequence()

	return c.WritePacket([]byte{
		0x05, //5 bytes long
		0x00,
		0x00,
		0x00, //sequence

		command,

		byte(arg),
		byte(arg >> 8),
		byte(arg >> 16),
		byte(arg >> 24),
	})
}

func (c *Conn) readResultStreaming(binary bool, result *Result, perRowCb SelectPerRowCallback, perResCb SelectPerResultCallback) error {
	bs := GetByteSlice(16)
	defer PutByteSlice(bs)
	var err error
	bs.B, err = c.ReadPacketReuseMem(bs.B[:0])
	if err != nil {
		return errors.Trace(err)
	}

	switch bs.B[0] {
	case OK_HEADER:
		// https://dev.mysql.com/doc/internals/en/com-query-response.html
		// 14.6.4.1 COM_QUERY Response
		// If the number of columns in the resultset is 0, this is a OK_Packet.

		okResult, err := c.handleOKPacket(bs.B)
		if err != nil {
			return errors.Trace(err)
		}

		result.Status = okResult.Status
		result.AffectedRows = okResult.AffectedRows
		result.InsertId = okResult.InsertId
		result.Warnings = okResult.Warnings
		if result.Resultset == nil {
			result.Resultset = NewResultset(0)
		} else {
			result.Reset(0)
		}
		return nil
	case ERR_HEADER:
		return c.handleErrorPacket(bytes.Repeat(bs.B, 1))
	case LocalInFile_HEADER:
		return ErrMalformPacket
	default:
		return c.readResultsetStreaming(bs.B, binary, result, perRowCb, perResCb)
	}
}

func (c *Conn) readResultsetStreaming(data []byte, binary bool, result *Result, perRowCb SelectPerRowCallback, perResCb SelectPerResultCallback) error {
	columnCount, _, n := LengthEncodedInt(data)

	if n-len(data) != 0 {
		return ErrMalformPacket
	}

	if result.Resultset == nil {
		result.Resultset = NewResultset(int(columnCount))
	} else {
		// Reuse memory if can
		result.Reset(int(columnCount))
	}

	// this is a streaming resultset
	result.Resultset.Streaming = StreamingSelect

	if err := c.readResultColumns(result); err != nil {
		return errors.Trace(err)
	}

	if perResCb != nil {
		if err := perResCb(result); err != nil {
			return err
		}
	}

	if err := c.readResultRowsStreaming(result, binary, perRowCb); err != nil {
		return errors.Trace(err)
	}

	// this resultset is done streaming
	result.Resultset.StreamingDone = true

	return nil
}

func (c *Conn) readResultRowsStreaming(result *Result, isBinary bool, perRowCb SelectPerRowCallback) (err error) {
	var (
		data []byte
		row  []FieldValue
	)

	for {
		data, err = c.ReadPacketReuseMem(data[:0])
		if err != nil {
			return
		}

		// EOF Packet
		if c.isEOFPacket(data) {
			if c.capability&CLIENT_PROTOCOL_41 > 0 {
				result.Warnings = binary.LittleEndian.Uint16(data[1:])
				// todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				c.status = result.Status
			}

			break
		}

		if data[0] == ERR_HEADER {
			return c.handleErrorPacket(data)
		}

		// Parse this row
		row, err = RowData(data).Parse(result.Fields, isBinary, row)
		if err != nil {
			return errors.Trace(err)
		}

		// Send the row to "userland" code
		err = perRowCb(row)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (c *Conn) Begin() error {
	_, err := c.exec("BEGIN")
	return errors.Trace(err)
}

func (c *Conn) Commit() error {
	_, err := c.exec("COMMIT")
	return errors.Trace(err)
}

func (c *Conn) Rollback() error {
	_, err := c.exec("ROLLBACK")
	return errors.Trace(err)
}

func (c *Conn) HandleErrorPacket(data []byte) error {
	return c.handleErrorPacket(data)
}

func (c *Conn) ReadOKPacket() (*Result, error) {
	return c.readOK()
}

func (c *Conn) SetCharset(charset string) error {
	if c.charset == charset {
		return nil
	}

	_, err := c.exec(fmt.Sprintf("SET NAMES %s", charset))
	if err != nil {
		return errors.Trace(err)
	}

	c.charset = charset
	return nil
}

func (c *Conn) GetConnectionID() uint32 {
	return c.connectionID
}

// SetTLSConfig: use user-specified TLS config
// pass to options when connect
func (c *Conn) SetTLSConfig(config *tls.Config) {
	c.tlsConfig = config
}
