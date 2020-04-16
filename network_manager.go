package netmgr

import (
	"github.com/godbus/dbus/v5"
	"github.com/nlepage/go-netmgr/internal/dbusext"
)

// BusName of NetworkManager.
const BusName = "org.freedesktop.NetworkManager"

// NetworkManagerPath is the Connection Manager path.
const NetworkManagerPath = "/org/freedesktop/NetworkManager"

// NetworkManagerInterface is the Connection Manager interface.
const NetworkManagerInterface = "org.freedesktop.NetworkManager"

type (
	// NetworkManager is the Connection Manager.
	//
	// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html for more information.
	NetworkManager interface {
		dbus.BusObject

		// Methods

		Reload(flags uint) error
		GetDevices() ([]Device, error)
		GetAllDevices() ([]Device, error)
		GetDeviceByIPIface(iface string) (Device, error)
		ActivateConnection(connection interface{}, device interface{}, specificObject interface{}) (ConnectionActive, error)
		DeactivateConnection(activeConnection interface{}) error

		// Properties
		Devices() ([]Device, error)
		AllDevices() ([]Device, error)
		CheckPoints() ([]Checkpoint, error)
		NetworkingEnabled() (bool, error)
		WirelessEnabled() (bool, error)
		SetWirelessEnabled(bool) error
		WirelessHardwareEnabled() (bool, error)
		WwanEnabled() (bool, error)
		SetWwanEnabled(bool) error
		WwanHardwareEnabled() (bool, error)
		WimaxEnabled() (bool, error)
		SetWimaxEnabled(bool) error
		WimaxHardwareEnabled() (bool, error)
		ActiveConnections() ([]ConnectionActive, error)
		PrimaryConnection() (ConnectionActive, error)
		PrimaryConnectionType() (string, error)
		Metered() (Metered, error)
		ActivatingConnection() (ConnectionActive, error)
		Startup() (bool, error)
		Version() (string, error)
		Capabilities() ([]Capability, error)
		State() (State, error)
		Connectivity() (ConnectivityState, error)
		ConnectivityCheckAvailable() (bool, error)
		ConnectivityCheckEnabled() (bool, error)
		SetConnectivityCheckEnabled(bool) error
		ConnectivityCheckUri() (string, error)
		GlobalDnsConfiguration() (map[string]interface{}, error)
		SetGlobalDnsConfiguration(map[string]interface{}) error
	}

	networkManager struct {
		dbusext.BusObject
	}
)

var _ NetworkManager = (*networkManager)(nil)

// New returns the Connection Manager from conn.
func New(conn *dbus.Conn) NetworkManager {
	return &networkManager{dbusext.NewBusObject(conn, BusName, NetworkManagerPath)}
}

// System returns the Connection Manager from the system bus.
//
// It is equivalent to:
//  conn, err := dbus.SystemBus()
//  if err != nil {
//      return nil, err
//  }
//  nm := netmgr.New(conn)
func System() (NetworkManager, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return New(conn), nil
}

func (nm *networkManager) Reload(flags uint) error {
	return nm.CallAndStore(NetworkManagerInterface+".Reload", dbusext.Args{flags}, nil)
}

// Reload NetworkManager's configuration and perform certain updates, like flushing a cache or rewriting external state to disk.
//
// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.Reload for more information.
func Reload(flags uint) error {
	nm, err := System()
	if err != nil {
		return err
	}
	return nm.Reload(flags)
}

func (nm *networkManager) GetDevices() ([]Device, error) {
	return nm.getDevices(NetworkManagerInterface + ".GetDevices")
}

// GetDevices gets the list of realized network devices.
//
// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.GetDevices for more information.
func GetDevices() ([]Device, error) {
	nm, err := System()
	if err != nil {
		return nil, err
	}
	return nm.GetDevices()
}

func (nm *networkManager) GetAllDevices() ([]Device, error) {
	return nm.getDevices(NetworkManagerInterface + ".GetAllDevices")
}

