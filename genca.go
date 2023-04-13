package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
)

func GenerateSKI(pub *rsa.PublicKey) []byte {
	var Raw_module []byte
	header := []byte{0x30, 0x81, 0x89, 0x02, 0x81, 0x81, 0x00}
	footer := []byte{0x02, 0x03, 0x01, 0x00, 0x01}
	module := pub.N.Bytes()

	sha1_hash := sha1.New()
	sha1_hash.Write(header)
	sha1_hash.Write(module)
	sha1_hash.Write(footer)
	sha1_hash.Write(Raw_module)

	return sha1_hash.Sum(nil)

}

func genkey(key_name string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
		return
	}
	keyOut, err := os.OpenFile(key_name+".pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("Failed to open "+key_name+".pem for writing:", err)
		return
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("Written " + key_name + ".pem\n")

	var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

	emailAddress := "client@nifcloud.local"
	subj := pkix.Name{
		CommonName:         key_name,
		Organization:       []string{key_name},
		Country:            []string{"JP"},
		Province:           []string{"Tokyo"},
		Locality:           []string{"Chuo-ku"},
		OrganizationalUnit: []string{"NIFCloud Private CA"},
	}
	rawSubj := subj.ToRDNSequence()
	rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
		{Type: oidEmailAddress, Value: emailAddress},
	})

	asn1Subj, _ := asn1.Marshal(rawSubj)

	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     []string{emailAddress},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, priv)

	crt_keyOut, err := os.OpenFile(key_name+".csr.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("Failed to open "+key_name+".csr.pem for writing:", err)
		return
	}
	pem.Encode(crt_keyOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	crt_keyOut.Close()
	log.Print("Written " + key_name + ".csr.pem\n")

	now := time.Now()
	var years = 4
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	cert_template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   key_name,
			Organization: []string{key_name},
			Country:      []string{"JP"},
		},
		NotBefore:             now.Add(-5 * time.Minute).UTC(),
		NotAfter:              now.AddDate(years, 0, 0).UTC(), // valid for years
		BasicConstraintsValid: true,
		IsCA:         false,
		SubjectKeyId: GenerateSKI(&priv.PublicKey),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	crt_Bytes, err := x509.CreateCertificate(
		rand.Reader, &cert_template, &cert_template, &priv.PublicKey, priv)

	if err != nil {
		log.Fatalf("Failed to create client Certificate: %s", err)
		return
	}
	client_certout, err := os.Create(key_name + ".crt.pem")
	if err != nil {
		log.Fatalf("Failed to open "+key_name+".crt.pem for writing: %s", err)
		return
	}
	pem.Encode(client_certout, &pem.Block{Type: "CERTIFICATE", Bytes: crt_Bytes})
	client_certout.Close()
	log.Print("Written " + key_name + ".crt.pem\n")

}

func readPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyData)
	if privateKeyBlock == nil {
		return nil, errors.New("invalid private key data")
	}
	if privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return nil, errors.New(fmt.Sprintf("invalid private key type : %s", privateKeyBlock.Type))
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, err
}

func readPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	if publicKeyBlock == nil {
		return nil, errors.New("invalid public key data")
	}
	if publicKeyBlock.Type != "PUBLIC KEY" {
		return nil, errors.New(fmt.Sprintf("invalid public key type : %s", publicKeyBlock.Type))
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return publicKey, nil
}

func readCertificateByte(path string) (*pem.Block, error) {
	certificateData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certificateBlock, _ := pem.Decode(certificateData)
	if certificateBlock == nil {
		return nil, errors.New("invalid certificate data")
	}
	if certificateBlock.Type != "CERTIFICATE" {
		return nil, errors.New(fmt.Sprintf("invalid certificate type : %s", certificateBlock.Type))
	}

	return certificateBlock, nil
}

