// Copyright (c) 2020 Proton Technologies AG
//
// This file is part of ProtonMail Bridge.Bridge.
//
// ProtonMail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// ProtonMail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with ProtonMail Bridge.  If not, see <https://www.gnu.org/licenses/>.

package liveapi

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (cntrl *Controller) TurnInternetConnectionOff() {
	cntrl.noInternetConnection = true
}

func (cntrl *Controller) TurnInternetConnectionOn() {
	cntrl.noInternetConnection = false
}

type fakeTransport struct {
	cntrl     *Controller
	transport http.RoundTripper
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.cntrl.noInternetConnection {
		return nil, errors.New("no route to host")
	}

	body := []byte{}
	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get body")
		}
		body, err = ioutil.ReadAll(bodyReader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read body")
		}
	}
	t.cntrl.recordCall(req.Method, req.URL.Path, body)

	return t.transport.RoundTrip(req)
}
