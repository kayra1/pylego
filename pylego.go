package main

import (
	"C"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns"
	"github.com/go-acme/lego/v4/registration"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func request_certificate(email string, server string, csr []byte, plugin string) error {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: email,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)
	config.CADirURL = server

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	dns_provider, err := dns.NewDNSChallengeProviderByName(plugin)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Challenge.SetDNS01Provider(dns_provider)
	if err != nil {
		log.Fatal(err)
	}


	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	decoded_pem, _ := pem.Decode(csr)
	if decoded_pem == nil {
		log.Fatal("Could not decode PEM CSR")
	}
	parsed_csr, err := x509.ParseCertificateRequest(decoded_pem.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	request := certificate.ObtainForCSRRequest{
		CSR: parsed_csr,
		Bundle:  true,
	}
	certificates, err := client.Certificate.ObtainForCSR(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", certificates)

	return nil
}

//export run
func run(email *C.char, server *C.char, csr *C.char, plugin *C.char, env **C.char, env_len int) int {
	goemail := C.GoString(email)
	goserver := C.GoString(server)
	gocsr := []byte(C.GoString(csr))
	goplugin := C.GoString(plugin)
	env_slice := make([]string, 0, env_len)
	for _, e := range unsafe.Slice(env, env_len) {
		env_slice = append(env_slice, C.GoString(e))
	}
	for i := 0; i < env_len; i = i + 2 {
		os.Setenv(env_slice[i], env_slice[i+1])
	}
	err := request_certificate(goemail, goserver, gocsr, goplugin)
	if err != nil {
		return 1
	}
	return 0
}

func main() {}
