package lib

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestUnity() *Unity {
	return NewUnity("10.0.0.5", "admin", "password")
}

func TestUnity_Ip(t *testing.T) {
	assert.Equal(t, getTestUnity().Ip(), "10.0.0.5")
}

func TestUnity_Username(t *testing.T) {
	assert.Equal(t, getTestUnity().Username(), "admin")
}

func TestGetSystem(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	system := NewUnityByConn(mockConn())
	Update(system)

	asserts := assert.New(t)
	asserts.Equal("0", system.Id)
	asserts.Equal("FNM00150600267", system.Name)
	asserts.Equal("Unity 500", system.Model)
	asserts.Equal("FNM00150600267", system.SerialNumber)
	asserts.Equal("Oberon_DualSP", system.Platform)
	asserts.Equal("08:00:1B:FF:16:F9", system.MacAddress)
	asserts.True(system.IsEULAAccepted)
	asserts.False(system.IsUpgradeComplete)
	asserts.True(system.IsAutoFailbackEnabled)
	asserts.Equal(503, system.CurrentPower)
	asserts.Equal(504, system.AvgPower)

	health := system.Health
	asserts.Equal(5, health.Value)
	asserts.Contains(health.DescriptionIds, "ALRT_SYSTEM_OK")
	asserts.Contains(health.Descriptions, "The system is operating normally.")
}
