package bsw

import (
	"crypto/tls"
	"net"
	"time"
)

// TLS attempts connection to an IP using TLS on port 443, and if successfull, will parse the server
// certificate for CommonName and SubjectAlt names.
func TLS(ip string, timeout int64) (Results, error) {
	results := Results{}
	tconn, err := net.Dial("tcp", ip+":443")
	if err != nil {
		return results, err
	}
	t := time.Duration(timeout)*time.Millisecond
	if err := tconn.SetDeadline(time.Now().Add(t)); err != nil {
		return results, err		
	}
	conn := tls.Client(tconn, &tls.Config{InsecureSkipVerify: true})
	defer conn.Close()
	if err := conn.Handshake(); err != nil {
		return results, err
	}
	state := conn.ConnectionState()
	cert := state.PeerCertificates[0]
	results = append(results, Result{Source: "TLS Certificate", IP: ip, Hostname: cert.Subject.CommonName})
	for _, name := range cert.DNSNames {
		results = append(results, Result{Source: "TLS Certificate", IP: ip, Hostname: name})
	}
	return results, nil
}
