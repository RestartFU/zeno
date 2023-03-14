package specs

// DeviceGroup represents a device group that a user may be in.
type DeviceGroup struct {
	unrestricted bool
	device       int
}

// DeviceGroupUnrestricted returns a device group with no restrictions.
func DeviceGroupUnrestricted() DeviceGroup {
	return DeviceGroup{unrestricted: true}
}

// DeviceGroupKeyboardMouse returns a device group that allows keyboard and mouse players.
func DeviceGroupKeyboardMouse() DeviceGroup {
	return DeviceGroup{device: 1}
}

// DeviceGroupMobile returns a device group that allows mobile players.
func DeviceGroupMobile() DeviceGroup {
	return DeviceGroup{device: 2}
}

// DeviceGroupController returns a device group that allows controller players.
func DeviceGroupController() DeviceGroup {
	return DeviceGroup{device: 3}
}

// DeviceGroups returns a list of all device groups.
func DeviceGroups() []DeviceGroup {
	return []DeviceGroup{
		DeviceGroupUnrestricted(),
		DeviceGroupKeyboardMouse(),
		DeviceGroupMobile(),
		DeviceGroupController(),
	}
}

// Name ...
func (d DeviceGroup) Name() string {
	switch d.device {
	case 1:
		return "Keyboard/Mouse"
	case 2:
		return "Mobile"
	case 3:
		return "Controller"
	}
	return "Unrestricted"
}

// Unrestricted ...
func (d DeviceGroup) Unrestricted() bool {
	return d.unrestricted
}

// Compare ...
func (d DeviceGroup) Compare(device int, unrestricted bool) bool {
	if d.unrestricted && unrestricted {
		// If both groups are unrestricted, we can match them together.
		return true
	}
	return d.device == device
}
