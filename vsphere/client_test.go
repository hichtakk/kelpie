package vsphere

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hichtakk/kelpie/vsphere/mock"
)

func SetupDefaultFixtures() {

}

func SetupOptionalFixtures() {

}

func TestSetCredential(t *testing.T) {
	vsc := NewVSphereClient(true)
	vsc.SetCredential("fake-user", "fake-password")
}

func TestValidatePath(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		{"/hello/", errors.New("request path must start with \"/api/\"")},
		{"/api/", nil},
		{"/api/session", nil},
		{"/api/vcenter", nil},
		{"/api/vcenter/vm", nil},
	}
	vsc := NewVSphereClient(true)
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := vsc.validatePath(tc.input)
			if got != tc.want {
				if got.Error() != tc.want.Error() {
					t.Errorf("error mismatch want: '%s', got: '%s'", tc.want.Error(), got.Error())
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	cases := []struct {
		name string
		res  *http.Response
		want error
	}{
		{
			"network error",
			&http.Response{},
			errors.New("http error"),
		},
		{
			"authentication error",
			&http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			errors.New("authentication failed"),
		},
		{
			"response error",
			&http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			fmt.Errorf("session id not found"),
		},
		{
			"normal",
			&http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(strings.NewReader("")),
				Header: http.Header{
					"Vmware-Api-Session-Id": []string{"test-token"},
				},
			},
			nil,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	vsc := NewVSphereClient(true)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockHttpClient := mock.NewMockHttpClient(mockCtrl)
			mockHttpClient.EXPECT().Do(gomock.Any()).Return(tc.res, tc.want)
			vsc.SetHttpClient(mockHttpClient)
			if got := vsc.Login(); got != tc.want {
				t.Errorf("test case failed %T", tc)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	cases := []struct {
		name string
		res  *http.Response
		want error
	}{
		{
			"network error",
			&http.Response{},
			errors.New("network error"),
		},
		{
			"http request error",
			&http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			errors.New("logout failed"),
		},
		{
			"normal",
			&http.Response{
				StatusCode: 204,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			nil,
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	vsc := NewVSphereClient(true)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockHttpClient := mock.NewMockHttpClient(mockCtrl)
			mockHttpClient.EXPECT().Do(gomock.Any()).Return(tc.res, tc.want)
			vsc.SetHttpClient(mockHttpClient)
			if got := vsc.Logout(); got != tc.want {
				t.Errorf("test case failed %T", tc)
			}
		})
	}
}