// GetAllDevices gets the list of all network devices.
//
// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.GetAllDevices for more information.
func GetAllDevices() ([]Device, error) {
	nm, err := System()
	if err != nil {
		return nil, err
	}
	return nm.GetAllDevices()
}

func (nm *networkManager) getDevices(method string) ([]Device, error) {
	var devicesPaths []dbus.ObjectPath
	if err := nm.CallAndStore(method, nil, dbusext.Args{&devicesPaths}); err != nil {
		return nil, err
	}

	devices := make([]Device, 0, len(devicesPaths))
	for _, path := range devicesPaths {
		devices = append(devices, NewDevice(nm.Conn, path))
	}

	return devices, nil
}

func (nm *networkManager) GetDeviceByIPIface(iface string) (Device, error) {
	var path dbus.ObjectPath
	if err := nm.CallAndStore(NetworkManagerInterface+".GetDeviceByIpIface", dbusext.Args{iface}, dbusext.Args{&path}); err != nil {
		return nil, err
	}
	return NewDevice(nm.Conn, path), nil
}

// GetDeviceByIPIface returns the object path of the network device referenced by its IP interface name.
//
// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.GetDeviceByIpIface for more information.
func GetDeviceByIPIface(iface string) (Device, error) {
	nm, err := System()
	if err != nil {
		return nil, err
	}
	return nm.GetDeviceByIPIface(iface)
}

func (nm *networkManager) ActivateConnection(connection interface{}, device interface{}, specificObject interface{}) (ConnectionActive, error) {
	connectionPath, err := dbusext.ObjectPath(connection)
	if err != nil {
		return nil, err
	}
	devicePath, err := dbusext.ObjectPath(device)
	if err != nil {
		return nil, err
	}
	specificObjectPath, err := dbusext.ObjectPath(specificObject)
	if err != nil {
		return nil, err
	}

	var connectionActivePath dbus.ObjectPath
	if err := nm.CallAndStore(NetworkManagerInterface+".ActivateConnection", dbusext.Args{connectionPath, devicePath, specificObjectPath}, dbusext.Args{&connectionActivePath}); err != nil {
		return nil, err
	}

	return NewConnectionActive(nm.Conn, connectionActivePath)
}

// ActivateConnection activates a connection using the supplied device.
//
// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.ActivateConnection for more information.
func ActivateConnection(connection interface{}, device interface{}, specificObject interface{}) (ConnectionActive, error) {
	nm, err := System()
	if err != nil {
		return nil, err
	}
	return nm.ActivateConnection(connection, device, specificObject)
}

func (nm *networkManager) DeactivateConnection(activeConnection interface{}) error {
	activeConnectionPath, err := dbusext.ObjectPath(activeConnection)
	if err != nil {
		return err
	}

	return nm.CallAndStore(NetworkManagerInterface+".DeactivateConnection", dbusext.Args{activeConnectionPath}, nil)
}

// DeactivateConnection deactivates an active connection.
//
// See https://developer.gnome.org/NetworkManager/stable/gdbus-org.freedesktop.NetworkManager.html#gdbus-method-org-freedesktop-NetworkManager.DeactivateConnection for more information.
func DeactivateConnection(activeConnection interface{}) error {
	nm, err := System()
	if err != nil {
		return err
	}
	return nm.DeactivateConnection(activeConnection)
}

// Capability names the numbers in the Capabilities property.
//
// See https://developer.gnome.org/NetworkManager/stable/nm-dbus-types.html#NMCapability for more information.
type Capability uint

const (
	// CapabilityTeam indicates teams can be managed. This means the team device plugin is loaded.
	CapabilityTeam Capability = iota + 1

	// CapabilityOVS indicates OpenVSwitch can be managed. This means the OVS device plugin is loaded. Since: 1.24, 1.22.2
	CapabilityOVS
)

// State values indicate the current overall networking state.
//
// See https://developer.gnome.org/NetworkManager/stable/nm-dbus-types.html#NMState for more information.
type State uint

// FIXME State consts

// ConnectivityState values indicate the connectivity state.
type ConnectivityState uint

// FIXME ConnectivityState consts
