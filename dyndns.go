package main

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

func parseFqdn(record string) (string, string, error) {
	parts := strings.SplitN(record, ".", 2)
	if len(parts) != 2 {
		msg := fmt.Sprintf("Not a fully-qualified domain name: %s", record)
		return "", "", errors.New(msg)
	}

	return parts[0], parts[1], nil
}

func dyndnsUpdate(token, record, kind string, ip net.IP) error {
	host, domain, err := parseFqdn(record)
	if err != nil {
		return err
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, NewTokenString(token))
	doClient := godo.NewClient(oauthClient)
	listOpts := &godo.ListOptions{}

	records, _, err := doClient.Domains.Records(domain, listOpts)
	if err != nil {
		return err
	}

	var id int
	for _, r := range records {
		if r.Type == kind && r.Name == host {
			if r.Data == ip.String() {
				fmt.Printf("%s %s up to date.\n", record, kind)
				return nil
			}
			id = r.ID
			break
		}
	}

	edit := &godo.DomainRecordEditRequest{
		Type: kind,
		Name: host,
		Data: ip.String(),
	}

	if id == 0 {
		// create
		_, _, err := doClient.Domains.CreateRecord(domain, edit)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s created.\n", record, kind)
	} else {
		_, _, err := doClient.Domains.EditRecord(domain, id, edit)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s updated.\n", record, kind)
	}

	return nil
}
