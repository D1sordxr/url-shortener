package sqlutil

import (
	"database/sql"
	"net"

	"github.com/sqlc-dev/pqtype"
)

func ToNullString(s *string) sql.NullString {
	if s != nil && *s != "" {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func ToInetFromIP(ip net.IP) pqtype.Inet {
	if len(ip) > 0 {
		if ipv4 := ip.To4(); ipv4 != nil {
			return pqtype.Inet{
				IPNet: net.IPNet{
					IP:   ipv4,
					Mask: net.CIDRMask(32, 32),
				},
				Valid: true,
			}
		} else {
			return pqtype.Inet{
				IPNet: net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(128, 128),
				},
				Valid: true,
			}
		}
	}
	return pqtype.Inet{Valid: false}
}
