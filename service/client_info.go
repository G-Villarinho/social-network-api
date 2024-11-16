package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/G-Villarinho/social-network/config"
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/internal"
	jsoniter "github.com/json-iterator/go"
	"github.com/mssola/user_agent"
)

type result struct {
	City    string `json:"city"`
	Region  string `json:"region_name"`
	Country string `json:"country_name"`
}

type clientInfoService struct {
	di             *internal.Di
	contextService domain.ContextService
}

func NewClientInfoService(di *internal.Di) (domain.ClientInfoService, error) {
	contextService, err := internal.Invoke[domain.ContextService](di)
	if err != nil {
		return nil, err
	}

	return &clientInfoService{
		di:             di,
		contextService: contextService,
	}, nil
}

func (c *clientInfoService) GetClientInfo(ctx context.Context) (*domain.ClientInfoResponse, error) {
	userAgent, err := c.contextService.GetUserAgent(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user-agent: %w", err)
	}

	clientIP, err := c.contextService.GetClientIP(ctx)
	if err != nil {
		return nil, fmt.Errorf("get client-ip: %w", err)
	}

	location, err := getLocation(clientIP)
	if err != nil {
		location = "Unknown"
	}

	clientInfo := domain.ClientInfoResponse{
		Device:    getDevice(userAgent),
		Location:  location,
		LoginTime: time.Now().UTC().Format("January 2, 2006, 3:04 PM"),
		IP:        clientIP,
	}

	return &clientInfo, nil
}

func getDevice(userAgent string) string {
	ua := user_agent.New(userAgent)
	name, version := ua.Browser()
	return fmt.Sprintf("%s (%s %s)", name, ua.OS(), version)
}

func getLocation(ip string) (string, error) {
	url := fmt.Sprintf("%s/%s?access_key=%s", config.Env.IpStacker.IpStackBaseURL, ip, config.Env.IpStacker.IpStackAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("fetch location: %w", err)
	}
	defer resp.Body.Close()

	var result result
	if err := jsoniter.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Warn("error decoding location response", slog.String("error", err.Error()))
		return "", fmt.Errorf("decode location response: %w", err)
	}

	return fmt.Sprintf("%s, %s, %s", result.City, result.Region, result.Country), nil
}
