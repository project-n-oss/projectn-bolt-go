package boltv4

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws/request"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"net/http"
	"time"
)

// SignBoltRequestHandler is a named request handler the Bolt Go SDK will use to sign
// service client request as a STS GetCallerIdentity Request using the V4 signature.
var SignBoltRequestHandler = request.NamedHandler{
	Name: "boltv4.SignBoltRequestHandler", Fn: SignBoltSDKRequest,
}

var (
	awsSTSRequestBody = []byte("Action=GetCallerIdentity&Version=2011-06-15")
	stsEndpointUrl = "https://sts.amazonaws.com/"
)

const (
	authorizationHeader = "Authorization"
	amzSecurityTokenHeader = "X-Amz-Security-Token"
	amzDateHeader = "X-Amz-Date"
)

// SignBoltSDKRequest signs an AWS request as a STS GetCallerIdentity request using V4 signature.
func SignBoltSDKRequest(req *request.Request) {

	v4signer := v4.NewSigner(req.Config.Credentials)

	// create STS GetCallerIdentity Request
	stsReq, err := http.NewRequest(http.MethodPost, stsEndpointUrl, nil)
	if err != nil {
		req.Error = err
		req.SignedHeaderVals = nil
		return
	}

	// signs the STS GetCallerIdentity Request using v4 signature.
	curTime := time.Now()
	signedHeaders, err := v4signer.Sign(stsReq, bytes.NewReader(awsSTSRequestBody), "sts", "us-east-1", curTime)
	if err != nil {
		req.Error = err
		req.SignedHeaderVals = nil
		return
	}

	// Add headers from the signed STS GetCallerIdentity Request to the original request.
	if secTokenHeader := stsReq.Header.Get(amzSecurityTokenHeader); len(secTokenHeader) > 0 {
		req.HTTPRequest.Header.Set(amzSecurityTokenHeader, secTokenHeader)
	}

	if dateHeader := stsReq.Header.Get(amzDateHeader); len(dateHeader) > 0 {
		req.HTTPRequest.Header.Set(amzDateHeader, dateHeader)
	}

	if authHeader := stsReq.Header.Get(authorizationHeader); len(authHeader) > 0 {
		req.HTTPRequest.Header.Set(authorizationHeader, authHeader)
	}

	req.SignedHeaderVals = signedHeaders
	req.LastSignedAt = curTime
}
