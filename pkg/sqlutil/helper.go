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

func ToNullInt32(i *int32) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: *i, Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

func ToNullInt64(i *int64) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: *i, Valid: true}
	}
	return sql.NullInt64{Valid: false}
}

func ToNullBool(b *bool) sql.NullBool {
	if b != nil {
		return sql.NullBool{Bool: *b, Valid: true}
	}
	return sql.NullBool{Valid: false}
}

func ToInet(ipAddr *net.IPAddr) pqtype.Inet {
	if ipAddr != nil && ipAddr.IP != nil {
		ip := ipAddr.IP
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

func ToInetFromIP(ip net.IP) pqtype.Inet {
	if ip != nil && len(ip) > 0 {
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

func ToInetFromString(ipStr string) pqtype.Inet {
	if ipStr == "" {
		return pqtype.Inet{Valid: false}
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return pqtype.Inet{Valid: false}
	}

	return ToInetFromIP(ip)
}

func ToIPAddrFromInet(inet pqtype.Inet) *net.IPAddr {
	if !inet.Valid {
		return nil
	}
	return &net.IPAddr{
		IP:   inet.IPNet.IP,
		Zone: "",
	}
}
