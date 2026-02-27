package tls

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "net"
    "os"
    "time"
)

// GenerateSelfSignedCert generates a new ECDSA certificate and corresponding
// private key suitable for use with a development TLS server. The cert will
// be valid for the provided hostnames and IPs and will be written to the
// provided file paths in PEM format. Existing files will be overwritten.
func GenerateSelfSignedCert(certPath, keyPath string, hosts []string) error {
    priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        return err
    }

    notBefore := time.Now()
    notAfter := notBefore.Add(365 * 24 * time.Hour)

    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        return err
    }

    tmpl := x509.Certificate{
        SerialNumber: serialNumber,
        Subject: pkix.Name{
            Organization: []string{"Evolutionary MCP Dev"},
        },
        NotBefore: notBefore,
        NotAfter:  notAfter,
        KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage: []x509.ExtKeyUsage{
            x509.ExtKeyUsageServerAuth,
        },
        BasicConstraintsValid: true,
    }

    for _, h := range hosts {
        if ip := net.ParseIP(h); ip != nil {
            tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
        } else {
            tmpl.DNSNames = append(tmpl.DNSNames, h)
        }
    }

    derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
    if err != nil {
        return err
    }

    certOut, err := os.Create(certPath)
    if err != nil {
        return err
    }
    defer certOut.Close()
    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

    keyOut, err := os.Create(keyPath)
    if err != nil {
        return err
    }
    defer keyOut.Close()
    b, err := x509.MarshalECPrivateKey(priv)
    if err != nil {
        return err
    }
    pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})

    return nil
}
