package mysql

import (
	"fmt"
	"net"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

type MySQL struct {
	DB *gorm.DB
}

type Model struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func New(cfg MySQLConfig) (*MySQL, error) {
	sql := &MySQL{}

	addr := cfg.Host
	if addr == "" {
		addr = "127.0.0.1:3306"
	}
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "3306")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, addr, cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql - New - gorm.Open: %w", err)
	}

	stdDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("mysql - New - get std db: %w", err)
	}
	if err := stdDB.Ping(); err != nil {
		return nil, fmt.Errorf("mysql - New - ping: %w", err)
	}

	sql.DB = db
	return sql, nil
}

func (p *MySQL) Close() {
	if p.DB != nil {
		if stdDB, err := p.DB.DB(); err == nil {
			_ = stdDB.Close()
		}
	}
}
