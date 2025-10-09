package service

import (
	"fmt"
	"time"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/logger"
	"github.com/mhsanaei/3x-ui/v2/xray"

	"gorm.io/gorm"
)

// OutboundService provides business logic for managing Xray outbound configurations.
// It handles CRUD operations for outbounds and traffic monitoring.
type OutboundService struct{}

func (s *OutboundService) AddTraffic(traffics []*xray.Traffic, clientTraffics []*xray.ClientTraffic) (error, bool) {
	var err error
	db := database.GetDB()
	tx := db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = s.addOutboundTraffic(tx, traffics)
	if err != nil {
		return err, false
	}

	return nil, false
}

func (s *OutboundService) addOutboundTraffic(tx *gorm.DB, traffics []*xray.Traffic) error {
	if len(traffics) == 0 {
		return nil
	}

	var err error

	for _, traffic := range traffics {
		if traffic.IsOutbound {

			var outbound model.OutboundTraffics

			err = tx.Model(&model.OutboundTraffics{}).Where("tag = ?", traffic.Tag).
				FirstOrCreate(&outbound).Error
			if err != nil {
				return err
			}

			outbound.Tag = traffic.Tag
			outbound.Up = outbound.Up + traffic.Up
			outbound.Down = outbound.Down + traffic.Down
			outbound.Total = outbound.Up + outbound.Down

			err = tx.Save(&outbound).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *OutboundService) GetOutboundsTraffic() ([]*model.OutboundTraffics, error) {
	db := database.GetDB()
	var traffics []*model.OutboundTraffics

	err := db.Model(model.OutboundTraffics{}).Find(&traffics).Error
	if err != nil {
		logger.Warning("Error retrieving OutboundTraffics: ", err)
		return nil, err
	}

	return traffics, nil
}

func (s *OutboundService) ResetOutboundTraffic(tag string) error {
	db := database.GetDB()

	whereText := "tag "
	if tag == "-alltags-" {
		whereText += " <> ?"
	} else {
		whereText += " = ?"
	}

	result := db.Model(model.OutboundTraffics{}).
		Where(whereText, tag).
		Updates(map[string]any{"up": 0, "down": 0, "total": 0})

	err := result.Error
	if err != nil {
		return err
	}

	return nil
}

// ========================= Outbound Configuration Management =========================

// GetOutbounds retrieves all outbounds for a specific user.
func (s *OutboundService) GetOutbounds(userId int) ([]*model.Outbound, error) {
	db := database.GetDB()
	var outbounds []*model.Outbound
	err := db.Model(model.Outbound{}).Where("user_id = ?", userId).Find(&outbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return outbounds, nil
}

// GetAllOutbounds retrieves all outbounds from the database.
func (s *OutboundService) GetAllOutbounds() ([]*model.Outbound, error) {
	db := database.GetDB()
	var outbounds []*model.Outbound
	err := db.Model(model.Outbound{}).Find(&outbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return outbounds, nil
}

// GetEnabledOutbounds retrieves all enabled outbounds.
func (s *OutboundService) GetEnabledOutbounds() ([]*model.Outbound, error) {
	db := database.GetDB()
	var outbounds []*model.Outbound
	err := db.Model(model.Outbound{}).Where("enable = ?", true).Find(&outbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return outbounds, nil
}

// GetOutbound retrieves a specific outbound by ID.
func (s *OutboundService) GetOutbound(id int) (*model.Outbound, error) {
	db := database.GetDB()
	outbound := &model.Outbound{}
	err := db.Model(model.Outbound{}).First(outbound, id).Error
	if err != nil {
		return nil, err
	}
	return outbound, nil
}

// GetOutboundByTag retrieves a specific outbound by tag.
func (s *OutboundService) GetOutboundByTag(tag string) (*model.Outbound, error) {
	db := database.GetDB()
	outbound := &model.Outbound{}
	err := db.Model(model.Outbound{}).Where("tag = ?", tag).First(outbound).Error
	if err != nil {
		return nil, err
	}
	return outbound, nil
}

// checkTagExist checks if an outbound tag already exists (excluding a specific ID).
func (s *OutboundService) checkTagExist(tag string, ignoreId int) (bool, error) {
	db := database.GetDB()
	var count int64
	err := db.Model(model.Outbound{}).Where("tag = ? and id != ?", tag, ignoreId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AddOutbound creates a new outbound configuration.
func (s *OutboundService) AddOutbound(outbound *model.Outbound) error {
	db := database.GetDB()

	// Check if tag already exists
	exists, err := s.checkTagExist(outbound.Tag, 0)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("outbound tag already exists: %s", outbound.Tag)
	}

	// Set timestamps
	now := time.Now().Unix()
	outbound.CreatedAt = now
	outbound.UpdatedAt = now

	return db.Create(outbound).Error
}

// UpdateOutbound updates an existing outbound configuration.
func (s *OutboundService) UpdateOutbound(outbound *model.Outbound) error {
	db := database.GetDB()

	// Check if new tag conflicts with existing outbound
	exists, err := s.checkTagExist(outbound.Tag, outbound.Id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("outbound tag already exists: %s", outbound.Tag)
	}

	// Update timestamp
	outbound.UpdatedAt = time.Now().Unix()

	// Update all fields
	return db.Save(outbound).Error
}

// DelOutbound deletes an outbound by ID.
func (s *OutboundService) DelOutbound(id int) error {
	db := database.GetDB()
	return db.Delete(model.Outbound{}, id).Error
}

// DelOutboundByTag deletes an outbound by tag.
func (s *OutboundService) DelOutboundByTag(tag string) error {
	db := database.GetDB()
	return db.Where("tag = ?", tag).Delete(model.Outbound{}).Error
}

// GetOutboundTags retrieves all outbound tags as a JSON array string.
func (s *OutboundService) GetOutboundTags() (string, error) {
	db := database.GetDB()
	var tags []string
	err := db.Model(model.Outbound{}).Pluck("tag", &tags).Error
	if err != nil {
		return "", err
	}

	// Build JSON array manually
	result := "["
	for i, tag := range tags {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("\"%s\"", tag)
	}
	result += "]"
	return result, nil
}
