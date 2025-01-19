package store

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/ip2location/ip2location-go/v9"
)

//	db, err := ip2location.OpenDB("./IP-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE-ADDRESSTYPE-CATEGORY-DISTRICT-ASN-DB26.BIN")
//
//	if err != nil {
//		fmt.Print(err)
//		return
//	}
//	ip := "45.5.164.0"
//	results, err := db.Get_all(ip)
//
//	if err != nil {
//		fmt.Print(err)
//		return
//	}
//
//	fmt.Printf("country_short: %s\n", results.Country_short)
//	fmt.Printf("country_long: %s\n", results.Country_long)
//	fmt.Printf("region: %s\n", results.Region)
//	fmt.Printf("city: %s\n", results.City)
//	fmt.Printf("isp: %s\n", results.Isp)

type IpStore interface {
	GetCountryByIP(ip string) (*models.IPInfo, error)
}

type ipStore struct {
	db *ip2location.DB
}

func NewIpStore(dbPath string) (IpStore, error) {
	db, err := ip2location.OpenDB(dbPath)
	if err != nil {
		return nil, err
	}
	return &ipStore{db: db}, nil
}

func (i *ipStore) GetCountryByIP(ip string) (*models.IPInfo, error) {
	results, err := i.db.Get_all(ip)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return &models.IPInfo{
		IP:          ip,
		CountryCode: results.Country_short,
		CountryName: results.Country_long,
	}, nil
}
