package internal

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfFile(t *testing.T) {
	config :=
		`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	}
}
	`
	t.Run("read config file ", func(t *testing.T) {
		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)

	})

	t.Run("change permissions of file", func(t *testing.T) {
		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), fs.FileMode(os.O_RDONLY))
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.Error(t, err)
		assert.Empty(t, data)

	})

	t.Run("no file exists", func(t *testing.T) {

		err := os.WriteFile("./config.json", []byte(config), fs.FileMode(os.O_RDONLY))
		assert.Error(t, err)

		data, err := ReadConfFile("./config.json")
		assert.Error(t, err)
		assert.Empty(t, data)

	})

}

func TestParseConf(t *testing.T) {

	t.Run("parse config file", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "email",
        "sendgrid_key": "my sendgrid_key",
        "timeout": 60 
    },
    "account": {
        "mnemonics": "my mnemonics",
		"network": "my network"
    },
	"token": {
        "secret": "secret",
        "timeout": 10
    },
	"database": {
        "file": "testing.db"
    },
	"version": "v1",
	"salt": "salt"
}
	`
		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		expected := Configuration{
			Server: Server{
				Host: "localhost",
				Port: ":3000",
			},
			MailSender: MailSender{
				Email:       "email",
				SendGridKey: "my sendgrid_key",
				Timeout:     60,
			},
			Account: GridAccount{
				Mnemonics: "my mnemonics",
			},
			Token: JwtToken{
				Secret:  "secret",
				Timeout: 10,
			},
			Database: DB{
				File: "testing.db",
			},
			Version: "v1",
		}

		got, err := ParseConf(data)
		assert.NoError(t, err)
		assert.Equal(t, got.Server, expected.Server)
		assert.Equal(t, got.MailSender.Email, expected.MailSender.Email)
		assert.Equal(t, got.Account.Mnemonics, expected.Account.Mnemonics)
		assert.Equal(t, got.Token, expected.Token)
		assert.Equal(t, got.Database, expected.Database)
		assert.Equal(t, got.Version, expected.Version)
	})

	t.Run("no file", func(t *testing.T) {
		_, err := ReadConfFile("config.json")
		assert.Error(t, err)

	})

	t.Run("no server configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "",
		"port": ""
	}
}
	`

		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "server configuration is required")

	})

	//TODO: error
	t.Run("no mail sender configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "",
        "sendgrid_key": "",
        "timeout": 0
    }
}
	`

		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "mail sender configuration is required")

	})

	t.Run("no database configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "email",
        "sendgrid_key": "my sendgrid_key",
        "timeout": 60 
    },
	"database": {
        "file": ""
    }
}
	`

		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "database configuration is required")

	})

	t.Run("no account configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "email",
        "sendgrid_key": "my sendgrid_key",
        "timeout": 60 
    },
	"database": {
        "file": "testing.db"
    },
    "account": {
        "mnemonics": "",
		"network": ""
    }
}
	`

		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "account configuration is required")

	})

	t.Run("no jwt token configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "email",
        "sendgrid_key": "my sendgrid_key",
        "timeout": 60 
    },
	"database": {
        "file": "testing.db"
    },
    "account": {
        "mnemonics": "my mnemonics",
		"network": "my network"
    },	
	"token": {
        "secret": "",
        "timeout": 0
    }
}
	`

		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "jwt token configuration is required")

	})

	t.Run("no version configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "email",
        "sendgrid_key": "my sendgrid_key",
        "timeout": 60 
    },
    "account": {
        "mnemonics": "my mnemonics",
		"network": "my network"
    },
	"database": {
        "file": "testing.db"
    },
	"token": {
        "secret": "secret",
        "timeout": 10
    },	
	"version": ""
}
	`

		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "version is required")

	})

	t.Run("no salt configuration", func(t *testing.T) {
		config :=
			`
{
	"server": {
		"host": "localhost",
		"port": ":3000"
	},
	"mailSender": {
        "email": "email",
        "sendgrid_key": "my sendgrid_key",
        "timeout": 60 
    },
    "account": {
        "mnemonics": "my mnemonics",
		"network": "my network"
    },
	"database": {
        "file": "testing.db"
    },
	"token": {
        "secret": "secret",
        "timeout": 10
    },	
	"version": "v1",	
	"salt": ""
}
	`
		dir := t.TempDir()
		configPath := dir + "/config.json"

		err := os.WriteFile(configPath, []byte(config), 0644)
		assert.NoError(t, err)

		data, err := ReadConfFile(configPath)
		assert.NoError(t, err)

		_, err = ParseConf(data)
		assert.Error(t, err, "salt is required")

	})
}
