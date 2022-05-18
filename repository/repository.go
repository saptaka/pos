package repository

import (
	"context"
	"time"

	"github.com/saptaka/pos/config"
)

type Repo interface {
	CashierRepo
	CategoryRepo
	ProductRepo
	PaymentRepo
	OrderRepo
	ReportRepo
	SetupTableStructure()
}

type repo struct {
	db DB
}

func NewRepository(cfg *config.Config) Repo {
	database := newDatabase(cfg.DBuser, cfg.DBPassword, cfg.DBHost, cfg.DBName,
		cfg.DBPort, cfg.DBMaxConnection, cfg.DBMaxIdle,
		time.Duration(cfg.DBConnectionTimeout))
	return &repo{database}
}

func (r repo) SetupTableStructure() {
	cashiersTable := `CREATE TABLE IF NOT EXISTS cashiers (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT 'DEFAULT',
		passcode varchar(255) CHARACTER SET utf8mb4  NOT NULL,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY id (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4;`

	categoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) CHARACTER SET utf8mb4  NOT NULL,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY id (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 ;
	`

	discountsTable := `
	CREATE TABLE IF NOT EXISTS discounts (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		qty int NOT NULL DEFAULT '0',
		types varchar(255) CHARACTER SET utf8mb4  DEFAULT NULL,
		result int DEFAULT NULL,
		expired_at timestamp NULL DEFAULT NULL,
		expired_at_format varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT 'DEFAULT',
		string_format varchar(255) CHARACTER SET utf8mb4  DEFAULT 'DEFAULT',
		UNIQUE KEY id (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 ;`

	ordersTable := `
	  CREATE TABLE  IF NOT EXISTS orders (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		cashier_id bigint unsigned DEFAULT NULL,
		payment_type_id bigint unsigned DEFAULT NULL,
		receipt_id varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT '',
		total_price int NOT NULL DEFAULT '0',
		total_paid int NOT NULL DEFAULT '0',
		total_return int NOT NULL DEFAULT '0',
		receipt_file_path varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT '',
		is_downloaded tinyint NOT NULL DEFAULT '0',
		UNIQUE KEY id (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ; 
	  `

	paymentsTable := `
	  CREATE TABLE  IF NOT EXISTS payments (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		types varchar(255) CHARACTER SET utf8mb4  DEFAULT NULL,
		logo varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT 'DEFAULT',
		name varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT 'DEFAULT',
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY id (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 ; 
	  `

	productsTable := `
	CREATE TABLE  IF NOT EXISTS products (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		sku varchar(5) CHARACTER SET utf8mb4  NOT NULL DEFAULT '' COMMENT '',
		stock int DEFAULT NULL,
		price int DEFAULT NULL,
		image varchar(255) CHARACTER SET utf8mb4  NOT NULL DEFAULT '',
		discount_id bigint unsigned DEFAULT NULL,
		category_id bigint unsigned DEFAULT NULL,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY id (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4;
	  `

	orderedProductsTable := `
	  CREATE TABLE  IF NOT EXISTS ordered_products (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		product_id bigint unsigned NOT NULL,
		order_id bigint unsigned NOT NULL,
		qty int DEFAULT NULL,
		total_normal_price int DEFAULT NULL,
		total_final_price int DEFAULT NULL,
		discount_id bigint unsigned DEFAULT NULL,
		price_product int DEFAULT NULL,
		UNIQUE KEY id (id),
		KEY fk_product_id (product_id),
		KEY fk_order_id (order_id),
		KEY fk_ordered_discount (discount_id),
		CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders (id),
		CONSTRAINT fk_ordered_discount FOREIGN KEY (discount_id) REFERENCES discounts (id),
		CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products (id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 ; 
	  `

	_, err := r.db.ExecContext(context.Background(), cashiersTable)
	if err != nil {
		panic(err)
	}
	_, err = r.db.ExecContext(context.Background(), categoriesTable)
	if err != nil {
		panic(err)
	}

	_, err = r.db.ExecContext(context.Background(), discountsTable)
	if err != nil {
		panic(err)
	}

	_, err = r.db.ExecContext(context.Background(), paymentsTable)
	if err != nil {
		panic(err)
	}

	_, err = r.db.ExecContext(context.Background(), ordersTable)
	if err != nil {
		panic(err)
	}

	_, err = r.db.ExecContext(context.Background(), productsTable)
	if err != nil {
		panic(err)
	}

	_, err = r.db.ExecContext(context.Background(), orderedProductsTable)
	if err != nil {
		panic(err)
	}
}
