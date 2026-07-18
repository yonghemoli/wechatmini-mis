package db

// AutoMigrate 自动迁移 MIS 基础表。
func AutoMigrate() error {
	if Get() == nil {
		return nil
	}
	err := Get().AutoMigrate(
		&AdminDO{},
		&CustomerDO{},
		&AppConfigDO{},
		&FAQDO{},
		&ChatSessionDO{},
		&ChatMessageDO{},
		&MiniServiceCategoryDO{},
		&CaregiverDO{},
		&DemandDO{},
		&ResumeDO{},
		&BusinessStatusHistoryDO{},
	)
	if err != nil {
		return err
	}

	return nil
}