func readCertificate(path string) (*x509.Certificate, error) {
	certificateBlock, err := readCertificateByte(path)
	certificate, err := x509.ParseCertificate(certificateBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func readCertificateRequest(path string) (*x509.CertificateRequest, error) {
	certificateRequestData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certificateRequestBlock, _ := pem.Decode(certificateRequestData)
	if certificateRequestBlock == nil {
		return nil, errors.New("invalid certificate request data")
	}
	if certificateRequestBlock.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New(fmt.Sprintf("invalid certificate request type : %s", certificateRequestBlock.Type))
	}

	certificaterequest, err := x509.ParseCertificateRequest(certificateRequestBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return certificaterequest, nil
}

func self_signed(crt_name string, ca_key_name string, ca_cert_name string, client_csr_name string) {
	ca_key, err := readPrivateKey(ca_key_name)
	if err != nil {
		fmt.Println(err)
	}
	ca_cert, err := readCertificate(ca_cert_name)
	client_csr, err := readCertificateRequest(client_csr_name)
	if err != nil {
		fmt.Println(err)
	}
	csr_pub, ok := client_csr.PublicKey.(*rsa.PublicKey)
	if !ok {
		errors.New("not RSA public key")
		return
	}

	now := time.Now()
	var years = 4
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	ComName := crt_name
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   ComName,
			Organization: []string{"nifcloud.local"},
			Country:      []string{"JP"},
		},
		NotBefore:             now.Add(-5 * time.Minute).UTC(),
		NotAfter:              now.AddDate(years, 0, 0).UTC(), // valid for years
		BasicConstraintsValid: true,
		IsCA:         false,
		SubjectKeyId: GenerateSKI(csr_pub),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
                ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
	crt_Bytes, err := x509.CreateCertificate(
		rand.Reader, &template, ca_cert, csr_pub, ca_key)

	if err != nil {
		log.Fatalf("Failed to create client Certificate: %s", err)
		return
	}
	client_certout, err := os.Create(crt_name + ".signed.crt.pem")
	if err != nil {
		log.Fatalf("Failed to open "+crt_name+".signed.crt.pem for writing: %s", err)
		return
	}
	pem.Encode(client_certout, &pem.Block{Type: "CERTIFICATE", Bytes: crt_Bytes})
	client_certout.Close()
	log.Print("Written " + crt_name + ".signed.crt.pem\n")

}


func genCACert(caName string, years int) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
		return
	}

	now := time.Now()
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   caName,
			Organization: []string{caName},
			Country:      []string{"JP"},
		},
		NotBefore:             now.Add(-5 * time.Minute).UTC(),
		NotAfter:              now.AddDate(years, 0, 0).UTC(), // valid for years
		BasicConstraintsValid: true,
		IsCA:           true,
		SubjectKeyId:   GenerateSKI(&priv.PublicKey),
		AuthorityKeyId: GenerateSKI(&priv.PublicKey),
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader, &template, &template, &priv.PublicKey, priv)

	if err != nil {
		log.Fatalf("Failed to create CA Certificate: %s", err)
		return
	}

	certOut, err := os.Create(caName + ".CAcert.pem")
	if err != nil {
		log.Fatalf("Failed to open "+caName+".CAcert.pem for writing: %s", err)
		return
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	certOut.Close()
	log.Print("Written " + caName + ".CAcert.pem\n")

	keyOut, err := os.OpenFile(caName+".CAkey.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("Failed to open "+caName+".CAkey.pem for writing:", err)
		return
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("Written " + caName + ".CAkey.pem\n")
}

func showCert(caName string) {
	certIn, err := ioutil.ReadFile(caName + ".CAcert.pem")
	if err != nil {
		log.Fatalf("Failed to open "+caName+".CAcert.pem for reading: %s", err)
		return
	}
	b, _ := pem.Decode(certIn)
	if b == nil {
		log.Fatalf("Failed to find a certificate in " + caName + ".CAcert.pem")
		return
	}
	caCert, err := x509.ParseCertificate(b.Bytes)
	log.Print("Readed cert "+caCert.Subject.CommonName+" isCA ", caCert.IsCA, "\n")
}

func main() {
	// Create CA.
	domain := "nifcloud.local"
	genCACert(domain, 4)
	showCert(domain)

	name := "client." + domain

	//Create client Private key , CSR.
	genkey(name)

	// Self Sign.
	self_signed(name, domain+".CAkey.pem", domain+".CAcert.pem", name+".csr.pem")

	name = "server." + domain

	//Create Server Cert Private key
	genkey(name)
	// Self Sign.
	self_signed(name, domain+".CAkey.pem", domain+".CAcert.pem", name+".csr.pem")
}
