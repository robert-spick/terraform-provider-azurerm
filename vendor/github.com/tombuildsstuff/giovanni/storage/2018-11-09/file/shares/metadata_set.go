package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

// SetMetaData sets the MetaData on the specified Storage Share
func (client Client) SetMetaData(ctx context.Context, accountName, shareName string, metaData map[string]string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "SetMetaData", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "SetMetaData", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "SetMetaData", "`shareName` must be a lower-cased string.")
	}
	if err := metadata.Validate(metaData); err != nil {
		return result, validation.NewError("shares.Client", "SetMetaData", fmt.Sprintf("`metadata` is not valid: %s.", err))
	}

	req, err := client.SetMetaDataPreparer(ctx, accountName, shareName, metaData)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "SetMetaData", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetMetaDataSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "SetMetaData", resp, "Failure sending request")
		return
	}

	result, err = client.SetMetaDataResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "SetMetaData", resp, "Failure responding to request")
		return
	}

	return
}

// SetMetaDataPreparer prepares the SetMetaData request.
func (client Client) SetMetaDataPreparer(ctx context.Context, accountName, shareName string, metaData map[string]string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("query", "share"),
		"comp":    autorest.Encode("query", "metadata"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}
	headers = metadata.SetIntoHeaders(headers, metaData)

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetPropertiesSetMetaDataSender sends the SetMetaData request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetMetaDataSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetMetaDataResponder handles the response to the SetMetaData request. The method always
// closes the http.Response Body.
func (client Client) SetMetaDataResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
