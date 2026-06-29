package db

// AutoMigrate 自动迁移 MIS 基础表。
func AutoMigrate() error {
	if Get() == nil {
		return nil
	}
	err := Get().AutoMigrate(
		&AdminDO{},
		&OrderDO{},
		&CustomerDO{},
		&ServiceTypeDO{},
		&ServiceDO{},
		&ShopDO{},
		&AddressDO{},
		&ServiceTargetDO{},
		&DishDO{},
		&MealPackageDO{},
		&AppConfigDO{},
		&FAQDO{},
		&ChatSessionDO{},
		&ChatMessageDO{},
	)
	if err != nil {
		return err
	}

	return nil
}
