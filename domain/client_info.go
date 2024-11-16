package domain

import "context"

type ClientInfoResponse struct {
	Device    string `json:"device"`
	Location  string `json:"location"`
	LoginTime string `json:"login_time"`
	IP        string `json:"ip"`
}

type ClientInfoService interface {
	GetClientInfo(ctx context.Context) (*ClientInfoResponse, error)
}
