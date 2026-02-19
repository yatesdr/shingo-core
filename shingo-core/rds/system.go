package rds

import (
	"encoding/json"
	"fmt"
)

// Ping checks RDS Core connectivity and returns product/version info.
func (c *Client) Ping() (*PingResponse, error) {
	var resp PingResponse
	if err := c.get("/ping", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetProfiles retrieves an RDS configuration file. Returns raw JSON content.
func (c *Client) GetProfiles(file string) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := c.post("/getProfiles", &GetProfilesRequest{File: file}, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// GetLicenseInfo retrieves the current RDS license information.
func (c *Client) GetLicenseInfo() (*LicenseInfo, error) {
	var resp LicenseResponse
	if err := c.get("/licInfo", &resp); err != nil {
		return nil, err
	}
	if err := checkResponse(&resp.Response); err != nil {
		return nil, err
	}
	if resp.Data == nil {
		c.dbg("!! /licInfo returned code=0 but data=null")
		return nil, fmt.Errorf("license info: empty response data")
	}
	return resp.Data, nil
}
