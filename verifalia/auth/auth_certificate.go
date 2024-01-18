package auth

/*
* Verifalia - Email list cleaning and real-time email verification service
* https://verifalia.com/
* support@verifalia.com
*
* Copyright (c) 2005-2024 Cobisi Research
*
* Cobisi Research
* Via Della Costituzione, 31
* 35010 Vigonza
* Italy - European Union
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in
* all copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
* THE SOFTWARE.
 */

import (
	"crypto/tls"
	"net/http"
	"time"
)

// X.509 mutual TLS client certificate authentication

type certificateAuthProvider struct {
	Certificate *tls.Certificate
	Client      *http.Client
}

func (provider certificateAuthProvider) Authenticate(request *http.Request) error {
	return nil
}

func (provider certificateAuthProvider) HandleUnauthorizedRequest() error {
	return nil
}

func (provider certificateAuthProvider) BuildClient() *http.Client {
	return provider.Client
}

func NewCertificateAuthProvider(certificate *tls.Certificate) Provider {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*certificate},
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &certificateAuthProvider{
		Certificate: certificate,
		Client:      client,
	}
}
