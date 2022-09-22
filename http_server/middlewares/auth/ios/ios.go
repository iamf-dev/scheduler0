package ios

import (
	"net/http"
	"scheduler0/http_server/middlewares/auth"
	"scheduler0/service"
	"scheduler0/utils"
)

// IsIOSClient returns true is the request is coming from an ios app
func IsIOSClient(req *http.Request) bool {
	apiKey := req.Header.Get(auth.APIKeyHeader)
	bundleID := req.Header.Get(auth.IOSBundleHeader)
	return len(apiKey) > 9 && len(bundleID) > 9
}

// IsAuthorizedIOSClient returns true if the credential is authorized ios app
func IsAuthorizedIOSClient(req *http.Request, credentialService service.Credential) (bool, *utils.GenericError) {
	apiKey := req.Header.Get(auth.APIKeyHeader)
	IOSBundleID := req.Header.Get(auth.IOSBundleHeader)

	return credentialService.ValidateIOSAPIKey(apiKey, IOSBundleID)
}