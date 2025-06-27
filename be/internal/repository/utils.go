package repository

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lolwierd/weatherboy/be/internal/db"
)

// GetConn retrieves a connection from the global pgx pool.
func GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	dbDriver := db.GetDBDriver()
	return dbDriver.ConnPool.Acquire(ctx)
}

// GetConnTransaction starts a transaction using the global connection pool.
// The returned connection is always nil because `pgxpool` manages connection lifecycles internally.
func GetConnTransaction(ctx context.Context) (conn *pgxpool.Conn, tx pgx.Tx, err error) {
	dbDriver := db.GetDBDriver()
	tx, err = dbDriver.ConnPool.BeginTx(ctx, pgx.TxOptions{})
	return nil, tx, err
}

func generateIPv6FromMAC(ipv6Cidr netip.Prefix, mac string) string {
	// Get the IPv6 prefix
	prefix := ipv6Cidr.Addr()

	// Convert MAC to EUI-64
	macAddr, _ := net.ParseMAC(mac)
	eui64 := macToEUI64(macAddr)

	// Combine prefix with EUI-64
	ipv6 := prefix.As16()
	copy(ipv6[8:], eui64[:])

	// Create and return the IPv6 address
	ip := netip.AddrFrom16(ipv6).String()
	return fmt.Sprintf("%s/%d", ip, ipv6Cidr.Bits())
}

// Helper function to convert MAC to EUI-64
func macToEUI64(mac net.HardwareAddr) [8]byte {
	var eui64 [8]byte
	copy(eui64[:3], mac[:3])
	copy(eui64[5:], mac[3:])
	eui64[3] = 0xFF
	eui64[4] = 0xFE
	eui64[0] ^= 0x02 // Flip Universal/Local bit
	return eui64
}

func getGatewayIP(cidr netip.Prefix) string {
	if !cidr.Addr().Is4() {
		return cidr.Addr().String()
	}
	gatewayIP := cidr.Addr().Next()
	return gatewayIP.String()
}
