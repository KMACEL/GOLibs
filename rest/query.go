package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/KMACEL/IITR/errc"
)

/*
 ██████╗ ██╗   ██╗███████╗██████╗ ██╗   ██╗
██╔═══██╗██║   ██║██╔════╝██╔══██╗╚██╗ ██╔╝
██║   ██║██║   ██║█████╗  ██████╔╝ ╚████╔╝
██║▄▄ ██║██║   ██║██╔══╝  ██╔══██╗  ╚██╔╝
╚██████╔╝╚██████╔╝███████╗██║  ██║   ██║
 ╚══▀▀═╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝   ╚═╝
*/

// Query is performs query util. These util are "Post", "Get", "Put".
// With this class, it makes a stable query everywhere.
// Retrieve values by query type. These values will be explained in the description of each questionnaire.
// The queries were taken from https://api.ardich.com/api/v3/apidoc/
// The queries return two values.
//    1: [] byte:
//           If the incoming data is meaningful, that is, "Status" as the response,
//           if "200 OK" is reached, the data that you read is sent as a byte array to the querying function.
//           If the message is empty or is problem, it returns the value "nil".
//    2: error:
//           If it encounters an error during the query, it sends back an error.
type Query struct{}

/*
 ██████╗ ███████╗████████╗         ██████╗ ██╗   ██╗███████╗██████╗ ██╗   ██╗
██╔════╝ ██╔════╝╚══██╔══╝        ██╔═══██╗██║   ██║██╔════╝██╔══██╗╚██╗ ██╔╝
██║  ███╗█████╗     ██║           ██║   ██║██║   ██║█████╗  ██████╔╝ ╚████╔╝
██║   ██║██╔══╝     ██║           ██║▄▄ ██║██║   ██║██╔══╝  ██╔══██╗  ╚██╔╝
╚██████╔╝███████╗   ██║           ╚██████╔╝╚██████╔╝███████╗██║  ██║   ██║
 ╚═════╝ ╚══════╝   ╚═╝            ╚══▀▀═╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝   ╚═╝
*/

// GetQuery is the type of query that does not capture information.
// It will not do operation. Only gets the existing information.
// Take two parameters.
//   1. "setQueryAdress":
//          The name of this query. This adrese request is sent.
//   2. "vasualFlag":
//            The task is visual. If it is sent as 0,
//            the value is only processed but not shown to the user.
//            But if it is 1, the incoming message is seen by the first
//            user and then processed.
func (q Query) GetQuery(setQueryAddress string, visualFlag bool) ([]byte, error) {
	// Query with the incoming address value and assign it to the variable "request"
	requestGet, errGet := http.NewRequest(GET, setQueryAddress, nil)
	errc.ErrorCenter(requestGetTag, errGet)

	if requestGet != nil {
		// This header is not automatically taken from the outside.

		responseGet, errDo := http.DefaultClient.Do(requestGet)
		errc.ErrorCenter(doGetTag, errDo)

		if responseGet != nil {
			defer responseGet.Body.Close()

			switch responseGet.StatusCode {
			case ResponseOKCode, ResponseCreatedCode:
				responseBodyGet, errBody := ioutil.ReadAll(responseGet.Body)
				errc.ErrorCenter(bodyGetTag, errBody)

				if visualFlag == Visible {
					fmt.Println(string(responseBodyGet))
				}
				return responseBodyGet, nil

			case ResponseBadRequestCode:
				return nil, ErrorResponseBadRequestCode400

			case ResponseUnauthorizedCode:
				return nil, ErrorResponseUnauthorizedCode401

			case ResponseForbiddenCode:
				return nil, ErrorResponseForbiddenCode403

			case ResponseNotFoundCode:
				return nil, ErrorNotFound404

			case ResponseServerProblemCode:
				return nil, ErrorServerProblemCode500

			default:
				return nil, ErrorElseProblem
			}
		}
		errc.ErrorCenter(requestGetTag, ErrorResponseNil)
		return nil, ErrorResponseNil
	}
	errc.ErrorCenter(requestGetTag, ErrorResponseNilRequest)
	return nil, ErrorResponseNilRequest
}

/*
██████╗  ██████╗ ███████╗████████╗     ██████╗ ██╗   ██╗███████╗██████╗ ██╗   ██╗
██╔══██╗██╔═══██╗██╔════╝╚══██╔══╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗╚██╗ ██╔╝
██████╔╝██║   ██║███████╗   ██║       ██║   ██║██║   ██║█████╗  ██████╔╝ ╚████╔╝
██╔═══╝ ██║   ██║╚════██║   ██║       ██║▄▄ ██║██║   ██║██╔══╝  ██╔══██╗  ╚██╔╝
██║     ╚██████╔╝███████║   ██║       ╚██████╔╝╚██████╔╝███████╗██║  ██║   ██║
╚═╝      ╚═════╝ ╚══════╝   ╚═╝        ╚══▀▀═╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝   ╚═╝
*/

