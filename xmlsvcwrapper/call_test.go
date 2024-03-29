package xmlsvcwrapper

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestRequest_Call(t *testing.T) {
	type fields struct {
		Url            string
		Header         http.Header
		client         *Client
		Ctx            context.Context
		SoapEnv        string
		SoapType       string
		BodyType       string
		BodyContent    string
		payloadRequest []byte
		SecurityHeader struct {
			UserName   string
			Password   string
			TimeToLive int64
		}
	}

	// Create dummy service
	go createMockWebService()

	headers := http.Header{}
	headers.Set("Content-Type", "text/xml; charset=utf-8")

	client := New()

	pRequest := `<soapenv:Envelope xmlns:soapenv="" xmlns:type=""><soapenv:Header><wsse:Security soapenv:mustUnderstand="1" xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"><wsu:Timestamp wsu:Id="Timestamp-20046406" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"><wsu:Created>2023-08-18T20:22:17Z</wsu:Created><wsu:Expires>2023-08-18T20:22:17Z</wsu:Expires></wsu:Timestamp><wsse:UsernameToken wsu:Id="UsernameToken-20914066" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"><wsse:Username></wsse:Username><wsse:Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordDigest">c9A2Uai15WulXRScmclx0mfpRRE=</wsse:Password><wsse:Nonce>UnC0v9MBo5g5Z/KPYB5GNA==</wsse:Nonce><wsu:Created>2023-08-18T20:22:17Z</wsu:Created></wsse:UsernameToken></wsse:Security></soapenv:Header><soapenv:Body><Content></Content></soapenv:Body></soapenv:Envelope>`

	responseText := `<soapenv:Envelope xmlns:soapenv=""><soapenv:Header></soapenv:Header><soapenv:Body>
				  <Response>
					 <string>Hello World!</string>
				  </Response>
			   </soapenv:Body></soapenv:Envelope>`

	testFields := fields{
		Url:    "http://127.0.0.1:3000/test",
		Header: headers,
		client: client,
		Ctx:    context.Background(),
	}

	test2Fields := fields{
		Url:    "http://127.0.0.3000:3000/test",
		Header: headers,
		client: client,
		Ctx:    context.Background(),
	}

	req := &Request{
		Url:            "http://127.0.0.1:3000/test",
		Header:         headers,
		client:         client,
		Ctx:            context.Background(),
		payloadRequest: []byte(pRequest),
	}

	resp := &Response{
		Request:         req,
		RawResponse:     nil,
		payloadResponse: []byte(responseText),
	}

	tests := []struct {
		name    string
		fields  fields
		want    *Response
		wantErr bool
	}{
		{
			name:    "Test Call()",
			fields:  testFields,
			want:    resp,
			wantErr: false,
		},
		{
			name:    "Test Request with wrong url",
			fields:  test2Fields,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Url:            tt.fields.Url,
				Header:         tt.fields.Header,
				client:         tt.fields.client,
				Ctx:            tt.fields.Ctx,
				SoapEnv:        tt.fields.SoapEnv,
				SoapType:       tt.fields.SoapType,
				BodyType:       tt.fields.BodyType,
				BodyContent:    tt.fields.BodyContent,
				payloadRequest: tt.fields.payloadRequest,
				SecurityHeader: tt.fields.SecurityHeader,
			}
			got, err := r.Call()
			if (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !reflect.DeepEqual(got.PayloadResponse(), tt.want.payloadResponse) {
					t.Errorf("Call() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
