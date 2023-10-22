package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	Config *StructConfig
)

type StructConfig struct {
	Health struct {
		HealingPotionAt     int `yaml:"healingPotionAt"`
		ManaPotionAt        int `yaml:"manaPotionAt"`
		RejuvPotionAtLife   int `yaml:"rejuvPotionAtLife"`
		RejuvPotionAtMana   int `yaml:"rejuvPotionAtMana"`
		MercHealingPotionAt int `yaml:"mercHealingPotionAt"`
		MercRejuvPotionAt   int `yaml:"mercRejuvPotionAt"`
	} `yaml:"health"`
	Bindings struct {
		PotionHP    int `yaml:"potionHP"`
		PotionMANA  int `yaml:"potionMANA"`
		PotionREJUV int `yaml:"potionREJUV"`
	} `yaml:"bindings"`
	Timings struct {
		RejuvInterval       float32 `yaml:"rejuvInterval"`
		HealingInterval     float32 `yaml:"healingInterval"`
		ManaInterval        float32 `yaml:"manaInterval"`
		HealingMercInterval float32 `yaml:"healingMercInterval"`
	} `yaml:"timings"`
}

// Load reads the config.ini file and returns a Config struct filled with data from the ini file
func Load() error {
	r, err := os.Open("config/config.yaml")
	if err != nil {
		return fmt.Errorf("error loading config.yaml: %w", err)
	}

	d := yaml.NewDecoder(r)
	if err = d.Decode(&Config); err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}
	return nil
}