// PostQuery is used to perform an event. For example turning off the device, running the application etc. .
// You need to be careful when using this query.
// Because it can make changes for many devices in the wrong steps.
// You need to send the parameters correctly. Otherwise, unwanted situations may occur.
// First test it on a test device in your hand.
// This query takes 4 parameters.
//    1: "setQueryAdress":
//           This query will show the address to be made.
//    2: "setBody":
//           This will show the parameters to be added to the body part of the address to be sent.
//    3: "setHeader":
//           Shows the portion of the address to be sent. The example access token is sent in this way.
//    4: "vasualFlag":
//           The task is visual. If it is not 0, the value is only processed and not shown to the user.
//           But if it is 1, the incoming message is seen by the first user and then processed.
func (q Query) PostQuery(setQueryAddress string, setBody string, setHeader map[string]string, visualFlag bool) ([]byte, error) {
	var (
		requestPost *http.Request
		errPost     error
	)

	// The possibility of whether or not the body is in question is checked.
	if setBody == "" {
		requestPost, errPost = http.NewRequest(POST, setQueryAddress, nil)
		errc.ErrorCenter(requestPostTag, errPost)
	} else {
		body := strings.NewReader(setBody)
		requestPost, errPost = http.NewRequest(POST, setQueryAddress, body)
		errc.ErrorCenter(requestPostTag, errPost)
	}
	if requestPost != nil {

		// The possibility of whether or not the header is in question is checked.
		if setHeader != nil {
			for key, value := range setHeader {
				requestPost.Header.Set(key, value)
			}
		}

		// Query based on given information
		responsePost, errDo := http.DefaultClient.Do(requestPost)
		errc.ErrorCenter(doPostTag, errDo)

		if responsePost != nil {
			defer responsePost.Body.Close()
			switch responsePost.StatusCode {
			case ResponseCreatedCode, ResponseOKCode:
				responseBodyPost, errBody := ioutil.ReadAll(responsePost.Body)
				errc.ErrorCenter(bodyPostTag, errBody)

				if visualFlag == Visible {
					fmt.Println(string(responseBodyPost))
				}
				return responseBodyPost, nil

			case ResponseBadRequestCode:
				return nil, ErrorResponseBadRequestCode400

			case ResponseUnauthorizedCode:
				return nil, ErrorResponseUnauthorizedCode401

			case ResponseForbiddenCode:
				return nil, ErrorResponseForbiddenCode403

			case ResponseNotFoundCode:
				return nil, ErrorNotFound404

			case ResponseServerProblemCode:
				return nil, ErrorServerProblemCode500

			default:
				return nil, ErrorElseProblem
			}
		}
		errc.ErrorCenter(requestPostTag, ErrorResponseNil)
		return nil, ErrorResponseNil
	}
	errc.ErrorCenter(requestPostTag, ErrorResponseNilRequest)
	return nil, ErrorResponseNilRequest
}

/*
██████╗ ██╗   ██╗████████╗     ██████╗ ██╗   ██╗███████╗██████╗ ██╗   ██╗
██╔══██╗██║   ██║╚══██╔══╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗╚██╗ ██╔╝
██████╔╝██║   ██║   ██║       ██║   ██║██║   ██║█████╗  ██████╔╝ ╚████╔╝
██╔═══╝ ██║   ██║   ██║       ██║▄▄ ██║██║   ██║██╔══╝  ██╔══██╗  ╚██╔╝
██║     ╚██████╔╝   ██║       ╚██████╔╝╚██████╔╝███████╗██║  ██║   ██║
╚═╝      ╚═════╝    ╚═╝        ╚══▀▀═╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝   ╚═╝
*/

// PutQuery is used to perform an event. For example turning off the device, running the application etc. .
// You need to be careful when using this query.
// Because it can make changes for many devices in the wrong steps.
// You need to send the parameters correctly. Otherwise, unwanted situations may occur.
// First test it on a test device in your hand.
// This query takes 4 parameters.
//    1: "setQueryAdress":
//           This query will show the address to be made.
//    2: "setBody":
//           This will show the parameters to be added to the body part of the address to be sent.
//    3: "setHeader":
//           Shows the portion of the address to be sent. The example access token is sent in this way.
//    4: "vasualFlag":
//           The task is visual. If it is not 0, the value is only processed and not shown to the user.
//           But if it is 1, the incoming message is seen by the first user and then processed.
func (q Query) PutQuery(setQueryAddress string, setBody string, setHeader map[string]string, visualFlag bool) ([]byte, error) {
	var (
		requestPut *http.Request
		errPut     error
	)

	if setBody == "" {
		requestPut, errPut = http.NewRequest(PUT, setQueryAddress, nil)
		errc.ErrorCenter(requestPutTag, errPut)
	} else {
		body := strings.NewReader(setBody)
		requestPut, errPut = http.NewRequest(PUT, setQueryAddress, body)
		errc.ErrorCenter(requestPutTag, errPut)
	}

	if requestPut != nil {

		if setHeader != nil {
			for key, value := range setHeader {
				requestPut.Header.Set(key, value)
			}
		}

		responsePut, errDo := http.DefaultClient.Do(requestPut)
		errc.ErrorCenter(doPutTag, errDo)
		defer responsePut.Body.Close()

		if responsePut != nil {

			switch responsePut.StatusCode {
			case ResponseCreatedCode, ResponseOKCode:
				responseBodyPut, errBody := ioutil.ReadAll(responsePut.Body)
				errc.ErrorCenter(bodyPutTag, errBody)

				if visualFlag == Visible {
					fmt.Println(string(responseBodyPut))
				}
				return responseBodyPut, nil

			case ResponseBadRequestCode:
				return nil, ErrorResponseBadRequestCode400
			case ResponseUnauthorizedCode:
				return nil, ErrorResponseUnauthorizedCode401

			case ResponseForbiddenCode:
				return nil, ErrorResponseForbiddenCode403

			case ResponseNotFoundCode:
				return nil, ErrorNotFound404

			case ResponseServerProblemCode:
				return nil, ErrorServerProblemCode500

			default:
				return nil, ErrorElseProblem
			}
		}
		errc.ErrorCenter(requestPutTag, ErrorResponseNil)
		return nil, ErrorResponseNil
	}
	errc.ErrorCenter(requestPutTag, ErrorResponseNilRequest)
	return nil, ErrorResponseNilRequest
}
